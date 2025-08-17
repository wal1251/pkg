package ftp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/textproto"
	"os"
	"path"
	"strings"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/files"
)

var _ Client = (*TestDoubleClient)(nil)

const (
	labelTestDouble = "TEST DOUBLE"
)

// TestDoubleClient тестовый двойник, заменяет реализацию клиента Client для использования в тестах.
type TestDoubleClient struct {
	core.Component
	homePath string
}

// Dial см. Client.Dial().
func (c *TestDoubleClient) Dial(ctx context.Context) error {
	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msg("connection established")

	return nil
}

// Quit см. Client.Quit().
func (c *TestDoubleClient) Quit(ctx context.Context) {
	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msg("connection closed")
}

// Walk см. Client.Walk().
func (c *TestDoubleClient) Walk(ctx context.Context, folderPath string, accept Visitor) error {
	fetchPath := path.Join(c.homePath, folderPath)

	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("traverse folder: %s", fetchPath)

	if err := files.DirTraverse(fetchPath, func(entry files.RelativeEntry, skip func()) error {
		if entry.IsHome() {
			return nil
		}

		dirEntry, err := makeEntryFromDirEntry(entry, path.Join(folderPath, entry.PathRelative))
		if err != nil {
			return err
		}

		return accept(dirEntry, skip)
	}); err != nil {
		if errors.Is(err, files.ErrNotExist) {
			return &textproto.Error{
				Code: ErrCodeNotFound,
				Msg:  err.Error(),
			}
		}

		return err
	}

	return nil
}

// Size см. Client.Size().
func (c *TestDoubleClient) Size(ctx context.Context, filePath string) (int64, error) {
	fetchPath := path.Join(c.homePath, filePath)

	logger := logs.FromContext(ctx)
	log := logs.LocalContext(c).To

	log(logger.Debug).Msgf("get file size: %s", fetchPath)

	entry, err := os.Open(fetchPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, &textproto.Error{
				Code: ErrCodeNotFound,
				Msg:  err.Error(),
			}
		}

		return 0, fmt.Errorf("can't open file %s: %w", fetchPath, err)
	}

	defer func() {
		if err := entry.Close(); err != nil {
			log(logger.Error).Err(err).Msgf("failed to close file: %s", fetchPath)
		}
	}()

	stats, err := entry.Stat()
	if err != nil {
		return 0, fmt.Errorf("can't read entry stats %s: %w", filePath, err)
	}

	return stats.Size(), nil
}

// List см. Client.List().
func (c *TestDoubleClient) List(ctx context.Context, folderPath string) ([]*Entry, error) {
	fetchPath := path.Join(c.homePath, folderPath)

	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("list folder: %s", fetchPath)

	entries, err := os.ReadDir(fetchPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, &textproto.Error{
				Code: ErrCodeNotFound,
				Msg:  err.Error(),
			}
		}

		return nil, fmt.Errorf("can't read dir %s: %w", fetchPath, err)
	}

	return collections.MapWithErr(entries, func(e os.DirEntry) (*Entry, error) {
		return makeEntryFromDirEntry(e, path.Join(folderPath, e.Name()))
	})
}

// NameList см. Client.NameList().
func (c *TestDoubleClient) NameList(ctx context.Context, folderPath string) ([]string, error) {
	fetchPath := path.Join(c.homePath, folderPath)

	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("list folder file names: %s", fetchPath)

	entries, err := os.ReadDir(fetchPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, &textproto.Error{
				Code: ErrCodeNotFound,
				Msg:  err.Error(),
			}
		}

		return nil, fmt.Errorf("can't read file %s: %w", fetchPath, err)
	}

	return collections.Map(entries, func(entry os.DirEntry) string { return entry.Name() }), nil
}

// Download см. Client.Download().
func (c *TestDoubleClient) Download(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fetchPath := path.Join(c.homePath, filePath)

	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("list folder file names: %s", fetchPath)

	content, err := os.Open(fetchPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, &textproto.Error{
				Code: ErrCodeNotFound,
				Msg:  err.Error(),
			}
		}

		return nil, fmt.Errorf("can't read file %s: %w", fetchPath, err)
	}

	return content, nil
}

// NewTestDouble возвращает нового тестового двойника для Client.
func NewTestDouble(homePath string) *TestDoubleClient {
	return &TestDoubleClient{
		Component: core.NewDefaultComponent(ComponentName, labelTestDouble),
		homePath:  homePath,
	}
}

func makeEntryFromDirEntry(entry os.DirEntry, entryPath string) (*Entry, error) {
	info, err := entry.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to read file info: %w", err)
	}

	entryType := EntryTypeFile
	if entry.Type().IsDir() {
		entryType = EntryTypeFolder
	}

	if !strings.HasPrefix(entryPath, "/") {
		entryPath = path.Join("/", entryPath)
	}

	return &Entry{
		Name: entry.Name(),
		Path: entryPath,
		Type: entryType,
		Size: uint64(info.Size()),
		Time: info.ModTime(),
	}, nil
}
