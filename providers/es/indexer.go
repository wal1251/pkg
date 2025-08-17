package es

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/tools/serial"
)

const (
	ActionIndex  = "index"
	ActionDelete = "delete"

	IndexerStatAdded    = "added"
	IndexerStatFlushed  = "flushed"
	IndexerStatFailed   = "failed"
	IndexerStatIndexed  = "indexed"
	IndexerStatCreated  = "created"
	IndexerStatUpdated  = "updated"
	IndexerStatDeleted  = "deleted"
	IndexerStatRequests = "requests"

	IndexerDefaultFlushInterval = 30 * time.Second
)

var _ Indexer = (*indexer)(nil)

type (
	Indexer interface {
		// Add добавляет элемент в Indexer.
		Add(ctx context.Context, item IndexerJobItem, callback IndexerJobItemCallback) error

		// Stats возвращает статистку Indexer.
		Stats() map[string]uint64

		// Stop ждет, пока все добавленные item не будут сброшены, а потом закроет Indexer.
		Stop(context.Context)
	}

	Document interface {
		fmt.Stringer

		DocumentID() string
		Index() string
	}
)

var ErrIndexerStopped = errors.New("indexer stopped")

func (i IndexerJobItem) String() string {
	return fmt.Sprintf("%s: %s %s", i.Index, i.Action, i.DocumentID)
}

func (c IndexerJobItemCallback) Done(ctx context.Context) {
	if c != nil {
		c(ctx, nil)
	}
}

func (c IndexerJobItemCallback) Failed(ctx context.Context, err error) {
	if c != nil {
		c(ctx, err)
	}
}

func (i *indexer) checkIsClosed() bool {
	i.isClosedFlagLock.RLock()
	defer i.isClosedFlagLock.RUnlock()

	return i.isClosedFlag
}

func (i *indexer) Add(ctx context.Context, item IndexerJobItem, callback IndexerJobItemCallback) error {
	if i.checkIsClosed() {
		return fmt.Errorf("can't add item to index: %w", ErrIndexerStopped)
	}

	var body io.ReadSeeker

	if item.Body != nil {
		jsonBody, err := serial.ToBytes(item.Body, serial.JSONEncode[any])
		if err != nil {
			return err
		}

		body = bytes.NewReader(jsonBody)
	}

	esItem := esutil.BulkIndexerItem{
		Index:      i.indexNameMapper.Map(item.Index),
		Action:     item.Action,
		DocumentID: item.DocumentID,
		Body:       body,
		OnSuccess: func(ctx context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem) {
			callback.Done(ctx)
		},
		OnFailure: func(ctx context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem, err error) {
			callback.Failed(ctx, err)
		},
	}

	if err := i.ES.Add(ctx, esItem); err != nil {
		return fmt.Errorf("failed to add new indexer job: %w", err)
	}

	return nil
}

func (i *indexer) Stop(ctx context.Context) {
	closeIndexer := func() bool {
		i.isClosedFlagLock.Lock()
		defer i.isClosedFlagLock.Unlock()

		if !i.isClosedFlag {
			i.isClosedFlag = true
		}

		return i.isClosedFlag
	}

	if !closeIndexer() {
		if err := i.ES.Close(ctx); err != nil {
			logs.FromContext(ctx).Err(err).Msg("unable to close indexer")
		}
	}
}

func (i *indexer) Stats() map[string]uint64 {
	stats := i.ES.Stats()

	return map[string]uint64{
		IndexerStatAdded:    stats.NumAdded,
		IndexerStatFlushed:  stats.NumFlushed,
		IndexerStatFailed:   stats.NumFailed,
		IndexerStatIndexed:  stats.NumIndexed,
		IndexerStatCreated:  stats.NumCreated,
		IndexerStatUpdated:  stats.NumUpdated,
		IndexerStatDeleted:  stats.NumDeleted,
		IndexerStatRequests: stats.NumRequests,
	}
}

func NewIndexerFactory(cfg *Config, elastic *elasticsearch.Client, debug bool) func(ctx context.Context) (*indexer, error) {
	return func(ctx context.Context) (*indexer, error) {
		logger := logs.FromContext(ctx)

		var debugLogger esutil.BulkIndexerDebugLogger
		if debug {
			debugLogger = logger
		}

		esIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client:        elastic,
			FlushInterval: cfg.IndexerFlushInterval,
			Timeout:       cfg.Timeout,
			DebugLogger:   debugLogger,
			OnFlushStart:  logger.WithContext,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create elastic search bulk indexer: %w", err)
		}

		return &indexer{
			ES:              esIndexer,
			indexNameMapper: NewIndexNameMapper(cfg),
		}, nil
	}
}

func NewDocumentIndexJob(doc Document) IndexerJobItem {
	return IndexerJobItem{
		Action:     ActionIndex,
		DocumentID: doc.DocumentID(),
		Index:      doc.Index(),
		Body:       doc,
	}
}

func NewDocumentDeleteJob(doc Document) IndexerJobItem {
	return IndexerJobItem{
		Action:     ActionDelete,
		DocumentID: doc.DocumentID(),
		Index:      doc.Index(),
	}
}
