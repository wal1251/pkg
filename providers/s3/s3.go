// Package s3
// Представляет собой адаптер s3-клиента
package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/size"
)

var _ ClientInterface = (*Client)(nil)

type (
	ClientInterface interface {
		// GetListFiles возвращает слайс объектов c Body по заданному списку ключей.
		GetListFiles(ctx context.Context, keys []string) ([]*StorageObject, error)

		// GetFile возвращает объект с Body по указанному key.
		GetFile(ctx context.Context, key string) (*StorageObject, error)

		// ListFiles возвращает список объектов совпадающих по префиксу.
		ListFiles(ctx context.Context, prefix string) ([]*StorageObject, error)

		// FindFile ищет файл по совпадению с key, возвращает первый в списке найденных, если совпадений больше одного.
		FindFile(ctx context.Context, key string) (*StorageObject, error)

		// GetMetadata возвращает метаданные объекта по его ключу.
		GetMetadata(ctx context.Context, key string) (*StorageObject, error)

		// UploadFile загрузка файла в s3.
		UploadFile(ctx context.Context, file FileObject) (*StorageObject, error)

		// UploadStart инициализирует multipart загрузку.
		UploadStart(ctx context.Context, file FileObject) (string, error)

		// UploadPart загрузка файла в multipart операции.
		// Вызывается после инициализации multipart загрузки.
		UploadPart(ctx context.Context, uploadID string, file FileObject) (*StorageObject, error)

		// UploadComplete вызывается для объединения всех загруженных компонент и создания объекта в s3.
		// Вызывается после инициализации multipart загрузки.
		UploadComplete(ctx context.Context, name, uploadID string, parts []*StorageObject) (*StorageObject, error)

		// UploadCancel прерывает multipart загрузку.
		UploadCancel(ctx context.Context, uploadID, name string) error

		// CopyFile создает новый объект по текущему с новым key и метаданными.
		CopyFile(ctx context.Context, destKey, srcKey string, metadata map[string]string) error

		// DeleteFile удаляет объект из хранилища по его ключу.
		DeleteFile(ctx context.Context, key string) error

		// FormationURL возвращает ссылку на файл.
		FormationURL(fileName string) *url.URL
	}

	// Uploader интерфейс для загрузки данных на S3.
	Uploader interface {
		// Write записывает часть данных (chunk) для загрузки на S3.
		Write(chunk []byte) (int, error)
		// Cancel отменяет текущую загрузку.
		Cancel() error
		// Complete завершает загрузку данных и возвращает ссылку на объект хранилища S3.
		Complete() (*StorageObject, error)
	}

	// Downloader интерфейс для скачивания данных с S3.
	Downloader interface {
		// Content возвращает объект типа io.ReadSeeker для чтения скачанного контента.
		Content() io.ReadSeeker
		// Load загружает данные из хранилища S3 по ключу (key).
		Load(key string) error
	}

	// ZipManager интерфейс для работы с ZIP-архивами с загрузкой в S3.
	ZipManager interface {
		// Name возвращает имя текущего ZIP-архива.
		Name() string
		// AddFromStorage добавляет файл из хранилища S3 в текущий архив.
		AddFromStorage(name, key string) error
		// Add добавляет файл в текущий архив, используя функцию writeTo для записи содержимого файла.
		Add(name string, writeTo func(io.Writer) error) error
		// AddDir добавляет пустую директорию в текущий архив.
		AddDir(dirName string) error
		// Complete завершает работу с архивом и возвращает ссылку на объект хранилища S3.
		Complete() (*StorageObject, error)
		// TotalUploaded возвращает общий размер загруженных данных в архив.
		TotalUploaded() size.Size
		// Cancel отменяет работу с текущим архивом.
		Cancel()
	}
)

func (c *Client) GetMetadata(ctx context.Context, key string) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("requesting object metadata '%s' (bucket '%s')", key, c.config.BucketName)

	response, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.config.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get object metadata: %w", err)
	}

	logger.Trace().Msgf("got object '%s' metadata (bucket '%s'): %v", key, c.config.BucketName, response.Metadata)

	return &StorageObject{
		Key:      key,
		ETag:     aws.ToString(response.ETag),
		Metadata: response.Metadata,
		Expires:  response.Expires,
	}, nil
}

