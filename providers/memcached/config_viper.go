package memcached

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return NewConfig(
		getHostsList(loader, keyMapping...),
		viperx.Get(loader, CfgKeyMemcachedMaxBulkRequestSize.Map(keyMapping...), CfgDefaultMaxBulkRequestSize),
	)
}

func getHostsList(loader *viper.Viper, keyMapping ...cfg.KeyMap) []string {
	hosts := loader.GetStringSlice(string(CfgKeyMemcachedHosts.Map(keyMapping...)))
	if len(hosts) == 0 {
		hosts = CfgDefaultHosts
	}

	return hosts
}
