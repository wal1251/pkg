package grpcx

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(v *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		Host:                  viperx.Get(v, CfgKeyHost.Map(keyMapping...), CfgDefaultHost),
		Port:                  viperx.Get(v, CfgKeyPort.Map(keyMapping...), CfgDefaultPort),
		Reflection:            viperx.Get(v, CfgKeyReflection.Map(keyMapping...), CfgDefaultReflection),
		MaxConnectionIdle:     viperx.Get(v, CfgKeyMaxConnectionIdle.Map(keyMapping...), CfgDefaultMaxConnectionIdle),
		MaxConnectionAge:      viperx.Get(v, CfgKeyMaxConnectionAge.Map(keyMapping...), CfgDefaultMaxConnectionAge),
		MaxConnectionAgeGrace: viperx.Get(v, CfgKeyMaxConnectionAgeGrace.Map(keyMapping...), CfgDefaultMaxConnectionAgeGrace),
		KeepAliveTime:         viperx.Get(v, CfgKeyKeepAliveTime.Map(keyMapping...), CfgDefaultKeepAliveTime),
		KeepAliveTimeout:      viperx.Get(v, CfgKeyKeepAliveTimeout.Map(keyMapping...), CfgDefaultKeepAliveTimeout),
		MaxRecvMsgSize:        viperx.Get(v, CfgKeyMaxRecvMsgSize.Map(keyMapping...), CfgDefaultMaxRecvMsgSize),
		MaxSendMsgSize:        viperx.Get(v, CfgKeyMaxSendMsgSize.Map(keyMapping...), CfgDefaultMaxSendMsgSize),
		TLSCertFile:           viperx.Get(v, CfgKeyTLSCertFile.Map(keyMapping...), CfgDefaultTLSCertFile),
		TLSKeyFile:            viperx.Get(v, CfgKeyTLSKeyFile.Map(keyMapping...), CfgDefaultTLSKeyFile),
	}
}
