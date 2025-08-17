package ftp

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/jlaffaye/ftp"

	"github.com/wal1251/pkg/core"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/collections"
)

var _ Client = (*DefaultClient)(nil)

type (
	// DefaultClient реализация клиента FTP Client по умолчанию.
	DefaultClient struct {
		core.Component
		connection    *ftp.ServerConn
		newConnection func(ctx context.Context) (*ftp.ServerConn, error)
	}
)

// Dial см. Client.Dial().
func (c *DefaultClient) Dial(ctx context.Context) error {
	logger := logs.FromContext(ctx)
	log := logs.LocalContext(c).To

	conn, err := c.newConnection(ctx)
	if err != nil {
		return err
	}

	c.connection = conn

	log(logger.Debug).Msg("connection established")

	return nil
}

// Quit см. Client.Quit().
func (c *DefaultClient) Quit(ctx context.Context) {
	logger := logs.FromContext(ctx)
	log := logs.LocalContext(c).To

	if c.connection != nil {
		if err := c.connection.Quit(); err != nil {
			log(logger.Error).Err(err).Msg("failed to quit server")

			return
		}
	}

	log(logger.Debug).Msg("connection closed")
}

// Walk см. Client.Walk().
func (c *DefaultClient) Walk(ctx context.Context, path string, accept Visitor) error {
	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("traverse folder: %s", path)

	w := c.connection.Walk(path)
	for w.Next() {
		if err := accept(MakeEntry(w.Path(), w.Stat()), w.SkipDir); err != nil {
			return err
		}
	}

	return nil
}

// List см. Client.List().
func (c *DefaultClient) List(ctx context.Context, folderPath string) ([]*Entry, error) {
	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("list folder: %s", folderPath)

	entries, err := c.connection.List(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error while list FTP folder %s: %w", folderPath, err)
	}

	return collections.Map(entries, func(e *ftp.Entry) *Entry { return MakeEntry(path.Join(folderPath, e.Name), e) }), nil
}

// Size см. Client.Size().
func (c *DefaultClient) Size(ctx context.Context, filePath string) (int64, error) {
	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("get file size: %s", filePath)

	size, err := c.connection.FileSize(filePath)
	if err != nil {
		return 0, fmt.Errorf("error while requsting FTP file size %s: %w", filePath, err)
	}

	return size, nil
}

// NameList см. Client.NameList().
func (c *DefaultClient) NameList(ctx context.Context, folderPath string) ([]string, error) {
	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("list folder file names: %s", folderPath)

	entries, err := c.connection.NameList(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error while list FTP folder %s: %w", folderPath, err)
	}

	return entries, nil
}

// Download см. Client.Download().
func (c *DefaultClient) Download(ctx context.Context, filePath string) (io.ReadCloser, error) {
	logs.LocalContext(c).To(logs.FromContext(ctx).Debug).Msgf("download file: %s", filePath)

	resp, err := c.connection.Retr(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch file %s from ftp: %w", filePath, err)
	}

	return resp, nil
}

// NewClient возвращает реализацию клиента FTP Client.
func NewClient(cfg *Config) *DefaultClient {
	return &DefaultClient{
		Component:     core.NewDefaultComponent(ComponentName, cfg.Address),
		newConnection: func(ctx context.Context) (*ftp.ServerConn, error) { return Connect(ctx, cfg) },
	}
}
