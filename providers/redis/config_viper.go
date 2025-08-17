package redis

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return NewConfig(
		viperx.Get(loader, CfgKeyRedisHost.Map(keyMapping...), CfgDefaultHost),
		viperx.Get(loader, CfgKeyRedisPort.Map(keyMapping...), CfgDefaultPort),
		viperx.Get(loader, CfgKeyRedisPassword.Map(keyMapping...), CfgDefaultPassword),
		viperx.Get(loader, CfgKeyRedisCertCa.Map(keyMapping...), CfgDefaultCertCa),
		viperx.Get(loader, CfgKeyRedisMasterName.Map(keyMapping...), CfgDefaultMasterName),
		viperx.Get(loader, CfgKeyRedisDatabase.Map(keyMapping...), CfgDefaultDatabase),
		viperx.Get(loader, CfgKeyRedisClusterEnabled.Map(keyMapping...), CfgDefaultClusterEnabled),
		viperx.Get(loader, CfgKeyRedisMaxBulkRequestSize.Map(keyMapping...), CfgDefaultMaxBulkRequestSize),
	)
}

func BusCfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *BusConfig {
	return NewBusConfig(
		viperx.Get(loader, BusCfgConsumerPrefetchLimit.Map(keyMapping...), BusCfgDefaultConsumerPrefetchLimit),
		viperx.Get(loader, BusCfgConsumerPollDuration.Map(keyMapping...), BusCfgDefaultConsumerPollDuration),
	)
}
