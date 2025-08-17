package ftp

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

// CfgFromViper возвращает конфигурацию клиента FTP, загруженную с помощью viper.
func CfgFromViper(v *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		Address:  viperx.Get(v, CfgKeyAddress.Map(keyMapping...), CfgDefaultAddress),
		User:     viperx.Get(v, CfgKeyUser.Map(keyMapping...), CfgDefaultUser),
		Password: viperx.Get(v, CfgKeyPassword.Map(keyMapping...), CfgDefaultPassword),
		Timeout:  viperx.Get(v, CfgKeyTimeout.Map(keyMapping...), CfgDefaultTimeout),
	}
}
