package ftp_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/providers/ftp"
	"github.com/wal1251/pkg/tools/files"
)

func TestTestDoubleWalk(t *testing.T) {
	ctx := context.TODO()

	vol := files.NewVolume()
	require.NoError(t, vol.CreateTemp("samples*"), "can't create fixtures dir")
	defer func() {
		require.NoError(t, vol.Remove(), "can't cleanup fixtures dir")
	}()

	require.NoError(t, files.DirCreateEmptyFiles(vol.Path(),
		"foo/bar/baz/dummy.txt",
		"root/foo1/bar/file1.txt",
		"root/foo1/bar/file2.txt",
		"root/foo1/baz/dummy3.txt",
		"root/foo2/bar/dummy2.txt",
		"root/foo3/baz/dummy.txt",
		"root/foo3/dummy1.txt",
	), "failed to create fixture files")

	client := ftp.NewTestDouble(vol.Path())

	tests := []struct {
		name      string
		basePath  string
		wantPaths []string
		wantNames []string
	}{
		{
			name: "empty base path is root of home",
			wantPaths: []string{
				"/foo",
				"/foo/bar",
				"/foo/bar/baz",
				"/foo/bar/baz/dummy.txt",
				"/root",
				"/root/foo1",
				"/root/foo1/bar",
				"/root/foo1/bar/file1.txt",
				"/root/foo1/bar/file2.txt",
				"/root/foo1/baz",
				"/root/foo1/baz/dummy3.txt",
				"/root/foo2",
				"/root/foo2/bar",
				"/root/foo2/bar/dummy2.txt",
				"/root/foo3",
				"/root/foo3/baz",
				"/root/foo3/baz/dummy.txt",
				"/root/foo3/dummy1.txt",
			},
			wantNames: []string{
				"foo",
				"bar",
				"baz",
				"dummy.txt",
				"root",
				"foo1",
				"bar",
				"file1.txt",
				"file2.txt",
				"baz",
				"dummy3.txt",
				"foo2",
				"bar",
				"dummy2.txt",
				"foo3",
				"baz",
				"dummy.txt",
				"dummy1.txt",
			},
		},
		{
			name:     "/ base path is root of home",
			basePath: "/",
			wantPaths: []string{
				"/foo",
				"/foo/bar",
				"/foo/bar/baz",
				"/foo/bar/baz/dummy.txt",
				"/root",
				"/root/foo1",
				"/root/foo1/bar",
				"/root/foo1/bar/file1.txt",
				"/root/foo1/bar/file2.txt",
				"/root/foo1/baz",
				"/root/foo1/baz/dummy3.txt",
				"/root/foo2",
				"/root/foo2/bar",
				"/root/foo2/bar/dummy2.txt",
				"/root/foo3",
				"/root/foo3/baz",
				"/root/foo3/baz/dummy.txt",
				"/root/foo3/dummy1.txt",
			},
			wantNames: []string{
				"foo",
				"bar",
				"baz",
				"dummy.txt",
				"root",
				"foo1",
				"bar",
				"file1.txt",
				"file2.txt",
				"baz",
				"dummy3.txt",
				"foo2",
				"bar",
				"dummy2.txt",
				"foo3",
				"baz",
				"dummy.txt",
				"dummy1.txt",
			},
		},
		{
			name:     "subfolder walk",
			basePath: "root",
			wantPaths: []string{
				"/root/foo1",
				"/root/foo1/bar",
				"/root/foo1/bar/file1.txt",
				"/root/foo1/bar/file2.txt",
				"/root/foo1/baz",
				"/root/foo1/baz/dummy3.txt",
				"/root/foo2",
				"/root/foo2/bar",
				"/root/foo2/bar/dummy2.txt",
				"/root/foo3",
				"/root/foo3/baz",
				"/root/foo3/baz/dummy.txt",
				"/root/foo3/dummy1.txt",
			},
			wantNames: []string{
				"foo1",
				"bar",
				"file1.txt",
				"file2.txt",
				"baz",
				"dummy3.txt",
				"foo2",
				"bar",
				"dummy2.txt",
				"foo3",
				"baz",
				"dummy.txt",
				"dummy1.txt",
			},
		},
		{
			name:     "subfolder walk",
			basePath: "/root",
			wantPaths: []string{
				"/root/foo1",
				"/root/foo1/bar",
				"/root/foo1/bar/file1.txt",
				"/root/foo1/bar/file2.txt",
				"/root/foo1/baz",
				"/root/foo1/baz/dummy3.txt",
				"/root/foo2",
				"/root/foo2/bar",
				"/root/foo2/bar/dummy2.txt",
				"/root/foo3",
				"/root/foo3/baz",
				"/root/foo3/baz/dummy.txt",
				"/root/foo3/dummy1.txt",
			},
			wantNames: []string{
				"foo1",
				"bar",
				"file1.txt",
				"file2.txt",
				"baz",
				"dummy3.txt",
				"foo2",
				"bar",
				"dummy2.txt",
				"foo3",
				"baz",
				"dummy.txt",
				"dummy1.txt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualPath := make([]string, 0)
			actualNames := make([]string, 0)
			err := client.Walk(ctx, tt.basePath, func(entry *ftp.Entry, skip func()) error {
				actualPath = append(actualPath, entry.Path)
				actualNames = append(actualNames, entry.Name)
				return nil
			})

			if assert.NoError(t, err, "walk failed") {
				assert.Equal(t, tt.wantNames, actualNames)
				assert.Equal(t, tt.wantPaths, actualPath)
			}
		})
	}
}

