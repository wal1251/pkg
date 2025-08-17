package size

import (
	"io"
)

var _ io.Writer = (*CountingWriter)(nil)

type (
	// CountingWriter позволяет отслеживать общее количество записанных байтов.
	CountingWriter struct {
		Written int64
	}
)

// Write метод реализует интерфейс io.Writer и позволяет отслеживать общее количество записанных байтов.
// Возвращает количество записанных байтов и ошибку, если она возникла.
func (w *CountingWriter) Write(p []byte) (int, error) {
	l := len(p)
	w.Written += int64(l)

	return l, nil
}
