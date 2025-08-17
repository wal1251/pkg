package files_test

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/files"
)

func TestIsEmptyDir(t *testing.T) {
	tests := []struct {
		name    string
		withDir func(act func(dir string))
		want    bool
	}{
		{
			name: "Пустая директория",
			withDir: func(act func(string)) {
				tempDir, err := os.MkdirTemp("", "sampleDir*")
				require.NoError(t, err)

				defer func() {
					require.NoError(t, os.RemoveAll(tempDir))
				}()

				act(tempDir)
			},
			want: true,
		},
		{
			name: "Не пустая директория",
			withDir: func(act func(string)) {
				tempDir, err := os.MkdirTemp("", "sampleDir*")
				require.NoError(t, err)

				_, err = os.CreateTemp(tempDir, "sampleFile*")
				require.NoError(t, err)

				defer func() {
					require.NoError(t, os.RemoveAll(tempDir))
				}()

				act(tempDir)
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.withDir(func(dir string) {
				isEmpty, err := files.DirIsEmpty(dir)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, isEmpty, "unexpected result")
				}
			})
		})
	}
}

func TestDirExists(t *testing.T) {
	tests := []struct {
		name    string
		withDir func(act func(dir string))
		want    bool
	}{
		{
			name: "Директория существует",
			withDir: func(act func(dir string)) {
				tempDir, err := os.MkdirTemp("", "sampleDir*")
				require.NoError(t, err)

				defer func() {
					require.NoError(t, os.RemoveAll(tempDir))
				}()

				act(tempDir)
			},
			want: true,
		},
		{
			name: "Директория не существует",
			withDir: func(act func(dir string)) {
				act(path.Join("/", uuid.NewString()))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.withDir(func(dir string) {
				exists, err := files.DirExists(dir)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, exists, "unexpected result")
				}
			})
		})
	}
}

func TestDirMustExist(t *testing.T) {
	tests := []struct {
		name    string
		withDir func(act func(dir string))
		wantErr error
	}{
		{
			name: "Директория существует",
			withDir: func(act func(dir string)) {
				tempDir, err := os.MkdirTemp("", "sampleDir*")
				require.NoError(t, err)

				defer func() {
					require.NoError(t, os.RemoveAll(tempDir))
				}()

				act(tempDir)
			},
		},
		{
			name: "Директория не существует",
			withDir: func(act func(dir string)) {
				act(path.Join("/", uuid.NewString()))
			},
			wantErr: files.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.withDir(func(dir string) {
				err := files.DirMustExist(dir)
				if tt.wantErr == nil {
					assert.NoError(t, err, "no error expected")
				} else {
					if assert.Error(t, err) {
						assert.ErrorIs(t, err, tt.wantErr)
					}
				}
			})
		})
	}
}

func TestDirCreateStructure(t *testing.T) {
	tests := []struct {
		name     string
		dirPaths []string
		wantDirs []string
	}{
		{
			name: "Базовый кейс",
			dirPaths: []string{
				"",
				"/",
				"foo",
				"bar/",
				"baz/foo/bar",
			},
			wantDirs: []string{
				"foo",
				"bar",
				"baz",
				"baz/foo",
				"baz/foo/bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath, err := os.MkdirTemp(os.TempDir(), "sampleDir_*")
			require.NoError(t, err, "failed to make test dir")

			actualContent := make([]string, 0, len(tt.dirPaths))

			err = files.DirCreateStructure(basePath, tt.dirPaths...)
			if assert.NoError(t, err) {
				err = filepath.WalkDir(basePath, func(filePath string, entry fs.DirEntry, err error) error {
					if filePath == basePath {
						return nil
					}

					if assert.True(t, entry.IsDir(), "directory expected") {
						relPath, err := filepath.Rel(basePath, filePath)
						require.NoError(t, err, "can't get relative path")

						actualContent = append(actualContent, relPath)
					}

					return nil
				})
				require.NoError(t, err, "failed to check dir")

				assert.ElementsMatch(t, tt.wantDirs, actualContent, "actual dirs doesn't match expected")
			}

			err = os.RemoveAll(basePath)
			require.NoError(t, err, "failed to cleanup test dir")
		})
	}
}

