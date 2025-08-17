package otelx

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(v *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		ExporterEndpoint:  viperx.Get(v, CfgKeyExporterEndpoint.Map(keyMapping...), CfgDefaultExporterEndpoint),
		ServiceName:       viperx.Get(v, CfgKeyServiceName.Map(keyMapping...), CfgDefaultServiceName),
		SamplerTraceRatio: viperx.Get(v, CfgKeySamplerTraceRatio.Map(keyMapping...), CfgDefaultSamplerTraceRatio),
		Env:               viperx.Environment(v).String(),
	}
}
