package es

import (
	"context"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type (
	IndexNameMapper func(name string) string

	Client struct {
		ES              *elasticsearch.Client
		indexNameMapper IndexNameMapper
		startIndexer    func(ctx context.Context) (*indexer, error)
	}

	IndexerJobItem struct {
		Action     string
		Index      string
		DocumentID string
		Body       any
	}

	IndexerJobItemCallback func(context.Context, error)

	indexer struct {
		ES               esutil.BulkIndexer
		indexNameMapper  IndexNameMapper
		isClosedFlagLock sync.RWMutex
		isClosedFlag     bool
	}
)
