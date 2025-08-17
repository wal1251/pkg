package docx_test

import (
	"archive/zip"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/docx"
)

func TestReadDocFromMemory(t *testing.T) {
	rawExampleContent, err := os.ReadFile("testdata/raw_example.xml")
	require.NoError(t, err)

	tests := []struct {
		name        string
		filePath    string
		expectedErr error
	}{
		{
			name:        "Success",
			filePath:    "testdata/example.docx",
			expectedErr: nil,
		},
		{
			name:        "File not found",
			filePath:    "testdata/not_found.docx",
			expectedErr: os.ErrNotExist,
		},
		{
			name:        "No docx file",
			filePath:    "testdata/example.txt",
			expectedErr: zip.ErrFormat,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			content, err := readDocWrapper(tc.filePath)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, string(rawExampleContent), content)
		})
	}
}

func readDocWrapper(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	size := fileInfo.Size()

	return docx.ReadDocFromMemory(file, size)
}
