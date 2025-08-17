package es

import (
	"time"

	"github.com/wal1251/pkg/core/cfg"
)

const (
	CfgKeyEsHosts                  cfg.Key = "ES_HOSTS"                      // Список хостов, используемых ES.
	CfgKeyEsCertPath               cfg.Key = "ES_CERT_PATH"                  // Путь к сертификату.
	CfgKeyEsTimeout                cfg.Key = "ES_TIMEOUT"                    // Время задержки повторного запроса.
	CfgKeyEsMaxIdleConnectsPerHost cfg.Key = "ES_MAX_IDLE_CONNECTS_PER_HOST" // Максимальное количество незанятых подключений на один хост.
	CfgKeyEsUsername               cfg.Key = "ES_USERNAME"                   // Имя пользователя для HTTP Basic Authentication.
	CfgKeyEsPassword               cfg.Key = "ES_PASSWORD"                   // Пароль пользователя для HTTP Basic Authentication.
	CfgKeyEsIndexerFlushInterval   cfg.Key = "ES_INDEXER_FLUSH_INTERVAL"     // Время сброса.
	CfgKeyIndexPrefix              cfg.Key = "ES_PREFIX"

	CfgDefaultTimeOut              = 5 * time.Second
	CfgDefaultEsCertPath           = ""
	CfgDefaultIndexerFlushInterval = 5 * time.Second
	CfgDefaultHost                 = "https://localhost:9200"
	CfgDefaultUsername             = "elastic"
	CfgDefaultPassword             = "elastic"
	CfgMaxIdleConnectsPerHost      = 1
	CfgDefaultIndexPrefix          = "domain"
)

type (
	Config struct {
		Hosts                  []string
		Timeout                time.Duration
		MaxIdleConnectsPerHost int
		Cert                   string
		Username               string
		Password               string
		IndexPrefix            string
		Environment            string
		IndexerFlushInterval   time.Duration
	}
)
