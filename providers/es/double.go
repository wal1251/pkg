package es

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/providers/es/indices"
	"github.com/wal1251/pkg/providers/es/search"
	"github.com/wal1251/pkg/tools/serial"
)

var (
	_ ElasticSearch = (*ClientTestDouble)(nil)
	_ Indexer       = (*indexerTestDouble)(nil)

	ErrIndexAlreadyExists = errors.New("index already exists")
	ErrIndexNotFound      = errors.New("index not found")
)

type (
	indexerTestDouble struct{}

	ClientTestDouble struct {
		searchCallback  func(search.MultiSearchItem) search.Response
		indexesLock     sync.RWMutex
		indexes         map[string]api.Object
		indexNameMapper IndexNameMapper
	}
)

func (i *indexerTestDouble) Add(ctx context.Context, item IndexerJobItem, callback IndexerJobItemCallback) error {
	logger := logs.FromContext(ctx)

	req, err := serial.ToBytes(item, serial.JSONEncode[IndexerJobItem])
	if err != nil {
		logger.Err(err).Msg("failed to marshal ES indexer job")

		return err
	}

	logger.Debug().Msgf("added ES indexer dummy job: %s", string(req))

	callback(ctx, nil)

	return nil
}

func (i *indexerTestDouble) Stats() map[string]uint64 {
	return nil
}

func (i *indexerTestDouble) Stop(ctx context.Context) {
	logger := logs.FromContext(ctx)
	logger.Debug().Msg("stop ES indexer")
}

func (t *ClientTestDouble) StartIndexer(ctx context.Context) (Indexer, error) {
	logger := logs.FromContext(ctx)
	logger.Debug().Msg("start ES indexer")

	return &indexerTestDouble{}, nil
}

func (t *ClientTestDouble) Ping(ctx context.Context) error {
	logs.FromContext(ctx).Debug().Msg("ping ES")

	return nil
}

func (t *ClientTestDouble) IndexGet(ctx context.Context, indexName string) (api.Object, error) {
	indexName = t.indexNameMapper.Map(indexName)

	logger := logs.FromContext(ctx)

	t.indexesLock.RLock()
	defer t.indexesLock.RUnlock()
	index, ok := t.indexes[indexName]
	if !ok {
		logger.Debug().Msgf("ES index not found: %s", indexName)

		return (api.Object)(nil), nil
	}

	index = index.Copy()
	logs.FromContext(ctx).Debug().Msgf("got ES index '%s': %v", indexName, index)

	return index, nil
}

func (t *ClientTestDouble) IndexCreate(ctx context.Context, indexName string, body api.Factory[indices.Create]) error {
	indexName = t.indexNameMapper.Map(indexName)

	logger := logs.FromContext(ctx)

	object, err := serial.JSONDecode[api.Object](body.JSONReader())
	if err != nil {
		logger.Err(err).Msgf("failed to decode body on ES index create request: %s", indexName)

		return err
	}

	logger.Debug().Msgf("creating ES index '%s': %v", indexName, object)

	t.indexesLock.Lock()
	defer t.indexesLock.Unlock()
	if _, ok := t.indexes[indexName]; ok {
		logger.Debug().Msgf("index already exists: %s", indexName)

		return fmt.Errorf("%w: %s", ErrIndexAlreadyExists, indexName)
	}

	t.indexes[indexName] = object

	logger.Debug().Msgf("successfully created ES index: %s", indexName)

	return nil
}

func (t *ClientTestDouble) IndexMappingUpdate(ctx context.Context, indexName string, body api.Factory[indices.Mapping]) error {
	indexName = t.indexNameMapper.Map(indexName)

	logger := logs.FromContext(ctx)
	logger.Debug().Msgf("updating es index mapping: %s", indexName)

	object, err := serial.JSONDecode[api.Object](body.JSONReader())
	if err != nil {
		logger.Err(err).Msgf("failed to decode body on ES index update mapping request: %s", indexName)

		return err
	}

	t.indexesLock.Lock()
	defer t.indexesLock.Unlock()

	index, ok := t.indexes[indexName]
	if !ok {
		logger.Debug().Msgf("index not found: %s", indexName)

		return fmt.Errorf("%w: %s", ErrIndexNotFound, indexName)
	}

	index.ApplyFrom(api.Object{"mappings": object})
	logger.Debug().Msgf("updated ES index mapping '%s': %v", indexName, index)

	return nil
}

func (t *ClientTestDouble) MultiSearch(ctx context.Context, body search.MultiSearch, _ ...string) (*search.MultiResponse, error) {
	logger := logs.FromContext(ctx)

	buf := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(buf, body.NewNDJSONReader()); err != nil {
		logger.Err(err).Msgf("failed to marshall ES MultiSearch request")

		return nil, fmt.Errorf("can't read search request: %w", err)
	}

	logger.Debug().Msgf("querying ES MultiSearch:\n%v", buf)

	t.indexesLock.RLock()
	defer t.indexesLock.RUnlock()

	responses := make([]search.Response, 0, len(body))
	for _, item := range body {
		var resp search.Response
		indexName := t.indexNameMapper.Map(item.Header.Index)
		if _, ok := t.indexes[indexName]; ok {
			if t.searchCallback != nil {
				resp = t.searchCallback(item)
			}
		} else {
			logger.Debug().Msgf("ES index not found: %s", indexName)
			resp = search.Response{Status: http.StatusNotFound}
		}
		responses = append(responses, resp)
	}

	response := search.MultiResponse{Responses: responses}
	logger.Debug().Msgf("query ES MultiSearch result: %v", response)

	return &response, nil
}

func (t *ClientTestDouble) DocumentIndex(context.Context, string, string, api.Object) (api.Object, error) {
	return nil, nil
}

func (t *ClientTestDouble) DocumentDelete(context.Context, string, string) error {
	return nil
}

func (t *ClientTestDouble) DocumentDeleteByQuery(context.Context, []string, search.Search) error {
	return nil
}

func NewClientTestDouble(cfg *Config, search func(search.MultiSearchItem) search.Response) (*ClientTestDouble, error) {
	return &ClientTestDouble{
		indexes:         make(map[string]api.Object),
		indexNameMapper: NewIndexNameMapper(cfg),
		searchCallback:  search,
	}, nil
}
