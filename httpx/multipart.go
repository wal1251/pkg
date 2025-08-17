package httpx

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/wal1251/pkg/core/logs"
)

const sniffSize = 512

type MultipartFile struct {
	logger      *zerolog.Logger
	file        multipart.File
	header      *multipart.FileHeader
	contentType string
}

func (f *MultipartFile) Content() io.ReadSeekCloser {
	return f.file
}

func (f *MultipartFile) Close() {
	if err := f.file.Close(); err != nil {
		f.logger.Err(err).Msg("failed to close form multipart file content")
	}
}

func (f *MultipartFile) ReadAndClose(read func(reader io.Reader) error) error {
	defer f.Close()

	return read(f.file)
}

func (f *MultipartFile) Size() int64 {
	return f.header.Size
}

func (f *MultipartFile) Filename() string {
	return f.header.Filename
}

func (f *MultipartFile) ContentType() string {
	return f.contentType
}

func NewMultipartFile(request *http.Request, key string) (*MultipartFile, error) {
	file, header, err := request.FormFile(key)
	if err != nil {
		return nil, fmt.Errorf("unable to read form file: %w", err)
	}

	bytesToSniff := make([]byte, sniffSize)
	if _, err = file.Read(bytesToSniff); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to read multipart file: %w", err)
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek multipart file: %w", err)
	}

	return &MultipartFile{
		logger:      logs.FromContext(request.Context()),
		contentType: http.DetectContentType(bytesToSniff),
		file:        file,
		header:      header,
	}, nil
}
