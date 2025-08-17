package mongo

import (
	"fmt"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"

	"github.com/spf13/viper"
)

const (
	CfgMongoUsername         cfg.Key = "MONGO_USERNAME"
	CfgMongoPassword         cfg.Key = "MONGO_PASSWORD"
	CfgMongoHost             cfg.Key = "MONGO_HOST"
	CfgMongoPort             cfg.Key = "MONGO_PORT"
	CfgMongoDatabase         cfg.Key = "MONGO_DATABASE"
	CfgMongoOptions          cfg.Key = "MONGO_OPTIONS"
	CfgMongoConnectionString cfg.Key = "MONGO_CONNECTION_STRING"

	defaultMongoConnectionString = "mongodb://%s:%s@%s:%s/%s?%s"
)

type (
	// Config представляет конфигурацию для подключения к MongoDB.
	Config struct {
		Username         string
		Password         string
		Host             string
		Port             string
		ConnectionString string
		Database         string
		Options          string
	}
)

// CfgFromViper загружает и возвращает конфигурацию подключения к MongoDB с помощью viper.
func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		Username:         viperx.Get(loader, CfgMongoUsername.Map(keyMapping...), ""),
		Password:         viperx.Get(loader, CfgMongoPassword.Map(keyMapping...), ""),
		Host:             viperx.Get(loader, CfgMongoHost.Map(keyMapping...), ""),
		Port:             viperx.Get(loader, CfgMongoPort.Map(keyMapping...), ""),
		Database:         viperx.Get(loader, CfgMongoDatabase.Map(keyMapping...), ""),
		Options:          viperx.Get(loader, CfgMongoOptions.Map(keyMapping...), ""),
		ConnectionString: viperx.Get(loader, CfgMongoConnectionString.Map(keyMapping...), defaultMongoConnectionString),
	}
}

func (c *Config) GetURI() string {
	return fmt.Sprintf(c.ConnectionString,
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Options) //nolint:nosprintfhostport
}
