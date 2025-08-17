package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/collections"
)

var _ ClientInterface = (*ClientTestDouble)(nil)

const (
	filePermission      = 0o600
	fileWritePermission = 0o644
	dirWritePermission  = 0o755
)

type (
	// ClientTestDouble тестовый двойник, который имитирует работу s3, храня все в файловой системе.
	// Все файлы создаются в отдельной тестовой директории.
	// Не забывайте очищать после работы директории через команды RemoveDir/ClearDir.
	ClientTestDouble struct {
		dir    string
		config *Config
	}
)

func NewTestDouble(config *Config) (*ClientTestDouble, error) {
	dir, err := os.MkdirTemp("./", config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("can't make temp dir %s: %w", config.BucketName, err)
	}

	return &ClientTestDouble{
		dir:    dir,
		config: config,
	}, nil
}

func (c *ClientTestDouble) RemoveDir() error {
	if err := os.RemoveAll(c.dir); err != nil {
		return fmt.Errorf("can't remove dir %s: %w", c.dir, err)
	}

	return nil
}

func (c *ClientTestDouble) ClearDir() error {
	err := c.RemoveDir()
	if err != nil {
		return err
	}

	dir, err := os.MkdirTemp("./", c.config.BucketName)
	if err != nil {
		return fmt.Errorf("can't make temp dir %s: %w", c.config.BucketName, err)
	}

	c.dir = dir

	return nil
}

func (c *ClientTestDouble) UploadFile(ctx context.Context, file FileObject) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("uploading object '%s' (bucket '%s')", file.Name, c.config.BucketName)

	content, _ := io.ReadAll(file.Body)
	if strings.Contains(file.Name, "/") {
		folders := strings.Split(file.Name, "/")
		folders = folders[:len(folders)-1]
		err := os.MkdirAll(filepath.Join(c.dir, strings.Join(folders, "/")), dirWritePermission)
		if err != nil {
			return nil, fmt.Errorf("unable to create folder: %w", err)
		}
	}

	err := os.WriteFile(filepath.Join(c.dir, file.Name), content, fileWritePermission)
	if err != nil {
		return nil, fmt.Errorf("unable to upload file to bucket: %w", err)
	}

	logger.Trace().Msgf("uploaded object '%s' (bucket '%s')", file.Name, c.config.BucketName)

	return &StorageObject{
		Key: file.Name,
	}, nil
}

func (c *ClientTestDouble) GetFile(ctx context.Context, key string) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("requesting object '%s' (bucket '%s')", key, c.config.BucketName)
	body, err := os.ReadFile(filepath.Join(c.dir, key))
	if err != nil {
		return nil, fmt.Errorf("unable to upload file to bucket: %w", err)
	}

	return &StorageObject{
		Key: key,
		Body: io.NopCloser(
			bytes.NewReader(body),
		),
	}, nil
}

func (c *ClientTestDouble) DeleteFile(ctx context.Context, key string) error {
	logger := c.logger(ctx)
	logger.Trace().Msgf("deleting object '%s' (bucket '%s')", key, c.config.BucketName)

	filePath := filepath.Join(c.dir, key)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("can't delete file %s: %w", filePath, err)
	}

	return nil
}

func (c *ClientTestDouble) GetListFiles(ctx context.Context, keys []string) ([]*StorageObject, error) {
	return collections.MapWithErr(keys, func(key string) (*StorageObject, error) { return c.GetFile(ctx, key) })
}

func (c *ClientTestDouble) FindFile(ctx context.Context, key string) (*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("trying to find object '%s' (bucket '%s')", key, c.config.BucketName)

	files, err := os.ReadDir(c.dir)
	if err != nil {
		return nil, fmt.Errorf("can't read dir %s: %w", c.dir, err)
	}

	if len(files) == 0 {
		logger.Trace().Msgf("not found object in storage '%s' (bucket '%s')", key, c.config.BucketName)

		return (*StorageObject)(nil), nil
	}

	file, found := collections.Find(files, func(file os.DirEntry) bool {
		return file.Name() == key
	})

	if !found {
		return (*StorageObject)(nil), nil
	}

	return c.GetFile(ctx, file.Name())
}

func (c *ClientTestDouble) ListFiles(ctx context.Context, prefix string) ([]*StorageObject, error) {
	logger := c.logger(ctx)
	logger.Trace().Msgf("listing objects %s (bucket '%s')", prefix, c.config.BucketName)

	files, err := os.ReadDir(c.dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %w", c.dir, err)
	}

	fileFilteredByPrefix := collections.Filter(files, func(file os.DirEntry) bool {
		return strings.HasPrefix(file.Name(), prefix)
	})

	logger.Trace().Msgf("found %d objects with prefix %s (bucket '%s')", len(fileFilteredByPrefix), prefix, c.config.BucketName)

	return collections.Map(fileFilteredByPrefix, func(file os.DirEntry) *StorageObject {
		return &StorageObject{
			Key: file.Name(),
		}
	}), nil
}

func (c *ClientTestDouble) FormationURL(fileName string) *url.URL {
	return &url.URL{
		Host: "localhost",
		Path: fileName,
	}
}

func (c *ClientTestDouble) GetMetadata(_ context.Context, _ string) (*StorageObject, error) {
	return &StorageObject{}, nil
}

func (c *ClientTestDouble) CopyFile(ctx context.Context, destKey, srcKey string, _ map[string]string) error {
	logger := c.logger(ctx)
	logger.Trace().Msgf("copying object '%s' to '%s' (bucket '%s')", srcKey, destKey, c.config.BucketName)

	file, err := c.GetFile(ctx, destKey)
	if err != nil {
		return err
	}

	_, err = c.UploadFile(ctx, FileObject{
		Name: srcKey,
		Body: file.Body,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientTestDouble) UploadStart(context.Context, FileObject) (string, error) {
	return "", nil
}

func (c *ClientTestDouble) UploadPart(context.Context, string, FileObject) (*StorageObject, error) {
	return &StorageObject{}, nil
}

func (c *ClientTestDouble) UploadComplete(context.Context, string, string, []*StorageObject) (*StorageObject, error) {
	return &StorageObject{}, nil
}

func (c *ClientTestDouble) UploadCancel(context.Context, string, string) error {
	return nil
}

func (c *ClientTestDouble) logger(ctx context.Context) *zerolog.Logger {
	return logs.FromContext(ctx)
}