func (c *Client) ListFiles(ctx context.Context, prefix string) ([]*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("listing objects %s (bucket '%s')", prefix, c.config.BucketName)

	objects := make([]*StorageObject, 0, DefaultMaxKeys)

	var continuationToken *string
	for fetchNext := true; fetchNext; {
		listObjects, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(c.config.BucketName),
			Prefix:            aws.String(prefix),
			MaxKeys:           aws.Int32(DefaultMaxKeys),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, fmt.Errorf("can't list objects with prefix %s in bucket %s: %w",
				prefix, c.config.BucketName, err)
		}

		fetchNext = aws.ToBool(listObjects.IsTruncated)
		continuationToken = listObjects.NextContinuationToken

		objects = append(objects, collections.Map(listObjects.Contents, func(object types.Object) *StorageObject {
			return &StorageObject{
				Key:          aws.ToString(object.Key),
				ETag:         aws.ToString(object.ETag),
				Size:         size.Size(aws.ToInt64(object.Size)),
				LastModified: object.LastModified,
			}
		})...)
	}

	logger.Trace().Msgf("found %d objects with prefix %s (bucket '%s')", len(objects), prefix, c.config.BucketName)

	return objects, nil
}

func (c *Client) FindFile(ctx context.Context, key string) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("trying to find object '%s' (bucket '%s')", key, c.config.BucketName)

	listObjects, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(c.config.BucketName),
		Prefix:  aws.String(key),
		MaxKeys: aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list bucket: %w", err)
	}

	if aws.ToInt32(listObjects.KeyCount) != 1 ||
		listObjects.Contents[0].Key == nil ||
		*listObjects.Contents[0].Key != key {
		logger.Trace().Msgf("not found object in storage '%s' (bucket '%s')", key, c.config.BucketName)

		return (*StorageObject)(nil), nil
	}

	object := listObjects.Contents[0]

	logger.Trace().Msgf("found object in storage '%s' (bucket '%s')", key, c.config.BucketName)

	return &StorageObject{
		Key:          aws.ToString(object.Key),
		ETag:         aws.ToString(object.ETag),
		Size:         size.Size(aws.ToInt64(object.Size)),
		LastModified: object.LastModified,
	}, nil
}

func (c *Client) GetListFiles(ctx context.Context, keys []string) ([]*StorageObject, error) {
	return collections.MapWithErr(keys, func(key string) (*StorageObject, error) { return c.GetFile(ctx, key) })
}

func (c *Client) CopyFile(ctx context.Context, destKey, srcKey string, metadata map[string]string) error {
	logger := c.logger(ctx)
	logger.Trace().Msgf("copying object '%s' to '%s' (bucket '%s')", srcKey, destKey, c.config.BucketName)

	out, err := c.client.CopyObject(ctx, &s3.CopyObjectInput{
		CopySource: aws.String(path.Join(c.config.BucketName, srcKey)),
		Bucket:     aws.String(c.config.BucketName),
		Key:        aws.String(destKey),
		Metadata:   metadata,
	})
	if err != nil {
		return fmt.Errorf("unable to copy file: %w", err)
	}

	logger.Trace().Msgf("object '%s' copied to '%s' (bucket '%s'): ETag: %s",
		srcKey, destKey, c.config.BucketName, aws.ToString(out.CopyObjectResult.ETag))

	return nil
}

func (c *Client) GetFile(ctx context.Context, key string) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("requesting object '%s' (bucket '%s')", key, c.config.BucketName)

	object, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.config.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get file from bucket: %w", err)
	}

	etag := ""
	if object.ETag != nil {
		etag = *object.ETag
	}

	logger.Trace().Msgf("got object '%s' (bucket '%s'): ETag %s", key, c.config.BucketName, etag)

	return &StorageObject{
		Key:          key,
		ETag:         aws.ToString(object.ETag),
		LastModified: object.LastModified,
		Metadata:     object.Metadata,
		Size:         size.Size(aws.ToInt64(object.ContentLength)),
		Body:         object.Body,
		Expires:      object.Expires,
	}, nil
}

func (c *Client) UploadFile(ctx context.Context, file FileObject) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("uploading object '%s' (bucket '%s')", file.Name, c.config.BucketName)

	object, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Key:         aws.String(file.Name),
		Bucket:      aws.String(c.config.BucketName),
		Body:        file.Body,
		Metadata:    file.Metadata,
		ContentType: file.ContentType,
		Expires:     file.Expires,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to upload file to bucket: %w", err)
	}

	logger.Trace().Msgf("uploaded object '%s' (bucket '%s')", file.Name, c.config.BucketName)

	return &StorageObject{
		Key:     file.Name,
		ETag:    aws.ToString(object.ETag),
		Expires: file.Expires,
	}, nil
}

