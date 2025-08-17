package es

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	loader.SetDefault(CfgKeyEsHosts.String(), CfgDefaultHost)

	return &Config{
		Hosts:                  loader.GetStringSlice(CfgKeyEsHosts.String()),
		Timeout:                viperx.Get(loader, CfgKeyEsTimeout.Map(keyMapping...), CfgDefaultTimeOut),
		MaxIdleConnectsPerHost: viperx.Get(loader, CfgKeyEsMaxIdleConnectsPerHost.Map(keyMapping...), CfgMaxIdleConnectsPerHost),
		Cert:                   viperx.Get(loader, CfgKeyEsCertPath.Map(keyMapping...), CfgDefaultEsCertPath),
		Username:               viperx.Get(loader, CfgKeyEsUsername.Map(keyMapping...), CfgDefaultUsername),
		Password:               viperx.Get(loader, CfgKeyEsPassword.Map(keyMapping...), CfgDefaultPassword),
		IndexPrefix:            viperx.Get(loader, CfgKeyIndexPrefix.Map(keyMapping...), CfgDefaultIndexPrefix),
		Environment:            viperx.Environment(loader).String(),
		IndexerFlushInterval:   viperx.Get(loader, CfgKeyEsIndexerFlushInterval.Map(keyMapping...), CfgDefaultIndexerFlushInterval),
	}
}
