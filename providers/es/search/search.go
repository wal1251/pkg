package search

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/elastic/go-elasticsearch/v8/esutil"

	"github.com/wal1251/pkg/providers/es/api"
)

const (
	FuzzinessAuto = "AUTO"
)

type (
	Query interface {
		json.Marshaler
		QueryType() string
		IsEmpty() bool
	}

	Queries []Query

	Aggregation interface {
		json.Marshaler
		AggregationType() string
	}

	Search struct {
		From  *int                   `json:"from,omitempty"`
		Size  *int                   `json:"size,omitempty"`
		Query Query                  `json:"query,omitempty"`
		Aggs  map[string]Aggregation `json:"aggs,omitempty"`
	}

	MultiSearchHeader struct {
		Index string `json:"index,omitempty"`
	}

	MultiSearchItem struct {
		Header MultiSearchHeader
		Search Search
	}

	MultiSearch []MultiSearchItem

	MultiResponse struct {
		Took      int64      `json:"took"`
		Responses []Response `json:"responses"`
	}

	Response struct {
		Status       int                        `json:"status"`
		Took         int64                      `json:"took"`
		TimedOut     bool                       `json:"timed_out"` //nolint: tagliatelle
		Hits         Hits                       `json:"hits"`
		Aggregations map[string]json.RawMessage `json:"aggregations"`
	}

	Hits struct {
		Total    HitsTotal `json:"total"`
		MaxScore float64   `json:"max_score"` //nolint: tagliatelle
		Hits     []Hit     `json:"hits"`
	}

	HitsTotal struct {
		Value    int64 `json:"value"`
		Relation string
	}

	Hit struct {
		ID     string          `json:"_id"`     //nolint: tagliatelle
		Index  string          `json:"_index"`  //nolint: tagliatelle
		Score  float64         `json:"_score"`  //nolint: tagliatelle
		Source json.RawMessage `json:"_source"` //nolint: tagliatelle
	}
)

func (s MultiSearch) NewNDJSONReader() io.Reader {
	buff := bytes.NewBuffer([]byte{})

	for i := range s {
		if _, err := io.Copy(buff, esutil.NewJSONReader(s[i].Header)); err != nil {
			return api.NewErrReader(err)
		}

		if _, err := io.Copy(buff, esutil.NewJSONReader(s[i].Search)); err != nil {
			return api.NewErrReader(err)
		}
	}

	return buff
}

func Aggregations(options ...func(map[string]Aggregation)) map[string]Aggregation {
	var aggregations map[string]Aggregation

	for _, opt := range options {
		if aggregations == nil {
			aggregations = make(map[string]Aggregation)
		}

		opt(aggregations)
	}

	return aggregations
}