func (c *Client) UploadStart(ctx context.Context, file FileObject) (string, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("starting object '%s' upload (bucket '%s')", file.Name, c.config.BucketName)

	out, err := c.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Key:         aws.String(file.Name),
		Bucket:      aws.String(c.config.BucketName),
		Metadata:    file.Metadata,
		ContentType: file.ContentType,
		Expires:     file.Expires,
	})
	if err != nil {
		return "", fmt.Errorf("unable to start object '%s' upload (bucket '%s'): %w", file.Name, c.config.BucketName, err)
	}

	logger.Trace().Msgf("started object '%s' upload (bucket '%s'): %s", file.Name, c.config.BucketName, aws.ToString(out.UploadId))

	return aws.ToString(out.UploadId), nil
}

func (c *Client) UploadPart(ctx context.Context, uploadID string, file FileObject) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("uploading object '%s' part (bucket '%s')", file.Name, c.config.BucketName)

	out, err := c.client.UploadPart(ctx, &s3.UploadPartInput{
		Key:        aws.String(file.Name),
		Bucket:     aws.String(c.config.BucketName),
		UploadId:   aws.String(uploadID),
		Body:       file.Body,
		PartNumber: aws.Int32(file.PartNumber),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to upload object '%s' part (bucket '%s'): %w", file.Name, c.config.BucketName, err)
	}

	logger.Trace().Msgf("uploaded object '%s' part (bucket '%s'): ETag %s", file.Name, c.config.BucketName, aws.ToString(out.ETag))

	return &StorageObject{
		Key:        file.Name,
		ETag:       aws.ToString(out.ETag),
		PartNumber: file.PartNumber,
	}, nil
}

func (c *Client) UploadComplete(ctx context.Context, uploadID, name string, parts []*StorageObject) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("completing object '%s' upload (bucket '%s')", name, c.config.BucketName)

	out, err := c.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Key:      aws.String(name),
		Bucket:   aws.String(c.config.BucketName),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: collections.Map(parts, func(part *StorageObject) types.CompletedPart {
				return types.CompletedPart{
					ETag:       aws.String(part.ETag),
					PartNumber: aws.Int32(part.PartNumber),
				}
			}),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to complete object '%s' upload (bucket '%s'): %w", name, c.config.BucketName, err)
	}

	logger.Trace().Msgf("completed object '%s' upload (bucket '%s'): ETag %s", name, c.config.BucketName, aws.ToString(out.ETag))

	return &StorageObject{
		Key:  aws.ToString(out.Key),
		ETag: aws.ToString(out.ETag),
	}, nil
}

func (c *Client) UploadCancel(ctx context.Context, uploadID, name string) error {
	logger := c.logger(ctx)
	logger.Trace().Msgf("aborting object '%s' upload (bucket '%s')", name, c.config.BucketName)

	_, err := c.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Key:      aws.String(name),
		Bucket:   aws.String(c.config.BucketName),
		UploadId: aws.String(uploadID),
	})
	if err != nil {
		return fmt.Errorf("unable to abort object '%s' upload (bucket '%s'): %w", name, c.config.BucketName, err)
	}

	logger.Trace().Msgf("aborted object '%s' upload (bucket '%s')", name, c.config.BucketName)

	return nil
}

func (c *Client) DeleteFile(ctx context.Context, key string) error {
	logger := c.logger(ctx)
	logger.Trace().Msgf("deleting object '%s' (bucket '%s')", key, c.config.BucketName)

	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(c.config.BucketName),
	})
	if err != nil {
		return fmt.Errorf("can't delete object %s from bucket: %s: %w", key, c.config.BucketName, err)
	}

	logger.Trace().Msgf("deleted object '%s' (bucket '%s')", key, c.config.BucketName)

	return nil
}

func (c *Client) FormationURL(fileName string) *url.URL {
	return c.urlProvider(fileName)
}

func (c *Client) logger(ctx context.Context) *zerolog.Logger {
	return logs.FromContext(ctx)
}

func (o FileObject) WithExpiration(expires time.Time) FileObject {
	o.Expires = &expires

	return o
}

func makeURLProvider(baseURL *url.URL, bucketName string) func(string) *url.URL {
	return func(fileName string) *url.URL {
		fileURL := *baseURL
		fileURL.Path = path.Join(baseURL.Path, bucketName, fileName)

		return &fileURL
	}
}
