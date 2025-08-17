package redis

import (
	"time"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/memorystore"
)

const (
	CfgDefaultHost               = "127.0.0.1"
	CfgDefaultPort               = "6379"
	CfgDefaultPassword           = ""
	CfgDefaultCertCa             = ""
	CfgDefaultDatabase           = 1
	CfgDefaultMasterName         = ""
	CfgDefaultClusterEnabled     = false
	CfgDefaultMaxBulkRequestSize = memorystore.DefaultMaxBulkRequestSize

	CfgKeyRedisHost               cfg.Key = "REDIS_HOST"
	CfgKeyRedisPort               cfg.Key = "REDIS_PORT"
	CfgKeyRedisPassword           cfg.Key = "REDIS_PASSWORD"
	CfgKeyRedisCertCa             cfg.Key = "REDIS_CERT_CA"
	CfgKeyRedisDatabase           cfg.Key = "REDIS_DATABASE"
	CfgKeyRedisMasterName         cfg.Key = "REDIS_MASTER_NAME"
	CfgKeyRedisClusterEnabled     cfg.Key = "REDIS_CLUSTER_ENABLED"
	CfgKeyRedisMaxBulkRequestSize cfg.Key = "REDIS_MAX_BULK_REQUEST_SIZE"

	BusCfgDefaultConsumerPrefetchLimit = 10
	BusCfgDefaultConsumerPollDuration  = time.Second

	BusCfgConsumerPrefetchLimit cfg.Key = "REDIS_BUS_CONSUMER_PREFETCH_LIMIT"
	BusCfgConsumerPollDuration  cfg.Key = "REDIS_BUS_CONSUMER_POLL_DURATION"

	redisFailedTxRetryInterval = time.Second
	RMQServiceName             = "RMQ"

	QueueConnectionCleanup = time.Minute
)

type (
	Config struct {
		Host               string
		Port               string
		Database           int
		Password           string
		CertCA             string // Локальный путь до сертификата.
		MasterName         string
		ClusterEnabled     bool
		MaxBulkRequestSize int
	}

	BusConfig struct {
		ConsumerPrefetchLimit int
		ConsumerPollDuration  time.Duration
	}
)

func NewConfig(
	host, port, password, certCa, masterName string,
	database int,
	clusterEnabled bool,
	bulkSize int,
) *Config {
	return &Config{
		Host:               host,
		Port:               port,
		Database:           database,
		Password:           password,
		CertCA:             certCa,
		MasterName:         masterName,
		ClusterEnabled:     clusterEnabled,
		MaxBulkRequestSize: bulkSize,
	}
}

func NewBusConfig(
	consumerPrefetchLimit int,
	consumerPollDuration time.Duration,
) *BusConfig {
	return &BusConfig{
		ConsumerPrefetchLimit: consumerPrefetchLimit,
		ConsumerPollDuration:  consumerPollDuration,
	}
}
