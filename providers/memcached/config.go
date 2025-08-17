package memcached

import (
	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/memorystore"
)

var CfgDefaultHosts = []string{"127.0.0.1:11211"} //nolint: gochecknoglobals

const (
	CfgDefaultMaxBulkRequestSize = memorystore.DefaultMaxBulkRequestSize

	CfgKeyMemcachedHosts              = cfg.Key("MEMCACHED_HOSTS")
	CfgKeyMemcachedMaxBulkRequestSize = cfg.Key("MEMCACHED_MAX_BULK_REQUEST_SIZE")
)

type Config struct {
	Hosts              []string
	MaxBulkRequestSize int
}

func NewConfig(hosts []string, bulkSize int) *Config {
	return &Config{
		Hosts:              hosts,
		MaxBulkRequestSize: bulkSize,
	}
}
