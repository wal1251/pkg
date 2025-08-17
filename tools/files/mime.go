package files

import (
	"path/filepath"
	"strings"

	gomime "github.com/cubewise-code/go-mime"
)

// MimeTypeByExtension определяет тип содержимого (mime type) файла по расширению.
func MimeTypeByExtension(fileName string) string {
	return gomime.TypeByExtension(filepath.Ext(fileName))
}

// MimeTypeIsImage является ли тип (mime type) изображением.
func MimeTypeIsImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}
