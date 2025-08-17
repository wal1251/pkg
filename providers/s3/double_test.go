package s3_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/s3"
	"github.com/wal1251/pkg/tools/collections"
)

func TestUploadGetFile(t *testing.T) {
	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)
	ctx := context.TODO()

	contentFile := "test content"
	storageObject, err := client.UploadFile(ctx, s3.FileObject{
		Name: uuid.New().String(),
		Body: bytes.NewReader([]byte(contentFile)),
	})
	require.NoError(t, err)

	storageObject, err = client.GetFile(ctx, storageObject.Key)
	require.NoError(t, err)

	result, err := io.ReadAll(storageObject.Body)
	require.NoError(t, err)

	require.Equal(t, contentFile, string(result))

	err = client.RemoveDir()
	require.NoError(t, err)
}

func TestGetFiles(t *testing.T) {
	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)
	ctx := context.TODO()

	storageObjectList := make([]*s3.StorageObject, 0)
	for _, file := range []s3.FileObject{
		{
			Name: uuid.New().String(),
			Body: bytes.NewReader([]byte{}),
		},
		{
			Name: uuid.New().String(),
			Body: bytes.NewReader([]byte{}),
		},
		{
			Name: uuid.New().String(),
			Body: bytes.NewReader([]byte{}),
		},
		{
			Name: uuid.New().String(),
			Body: bytes.NewReader([]byte{}),
		},
	} {
		storageObject, err := client.UploadFile(ctx, file)
		require.NoError(t, err)

		storageObjectList = append(storageObjectList, storageObject)
	}

	getKeys := func(objects []*s3.StorageObject) []string {
		return collections.Map(storageObjectList, func(object *s3.StorageObject) string {
			return object.Key
		})
	}

	files, err := client.GetListFiles(
		ctx,
		getKeys(storageObjectList),
	)

	require.Equal(t, len(files), len(storageObjectList))
	require.ElementsMatch(t, getKeys(storageObjectList), getKeys(files))

	err = client.RemoveDir()
	require.NoError(t, err)
}

func TestDeleteFile(t *testing.T) {
	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)
	ctx := context.TODO()

	storageObject, err := client.UploadFile(ctx, s3.FileObject{
		Name: uuid.New().String(),
		Body: bytes.NewReader([]byte{}),
	})
	require.NoError(t, err)

	err = client.DeleteFile(ctx, storageObject.Key)
	require.NoError(t, err)

	_, err = client.GetFile(ctx, storageObject.Key)
	require.Error(t, err)

	err = client.RemoveDir()
	require.NoError(t, err)
}

func TestFindFile(t *testing.T) {
	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)
	ctx := context.TODO()

	content := "content"
	storageObject, err := client.UploadFile(ctx, s3.FileObject{
		Name: uuid.New().String(),
		Body: bytes.NewReader([]byte(content)),
	})
	require.NoError(t, err)

	result, err := client.FindFile(ctx, storageObject.Key)
	require.NoError(t, err)

	resultContent, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	require.Equal(t, content, string(resultContent))

	err = client.RemoveDir()
	require.NoError(t, err)
}

func TestFindFile_NotFile(t *testing.T) {
	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)
	ctx := context.TODO()

	result, err := client.FindFile(ctx, uuid.New().String())
	require.NoError(t, err)

	require.Nil(t, result)

	err = client.RemoveDir()
	require.NoError(t, err)
}

func TestListFiles(t *testing.T) {
	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)
	ctx := context.TODO()

	for _, file := range []s3.FileObject{
		{
			Name: "file_1",
			Body: bytes.NewReader([]byte{}),
		},
		{
			Name: "file_2",
			Body: bytes.NewReader([]byte{}),
		},
		{
			Name: "1_file_3",
			Body: bytes.NewReader([]byte{}),
		},
		{
			Name: "1_file_4",
			Body: bytes.NewReader([]byte{}),
		},
	} {
		_, err = client.UploadFile(ctx, file)
		require.NoError(t, err)
	}

	listFiles, err := client.ListFiles(ctx, "1_")
	require.NoError(t, err)

	require.ElementsMatch(t, []string{"1_file_3", "1_file_4"}, collections.Map(listFiles, func(o *s3.StorageObject) string {
		return o.Key
	}))

	err = client.RemoveDir()
	require.NoError(t, err)
}

func TestCopyFile(t *testing.T) {
	client, err := s3.NewTestDouble(&s3.Config{
		BucketName: "double_bucket",
	})
	require.NoError(t, err)
	ctx := context.TODO()

	content := "content"
	storageObject, err := client.UploadFile(ctx, s3.FileObject{
		Name: uuid.New().String(),
		Body: bytes.NewReader([]byte(content)),
	})
	require.NoError(t, err)

	newKey := fmt.Sprintf(
		"%s-%d-%d-%d",
		storageObject.Key,
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
	)
	err = client.CopyFile(
		ctx,
		storageObject.Key,
		newKey,
		nil,
	)
	require.NoError(t, err)

	file, err := client.GetFile(ctx, newKey)
	require.NoError(t, err)

	require.Equal(t, newKey, file.Key)

	resultContent, err := io.ReadAll(file.Body)
	require.NoError(t, err)

	require.Equal(t, content, string(resultContent))

	err = client.RemoveDir()
	require.NoError(t, err)
}
