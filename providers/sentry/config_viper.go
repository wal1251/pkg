package sentry

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		Toggle:      viperx.Get(loader, CfgKeySentryToggle.Map(keyMapping...), CfgDefaultToggle),
		DSN:         viperx.Get(loader, CfgKeySentryDSN.Map(keyMapping...), CfgDefaultDSN),
		Debug:       viperx.Get(loader, CfgKeySentryDebug.Map(keyMapping...), CfgDefaultDebug),
		Environment: viperx.Environment(loader).String(),
		Timeout:     viperx.Get(loader, CfgKeySentryTimeout.Map(keyMapping...), CfgDefaultTimeout),
	}
}