func TestTestDoubleClient_Size(t *testing.T) {
	ctx := context.TODO()

	vol := files.NewVolume()
	require.NoError(t, vol.CreateTemp("samples*"), "can't create fixtures dir")
	defer func() {
		require.NoError(t, vol.Remove(), "can't cleanup fixtures dir")
	}()

	file, err := os.OpenFile(path.Join(vol.Path(), "dummy.txt"), os.O_CREATE|os.O_WRONLY, files.DefaultDirPermissions)
	require.NoError(t, err, "failed to create sample file")

	n, err := io.Copy(file, bytes.NewBufferString("fake data"))
	require.NoError(t, err, "failed to write sample file")

	require.NoError(t, file.Close(), "failed to close sample file")

	client := ftp.NewTestDouble(vol.Path())

	tests := []struct {
		name     string
		filePath string
		want     int64
	}{
		{
			name:     "Get file size",
			filePath: "/dummy.txt",
			want:     n,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := client.Size(ctx, tt.filePath)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, size)
			}
		})
	}
}

func TestTestDoubleClient_List(t *testing.T) {
	ctx := context.TODO()

	vol := files.NewVolume()
	require.NoError(t, vol.CreateTemp("samples*"), "can't create fixtures dir")
	defer func() {
		require.NoError(t, vol.Remove(), "can't cleanup fixtures dir")
	}()

	require.NoError(t, files.DirCreateEmptyFiles(vol.Path(),
		"root/foo1/bar/file1.txt",
		"root/foo1/bar/file2.txt",
		"root/foo2/bar/dummy2.txt",
		"root/dummy1.txt",
	), "failed to create fixture files")

	client := ftp.NewTestDouble(vol.Path())

	tests := []struct {
		name      string
		basePath  string
		wantPaths []string
	}{
		{
			name:     "List files in directory",
			basePath: "root/foo1/bar",
			wantPaths: []string{
				"/root/foo1/bar/file1.txt",
				"/root/foo1/bar/file2.txt",
			},
		},
		{
			name:     "List folders and files in directory",
			basePath: "root",
			wantPaths: []string{
				"/root/dummy1.txt",
				"/root/foo1",
				"/root/foo2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualPaths := make([]string, 0)
			list, err := client.List(ctx, tt.basePath)
			for _, entry := range list {
				actualPaths = append(actualPaths, entry.Path)
			}

			if assert.NoError(t, err, "list failed") {
				assert.Equal(t, tt.wantPaths, actualPaths)
			}
		})
	}
}

func TestTestDoubleClient_NameList(t *testing.T) {
	ctx := context.TODO()

	vol := files.NewVolume()
	require.NoError(t, vol.CreateTemp("samples*"), "can't create fixtures dir")
	defer func() {
		require.NoError(t, vol.Remove(), "can't cleanup fixtures dir")
	}()

	require.NoError(t, files.DirCreateEmptyFiles(vol.Path(),
		"root/foo1/bar/file1.txt",
		"root/foo1/bar/file2.txt",
		"root/foo2/bar/dummy2.txt",
		"root/dummy1.txt",
	), "failed to create fixture files")

	client := ftp.NewTestDouble(vol.Path())

	tests := []struct {
		name      string
		basePath  string
		wantPaths []string
	}{
		{
			name:     "List files in directory",
			basePath: "root/foo1/bar",
			wantPaths: []string{
				"file1.txt",
				"file2.txt",
			},
		},
		{
			name:     "List folders and files in directory",
			basePath: "root",
			wantPaths: []string{
				"dummy1.txt",
				"foo1",
				"foo2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := client.NameList(ctx, tt.basePath)

			if assert.NoError(t, err, "list names failed") {
				assert.Equal(t, tt.wantPaths, list)
			}
		})
	}
}

func TestTestDoubleClient_Download(t *testing.T) {
	ctx := context.TODO()

	vol := files.NewVolume()
	require.NoError(t, vol.CreateTemp("samples*"), "can't create fixtures dir")
	defer func() {
		require.NoError(t, vol.Remove(), "can't cleanup fixtures dir")
	}()

	file, err := os.OpenFile(path.Join(vol.Path(), "dummy.txt"), os.O_CREATE|os.O_WRONLY, files.DefaultDirPermissions)
	require.NoError(t, err, "failed to create sample file")

	_, err = io.Copy(file, bytes.NewBufferString("fake data"))
	require.NoError(t, err, "failed to write sample file")

	require.NoError(t, file.Close(), "failed to close sample file")

	client := ftp.NewTestDouble(vol.Path())

	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "Get file size",
			filePath: "/dummy.txt",
			want:     "fake data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := client.Download(ctx, tt.filePath)
			if assert.NoError(t, err) {
				defer func() {
					require.NoError(t, content.Close(), "can't close content")
				}()

				raw, err := io.ReadAll(content)
				require.NoError(t, err)

				assert.Equal(t, tt.want, string(raw))
			}
		})
	}
}
