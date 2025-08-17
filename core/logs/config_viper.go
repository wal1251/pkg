package logs

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

// CfgFromViper загружает конфиг с помощью viper.
func CfgFromViper(v *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		Level:  viperx.Get(v, CfgKeyLevel.Map(keyMapping...), CfgDefaultLevel),
		Pretty: viperx.Get(v, CfgKeyPretty.Map(keyMapping...), false),
	}
}
