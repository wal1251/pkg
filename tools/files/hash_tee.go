package files

import (
	"crypto/md5" //nolint:gosec
	"hash"
	"io"
)

var _ io.Reader = (*HashTeeReader)(nil)

// HashTeeReader обертка интерфейса io.Reader, при чтении из интерфейса подсчитывается хеш сумма. Хеш функция задается
// при инициализации.
type HashTeeReader struct {
	io.Reader
	hash hash.Hash
}

// Hash возвращает хэш сумму.
func (h *HashTeeReader) Hash() []byte {
	return h.hash.Sum(make([]byte, 0, h.hash.Size()))
}

func NewHashTeeReader(src io.Reader, algorithm func() hash.Hash) *HashTeeReader {
	h := algorithm()

	return &HashTeeReader{
		Reader: io.TeeReader(src, h),
		hash:   h,
	}
}

func NewMD5TeeReader(src io.Reader) *HashTeeReader {
	return NewHashTeeReader(src, md5.New)
}
