package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/clock"
	"github.com/wal1251/pkg/tools/size"
)

const (
	DefaultRetryCount   = 3
	DefaultRetryTimeout = time.Second
)

type Loader struct {
	client ClientInterface
	buf    *bytes.Buffer
}

// Content возвращает объект типа io.ReadSeeker для чтения скачанных данных из буфера.
func (l *Loader) Content() io.ReadSeeker {
	return bytes.NewReader(l.buf.Bytes())
}

// Load загружает данные из хранилища S3 по ключу (key) и сохраняет их в буфере.
func (l *Loader) Load(ctx context.Context, key string) error {
	return CopyToBuffer(ctx, l.client, key, l.buf)
}

// NewLoader создает новый экземпляр Loader с указанным клиентом S3 и размером буфера.
func NewLoader(client ClientInterface, bufSize size.Size) *Loader {
	return &Loader{
		client: client,
		buf:    bytes.NewBuffer(make([]byte, 0, bufSize.Int())),
	}
}

// Download выполняет загрузку файла с S3 по заданному имени (name) и передает его в функцию accept.
func Download(ctx context.Context, client ClientInterface, name string, accept func(*StorageObject) (bool, error)) error {
	logger := logs.FromContext(ctx)
	retry := clock.RetryingWrapperWE(DefaultRetryCount, DefaultRetryTimeout)

	return retry(func() (bool, error) {
		obj, err := client.GetFile(ctx, name)
		if err != nil {
			return true, err
		}
		defer func() {
			if err := obj.Body.Close(); err != nil {
				logger.Err(err).Msgf("can't close downloaded storage object content: %s", obj.Key)
			}
		}()

		return accept(obj)
	})
}

// CopyToBuffer загружает файл из хранилища S3 по имени (name) и копирует его содержимое в буфер (buf).
func CopyToBuffer(ctx context.Context, client ClientInterface, name string, buf *bytes.Buffer) error {
	logger := logs.FromContext(ctx)

	return Download(ctx, client, name, func(object *StorageObject) (bool, error) {
		buf.Reset()
		if _, err := io.Copy(buf, object.Body); err != nil {
			return true, fmt.Errorf("failed to copy downloaded storage object %s content to buffer: %w",
				object.Key, err)
		}
		logger.Debug().Msgf("storage object copied to buffer (%v): %s", object.Size, object.Key)

		return false, nil
	})
}