func TestDirCreateEmptyFiles(t *testing.T) {
	tests := []struct {
		name      string
		filePaths []string
		wantDirs  []string
		wantFiles []string
	}{
		{
			name: "Базовый кейс",
			filePaths: []string{
				"file1.txt",
				"/file2.txt",
				"foo/file3.txt",
				"baz/foo/bar/file4.txt",
			},
			wantDirs: []string{
				"foo",
				"baz",
				"baz/foo",
				"baz/foo/bar",
			},
			wantFiles: []string{
				"file1.txt",
				"file2.txt",
				"foo/file3.txt",
				"baz/foo/bar/file4.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath, err := os.MkdirTemp(os.TempDir(), "sampleDir_*")
			require.NoError(t, err, "failed to make test dir")

			actualFiles := make([]string, 0, len(tt.filePaths))
			actualDirs := make([]string, 0, len(tt.filePaths))

			err = files.DirCreateEmptyFiles(basePath, tt.filePaths...)
			if assert.NoError(t, err) {
				err = filepath.WalkDir(basePath, func(filePath string, entry fs.DirEntry, err error) error {
					if filePath == basePath {
						return nil
					}

					relPath, err := filepath.Rel(basePath, filePath)
					require.NoError(t, err, "can't get relative path")

					if entry.IsDir() {
						actualDirs = append(actualDirs, relPath)
					} else {
						actualFiles = append(actualFiles, relPath)
					}

					return nil
				})
				require.NoError(t, err, "failed to check dir")

				assert.ElementsMatch(t, tt.wantDirs, actualDirs, "actual dirs doesn't match expected")
				assert.ElementsMatch(t, tt.wantFiles, actualFiles, "actual files doesn't match expected")
			}

			err = os.RemoveAll(basePath)
			require.NoError(t, err, "failed to cleanup test dir")
		})
	}
}

func TestDirTraverse(t *testing.T) {
	var actualContent []string

	with := func(samples []string, act func(string)) {
		tempDir, err := os.MkdirTemp("", "sampleDir*")
		require.NoError(t, err)

		err = files.DirCreateEmptyFiles(tempDir, samples...)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.RemoveAll(tempDir))
		}()

		act(tempDir)
	}

	tests := []struct {
		name     string
		files    []string
		callback files.DirVisitor
		want     []string
	}{
		{
			name: "Сканирование директории полностью",
			files: []string{
				"foo/bar/file1.txt",
				"baz/file2.txt",
				"file3.txt",
			},
			callback: func(entry files.RelativeEntry, _ func()) error {
				actualContent = append(actualContent, entry.PathRelative)

				return nil
			},
			want: []string{
				".",
				"baz",
				"baz/file2.txt",
				"file3.txt",
				"foo",
				"foo/bar",
				"foo/bar/file1.txt",
			},
		},
		{
			name: "Сканирование директории частично",
			files: []string{
				"foo/bar/file1.txt",
				"baz/file2.txt",
				"foo/bar/baz/file3.txt",
				"file3.txt",
			},
			callback: func(entry files.RelativeEntry, skip func()) error {
				actualContent = append(actualContent, entry.PathRelative)

				if entry.PathRelative == "foo/bar" {
					skip()
				}

				return nil
			},
			want: []string{
				".",
				"baz",
				"baz/file2.txt",
				"file3.txt",
				"foo",
				"foo/bar",
			},
		},
		{
			name: "Сканирование пустой директории",
			callback: func(entry files.RelativeEntry, skip func()) error {
				actualContent = append(actualContent, entry.PathRelative)

				return nil
			},
			want: []string{"."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			with(tt.files, func(dir string) {
				actualContent = make([]string, 0, 1)
				err := files.DirTraverse(dir, tt.callback)
				if assert.NoError(t, err) {
					assert.ElementsMatch(t, tt.want, actualContent, "callback result doesn't matches expected")
				}
			})
		})
	}
}

func TestDoCopy(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "SrcDir*")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(srcDir))
	}()

	destDir, err := os.MkdirTemp("", "DestDir*")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(destDir))
	}()

	err = files.DirCreateEmptyFiles(srcDir,
		"file.txt",
		"foo/bar/file2.txt",
		"baz/file3.txt",
	)
	require.NoError(t, err)

	err = files.DirTraverse(srcDir, files.DoCopy(destDir))
	if assert.NoError(t, err) {
		content := make([]string, 0, 1)
		err = files.DirTraverse(destDir, func(entry files.RelativeEntry, _ func()) error {
			content = append(content, entry.PathRelative)
			return nil
		})
		require.NoError(t, err, "failed to check dest folder")

		assert.ElementsMatch(t, content, []string{
			".",
			"baz",
			"baz/file3.txt",
			"file.txt",
			"foo",
			"foo/bar",
			"foo/bar/file2.txt",
		})
	}
}
