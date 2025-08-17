package files_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/files"
)

func TestHashTeeReader(t *testing.T) {
	tests := []struct {
		name     string
		sample   string
		wantHash []byte
	}{
		{
			name:     "Basic case",
			sample:   "foo bar baz",
			wantHash: []byte{0xab, 0x7, 0xac, 0xbb, 0x1e, 0x49, 0x68, 0x1, 0x93, 0x7a, 0xdf, 0xa7, 0x72, 0x42, 0x4b, 0xf7},
		},
		{
			name:     "Nothing to read",
			sample:   "",
			wantHash: []byte{0xd4, 0x1d, 0x8c, 0xd9, 0x8f, 0x0, 0xb2, 0x4, 0xe9, 0x80, 0x9, 0x98, 0xec, 0xf8, 0x42, 0x7e},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := bytes.NewBufferString(tt.sample)
			wantN := int64(src.Len())
			reader := files.NewMD5TeeReader(src)
			dest := bytes.NewBuffer(make([]byte, 0, src.Len()))

			n, err := io.Copy(dest, reader)
			if assert.NoError(t, err) {
				assert.Equalf(t, wantN, n, "copied quantity of bytes doesn't match")
				assert.Equalf(t, tt.sample, dest.String(), "destination doesn't match its source")
				assert.Equalf(t, tt.wantHash, reader.Hash(), "wantHash doesn't match")
			}
		})
	}
}
