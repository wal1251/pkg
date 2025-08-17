// Package ftp содержит компоненты для работы с FTP сервером.
package ftp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/textproto"

	"github.com/jlaffaye/ftp"
)

const (
	ComponentName   = "ftp" // Имя компонента.
	ErrCodeNotFound = 550   // Код ошибки FTP NOT FOUND.
)

type (
	// Client предоставляем методы для работы с файлами на FTP сервере.
	Client interface {
		// Dial устанавливает соединение с сервером.
		Dial(context.Context) error

		// Quit закрывает открытое соединение с сервером, см. Client.Dial().
		Quit(context.Context)

		// Size получает размер указанного файла на сервере.
		Size(ctx context.Context, filePath string) (int64, error)

		// Walk рекурсивно обходит содержимое указанного каталога на сервере.
		Walk(ctx context.Context, folderPath string, accept Visitor) error

		// List вернет содержимое указанного каталога на сервере (не рекурсивно).
		List(ctx context.Context, folderPath string) ([]*Entry, error)

		// NameList вернет имена файлов в указанном каталоге на сервере (не рекурсивно).
		NameList(ctx context.Context, folderPath string) ([]string, error)

		// Download вернет io.ReadCloser указанного файла на сервере для скачивания его содержимого.
		Download(ctx context.Context, filePath string) (io.ReadCloser, error)
	}

	// Visitor callback функция для оповещения клиента о найденном содержимом в каталоге сервера. Если не требуется
	// рекурсивный обход содержимого текущего каталога, в теле функции необходимо вызвать skip().
	Visitor func(entry *Entry, skip func()) error
)

// Connect возвращает аутентифицированное соединение с FTP сервером.
func Connect(ctx context.Context, cfg *Config) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(cfg.Address,
		ftp.DialWithTimeout(cfg.Timeout),
		ftp.DialWithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("can't dial server: %w", err)
	}

	err = conn.Login(cfg.User, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("can't login server: %w", err)
	}

	return conn, nil
}

// IsNotFound вернет true, если указанная ошибка является ошибкой FTP NOT FOUND (550).
func IsNotFound(err error) bool {
	if err != nil {
		var ftpErr *textproto.Error
		if ok := errors.As(err, &ftpErr); ok && ftpErr.Code == ErrCodeNotFound {
			return true
		}
	}

	return false
}
