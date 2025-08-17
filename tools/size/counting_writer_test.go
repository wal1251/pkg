package size_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/size"
)

func TestCountingWriter_Write(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantErr     bool
		wantWritten int64
	}{
		{
			name:        "Empty input",
			input:       []byte(""),
			wantErr:     false,
			wantWritten: 0,
		},
		{
			name:        "Non-empty input - 5 bytes",
			input:       []byte("hello"),
			wantErr:     false,
			wantWritten: 5,
		},
		{
			name:        "Non-empty input - 11 bytes",
			input:       []byte("hello world"),
			wantErr:     false,
			wantWritten: 11,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			writer := size.CountingWriter{}

			_, err := writer.Write(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			}

			assert.Equal(t, tc.wantWritten, writer.Written)
		})
	}
}
