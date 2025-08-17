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

var _ io.Writer = (*MultipartUploader)(nil)

type MultipartUploader struct {
	name           string
	uploadID       string
	parts          []*StorageObject
	partNumber     int32
	partSize       int
	buf            *bytes.Buffer
	ctx            context.Context // nolint:containedctx
	client         ClientInterface
	retriesCount   int
	retriesTimeout time.Duration
	TotalUploaded  int64
}

// Complete завершает загрузку файла на хранилище S3.
// Возвращает ссылку на загруженный объект и ошибку, если возникла проблема во время загрузки.
func (w *MultipartUploader) Complete() (*StorageObject, error) {
	if err := w.flush(); err != nil {
		return nil, err
	}

	return w.client.UploadComplete(w.ctx, w.uploadID, w.name, w.parts)
}

// Cancel отменяет текущую многопоточную загрузку файла на хранилище S3.
// Возвращает ошибку, если возникла проблема при отмене загрузки.
func (w *MultipartUploader) Cancel() error {
	return w.client.UploadCancel(w.ctx, w.uploadID, w.name)
}

// Write записывает данные из среза chunk в буфер MultipartUploader.
// Если размер данных в буфере превышает размер части, производится загрузка в хранилище S3.
// Возвращает общее количество записанных байт и ошибку, если возникла проблема во время записи.
func (w *MultipartUploader) Write(chunk []byte) (int, error) {
	writtenTotal := 0
	for chunkSize := len(chunk); w.buf.Len()+chunkSize > w.partSize; {
		bytesRemain := w.partSize - w.buf.Len()
		if chunkSize >= bytesRemain {
			bytesWritten, err := w.buf.Write(chunk[:bytesRemain])
			if err != nil {
				return 0, fmt.Errorf("failed to write to buffer: %w", err)
			}
			writtenTotal += bytesWritten
			chunk = chunk[bytesWritten:]
			chunkSize -= bytesWritten
		}
		if err := w.flush(); err != nil {
			return 0, err
		}
	}
	if len(chunk) == 0 {
		return writtenTotal, nil
	}

	n, err := w.buf.Write(chunk)
	if err != nil {
		return 0, fmt.Errorf("failed to write to buffer: %w", err)
	}
	writtenTotal += n

	return writtenTotal, nil
}

// flush загружает данные из буфера MultipartUploader на хранилище S3 в виде отдельной части (в multipart операции).
// Если буфер пуст, функция ничего не делает и возвращает nil.
// В случае ошибки при загрузке части, функция пытается выполнить повторную загрузку согласно параметрам retriesCount и retriesTimeout.
// Возвращает ошибку, если возникла проблема во время загрузки части.
func (w *MultipartUploader) flush() error {
	bufSize := size.Size(w.buf.Len())
	if bufSize == 0 {
		return nil
	}
	w.partNumber++
	w.TotalUploaded += int64(bufSize)

	logger := logs.FromContext(w.ctx)
	logger.Debug().Msgf("flushing output buffer (%v): %s", bufSize, w.name)

	retry := clock.RetryingWrapperWE(w.retriesCount, w.retriesTimeout)

	return retry(func() (bool, error) {
		part, err := w.client.UploadPart(w.ctx, w.uploadID, FileObject{
			Name:       w.name,
			PartNumber: w.partNumber,
			Body:       bytes.NewReader(w.buf.Bytes()),
		})
		if err != nil {
			return true, err
		}
		w.buf.Reset()
		w.parts = append(w.parts, part)

		logger.Debug().Msgf("uploaded part %d (%v): %s", w.partNumber, bufSize, w.name)

		return false, nil
	})
}

// StartMultipartUpload инициирует multipart загрузку файла на хранилище S3.
func StartMultipartUpload(ctx context.Context, client ClientInterface, file FileObject, size size.Size) (*MultipartUploader, error) {
	id, err := client.UploadStart(ctx, file)
	if err != nil {
		return nil, err
	}

	return &MultipartUploader{
		buf:            bytes.NewBuffer(make([]byte, 0, size.Int())),
		partSize:       size.Int(),
		name:           file.Name,
		parts:          make([]*StorageObject, 0, 1),
		ctx:            ctx,
		client:         client,
		uploadID:       id,
		retriesCount:   DefaultRetryCount,
		retriesTimeout: DefaultRetryTimeout,
	}, nil
}
