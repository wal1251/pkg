package s3

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/size"
)

const (
	dirDelimiter        = "/"
	defaultBufferFactor = 1.75
)

type ZipCopy struct {
	name     string
	arch     *zip.Writer
	method   uint16
	loader   *Loader
	uploader *MultipartUploader
}

// Name возвращает название текущего ZIP-архива.
func (z *ZipCopy) Name() string {
	return z.name
}

// AddFromStorage добавляет файл из хранилища S3 в текущий архив.
func (z *ZipCopy) AddFromStorage(ctx context.Context, name, key string) error {
	if err := z.loader.Load(ctx, key); err != nil {
		return err
	}

	return z.Add(name, func(w io.Writer) error {
		_, err := io.Copy(w, z.loader.Content())

		return fmt.Errorf("failed to copy content: %w", err)
	})
}

// Add добавляет файл в текущий архив, используя функцию writeTo для записи содержимого файла.
func (z *ZipCopy) Add(name string, writeTo func(io.Writer) error) error {
	w, err := z.arch.CreateHeader(&zip.FileHeader{Name: name, Method: z.method})
	if err != nil {
		return fmt.Errorf("failed to created header: %w", err)
	}
	if err = writeTo(w); err != nil {
		return err
	}
	if err = z.arch.Flush(); err != nil {
		return fmt.Errorf("failed to flush zip: %w", err)
	}

	return nil
}

// AddDir добавляет пустую директорию в текущий архив.
func (z *ZipCopy) AddDir(dirName string) error {
	if !strings.HasSuffix(dirName, dirDelimiter) {
		dirName += dirDelimiter
	}
	if _, err := z.arch.CreateHeader(&zip.FileHeader{Name: dirName, Method: z.method}); err != nil {
		return fmt.Errorf("failed to created header: %w", err)
	}

	return nil
}

// Complete завершает работу с архивом и возвращает ссылку на объект хранилища S3.
func (z *ZipCopy) Complete() (*StorageObject, error) {
	if err := z.arch.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush zip: %w", err)
	}
	if err := z.arch.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip: %w", err)
	}
	object, err := z.uploader.Complete()
	if err != nil {
		return nil, err
	}

	return object, nil
}

// TotalUploaded возвращает общий размер загруженных данных в архив.
func (z *ZipCopy) TotalUploaded() size.Size {
	return size.Size(z.uploader.TotalUploaded)
}

// Cancel закрывает текущий архив и отменяет загрузку.
func (z *ZipCopy) Cancel(ctx context.Context) {
	logger := logs.FromContext(ctx)
	defer func() {
		if err := z.arch.Close(); err != nil {
			logger.Err(err).Msgf("failed to close archive")
		}
	}()
	if err := z.uploader.Cancel(); err != nil {
		logger.Err(err).Msgf("failed to cancel file upload")
	}
}

// NewZipCopy инициирует multipart загрузку, создает и возвращает новый экземпляр ZipCopy.
func NewZipCopy(ctx context.Context, client ClientInterface, file FileObject, bufferSize size.Size) (*ZipCopy, error) {
	upload, err := StartMultipartUpload(ctx, client, file, size.Size(defaultBufferFactor*float64(bufferSize)))
	if err != nil {
		return nil, err
	}

	return &ZipCopy{
		name:     file.Name,
		arch:     zip.NewWriter(upload),
		uploader: upload,
		loader:   NewLoader(client, bufferSize),
		method:   zip.Store,
	}, nil
}
