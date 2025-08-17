package presenters

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	loader.SetDefault(CfgKeySecuredKeywords.String(), CfgDefaultKeySecuredKeyword)

	return &Config{
		SecuredKeywords: loader.GetStringSlice(CfgKeySecuredKeywords.String()),
		MaxStringLength: viperx.Get(loader, CfgKeyMaxStringLength.Map(keyMapping...), CfgDefaultMaxStringLength),
	}
}
