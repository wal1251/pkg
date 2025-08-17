package otp

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(loader *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return NewConfig(
		viperx.Get(loader, CfgKeyOTPLength.Map(keyMapping...), CfgDefaultOTPLength),
		viperx.Get(loader, CfgKeyOTPLifetime.Map(keyMapping...), CfgDefaultOTPLifetime),
		viperx.Get(loader, CfgKeyOTPNextCodeDelay.Map(keyMapping...), CfgDefaultOTPNextCodeDelay),
		viperx.Get(loader, CfgKeyOTPMaxSendAttempts.Map(keyMapping...), CfgDefaultOTPMaxSendAttempts),
		viperx.Get(loader, CfgKeyOTPSendBlockDuration.Map(keyMapping...), CfgDefaultOTPSendBlockDuration),
		viperx.Get(loader, CfgKeyOTPMaxValidateAttempts.Map(keyMapping...), CfgDefaultOTPMaxValidateAttempts),
		viperx.Get(loader, CfgKeyOTPValidateBlockDuration.Map(keyMapping...), CfgDefaultOTPValidateBlockDuration),
	)
}
