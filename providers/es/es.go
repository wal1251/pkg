// Package es предоставляет адаптер для работы с ElasticSearch.
package es

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"

	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/providers/es/api"
	"github.com/wal1251/pkg/providers/es/indices"
	"github.com/wal1251/pkg/providers/es/search"
	"github.com/wal1251/pkg/tools/collections"
	"github.com/wal1251/pkg/tools/serial"
)

var (
	ErrESNotFound           = errors.New("not found")
	ErrESErrorResponse      = errors.New("error ES response")
	ErrESUnexpectedResponse = errors.New("unexpected ES response")
)

var _ ElasticSearch = (*Client)(nil)

type (
	ElasticSearch interface {
		StartIndexer(ctx context.Context) (Indexer, error)

		// Ping проверяет работу кластер.
		Ping(ctx context.Context) error

		// IndexGet возвращает информацию об индексе.
		IndexGet(ctx context.Context, indexName string) (api.Object, error)

		// IndexCreate создает индекс.
		IndexCreate(ctx context.Context, index string, body api.Factory[indices.Create]) error

		// IndexMappingUpdate обновляет сопоставление с индексом.
		IndexMappingUpdate(ctx context.Context, index string, body api.Factory[indices.Mapping]) error

		// MultiSearch позволяет сделать несколько операций поиска в один вызов.
		MultiSearch(ctx context.Context, body search.MultiSearch, filterPath ...string) (*search.MultiResponse, error)

		// DocumentIndex создает или обновляет документ в индексе.
		DocumentIndex(ctx context.Context, index string, id string, object api.Object) (api.Object, error)

		// DocumentDelete удаляет документ из индекса.
		DocumentDelete(ctx context.Context, index string, id string) error

		// DocumentDeleteByQuery удаляет документы, сопоставленные запросу.
		DocumentDeleteByQuery(ctx context.Context, index []string, s search.Search) error
	}
)

func (c *Client) StartIndexer(ctx context.Context) (Indexer, error) {
	if c.startIndexer != nil {
		return c.startIndexer(ctx)
	}

	return (Indexer)(nil), nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := execute[any](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.Ping(
			c.ES.API.Ping.WithContext(ctx),
			c.ES.API.Ping.WithOpaqueID(logs.RequestID(ctx)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to check es connection: %w", err)
		}

		return resp, nil
	})

	return err
}

func (c *Client) IndexGet(ctx context.Context, indexName string) (api.Object, error) {
	indexName = c.indexNameMapper.Map(indexName)

	resp, err := execute[api.Object](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.Indices.Get(
			[]string{indexName},
			c.ES.API.Indices.Get.WithContext(ctx),
			c.ES.API.Indices.Get.WithOpaqueID(logs.RequestID(ctx)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch es index %s: %w", indexName, err)
		}

		return resp, nil
	})
	if err != nil {
		if errors.Is(err, ErrESNotFound) {
			return nil, nil
		}

		return nil, err
	}

	if index, ok := resp[indexName].(map[string]any); ok {
		return index, nil
	}

	return nil, ErrESUnexpectedResponse
}

func (c *Client) IndexCreate(ctx context.Context, index string, body api.Factory[indices.Create]) error {
	_, err := execute[any](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.Indices.Create(
			c.indexNameMapper.Map(index),
			c.ES.API.Indices.Create.WithBody(body.JSONReader()),
			c.ES.API.Indices.Create.WithContext(ctx),
			c.ES.API.Indices.Create.WithOpaqueID(logs.RequestID(ctx)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create es index %s: %w", index, err)
		}

		return resp, nil
	})

	return err
}

func (c *Client) IndexMappingUpdate(ctx context.Context, index string, body api.Factory[indices.Mapping]) error {
	_, err := execute[any](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.Indices.PutMapping(
			[]string{c.indexNameMapper.Map(index)},
			body.JSONReader(),
			c.ES.API.Indices.PutMapping.WithContext(ctx),
			c.ES.API.Indices.PutMapping.WithOpaqueID(logs.RequestID(ctx)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to put mapping for es index %s: %w", index, err)
		}

		return resp, nil
	})

	return err
}

func (c *Client) MultiSearch(ctx context.Context, body search.MultiSearch, filterPath ...string) (*search.MultiResponse, error) {
	indexMap := func(item search.MultiSearchItem) search.MultiSearchItem {
		if item.Header.Index != "" {
			item.Header.Index = c.indexNameMapper.Map(item.Header.Index)
		}

		return item
	}

	for i, item := range body {
		body[i] = indexMap(item)
	}

	return execute[*search.MultiResponse](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.Msearch(
			body.NewNDJSONReader(),
			c.ES.API.Msearch.WithContext(ctx),
			c.ES.API.Msearch.WithOpaqueID(logs.RequestID(ctx)),
			c.ES.API.Msearch.WithFilterPath(filterPath...),
			c.ES.API.Msearch.WithPretty(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query multi search es request: %w", err)
		}

		return resp, nil
	})
}

