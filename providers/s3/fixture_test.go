package s3_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	s3 "github.com/wal1251/pkg/providers/s3"
	"github.com/wal1251/pkg/tools/collections"
)

func TestInitFixture(t *testing.T) {
	fileNameList := []string{
		"fileName1",
		"fileName2",
		"fileName3",
	}

	for _, fileName := range fileNameList {
		err := os.WriteFile(fileName, []byte("content"), os.FileMode(0o600))
		require.NoError(t, err)
	}

	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)

	err = client.InitFilesFixture(fileNameList)
	require.NoError(t, err)

	for _, fileName := range fileNameList {
		err = os.Remove(fileName)
		require.NoError(t, err)
	}

	files, err := client.ListFiles(context.TODO(), "fileName")
	require.NoError(t, err)

	require.ElementsMatch(t, fileNameList, collections.Map(files, func(o *s3.StorageObject) string {
		return o.Key
	}))

	err = client.RemoveDir()
	require.NoError(t, err)
}
