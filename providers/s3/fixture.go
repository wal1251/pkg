package s3

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

// InitFilesFixture инициализирует фикстуры для работы с s3 двойником.
func (c *ClientTestDouble) InitFilesFixture(filePaths []string) error {
	for _, filePath := range filePaths {
		bytesRead, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read fixtures file %s: %w", filePath, err)
		}

		parts := strings.Split(filePath, "/")
		_, err = c.UploadFile(context.TODO(), FileObject{
			Name: parts[len(parts)-1],
			Body: bytes.NewReader(bytesRead),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