func (c *Client) DocumentIndex(ctx context.Context, index string, id string, object api.Object) (api.Object, error) {
	return execute[api.Object](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.Index(
			index,
			object.JSONReader(),
			c.ES.API.Index.WithContext(ctx),
			c.ES.API.Index.WithOpaqueID(logs.RequestID(ctx)),
			c.ES.API.Index.WithDocumentID(id),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to index es document %s to index %s: %w", id, index, err)
		}

		return resp, nil
	})
}

func (c *Client) DocumentDelete(ctx context.Context, index string, id string) error {
	_, err := execute[any](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.Delete(
			index,
			id,
			c.ES.API.Delete.WithContext(ctx),
			c.ES.API.Delete.WithOpaqueID(logs.RequestID(ctx)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to delete es document %s from index %s: %w", id, index, err)
		}

		return resp, nil
	})

	return err
}

func (c *Client) DocumentDeleteByQuery(ctx context.Context, index []string, s search.Search) error {
	_, err := execute[any](ctx, func() (*esapi.Response, error) {
		resp, err := c.ES.API.DeleteByQuery(
			collections.Map(index, c.indexNameMapper.Map),
			esutil.NewJSONReader(s),
			c.ES.API.DeleteByQuery.WithContext(ctx),
			c.ES.API.DeleteByQuery.WithOpaqueID(logs.RequestID(ctx)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to delete es document by query index %s: %w", index, err)
		}

		return resp, nil
	})

	return err
}

func NewClient(ctx context.Context, cfg *Config, debug bool) (*Client, error) {
	var esLogger elastictransport.Logger

	logger := logs.FromContext(ctx)

	caCert, err := os.ReadFile(cfg.Cert)
	if err != nil {
		return nil, fmt.Errorf("error while reading ES certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	if debug {
		esLogger = &elastictransport.ColorLogger{
			Output:             logger,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}
	}

	esConfig := elasticsearch.Config{
		Addresses: cfg.Hosts,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ //nolint
				RootCAs: caCertPool,
			},
			MaxIdleConnsPerHost:   cfg.MaxIdleConnectsPerHost,
			ResponseHeaderTimeout: CfgDefaultTimeOut,
		},
		Username: cfg.Username,
		Password: cfg.Password,
		RetryBackoff: func(attempt int) time.Duration {
			return time.Duration(attempt) * CfgDefaultTimeOut
		},
		EnableDebugLogger: debug,
		Logger:            esLogger,
	}

	elastic, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic search cleent: %w", err)
	}

	return &Client{
		ES:              elastic,
		startIndexer:    NewIndexerFactory(cfg, elastic, debug),
		indexNameMapper: NewIndexNameMapper(cfg),
	}, nil
}

func (m IndexNameMapper) Map(name string) string {
	if m == nil {
		return name
	}

	return m(name)
}

func NewIndexNameMapper(cfg *Config) IndexNameMapper {
	return func(name string) string {
		var builder strings.Builder

		if cfg.IndexPrefix != "" {
			builder.WriteString(cfg.IndexPrefix)
		}

		if cfg.Environment != "" {
			if builder.Len() != 0 {
				builder.WriteRune('-')
			}
			builder.WriteString(cfg.Environment)
		}

		if builder.Len() != 0 {
			builder.WriteRune('-')
		}

		builder.WriteString(name)

		return builder.String()
	}
}

func execute[T any](ctx context.Context, executor func() (*esapi.Response, error)) (T, error) {
	var blank T

	esResponse, err := executor()
	if err != nil {
		return blank, fmt.Errorf("failed to invoke es: %w", err)
	}

	defer func() {
		if err := esResponse.Body.Close(); err != nil {
			logs.FromContext(ctx).Err(err).Msg("failed to close es response body")
		}
	}()

	if esResponse.HasWarnings() {
		logs.FromContext(ctx).Warn().Msgf("es warns: %v", strings.Join(esResponse.Warnings(), "; "))
	}

	if esResponse.IsError() {
		if esResponse.StatusCode == http.StatusNotFound {
			return blank, ErrESNotFound
		}

		return blank, fmt.Errorf("%w: status %s", ErrESErrorResponse, esResponse.Status())
	}

	return serial.JSONDecode[T](esResponse.Body)
}
