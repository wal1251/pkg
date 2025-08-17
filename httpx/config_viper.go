package httpx

import (
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func CfgFromViper(v *viper.Viper, keyMapping ...cfg.KeyMap) *Config {
	return &Config{
		Port:              viperx.Get(v, CfgKeyPort.Map(keyMapping...), CfgDefaultPort),
		ReadHeaderTimeout: CfgDefaultReadHeaderTimeout,
		Timeout:           viperx.Get(v, CfgKeyTimeout.Map(keyMapping...), CfgDefaultTimeout),
		CacheTTL:          viperx.Get(v, CfgKeyCacheTTL.Map(keyMapping...), CfgDefaultCacheTTL),
		BaseURL:           viperx.Get(v, CfgKeyBaseURL.Map(keyMapping...), CfgDefaultBaseURL),
		Rate: RateLimit{
			MaxRate: viperx.Get(v, CfgKeyMaxRate.Map(keyMapping...), CfgDefaultMaxRate),
			Every:   viperx.Get(v, CfgKeyRatePeriod.Map(keyMapping...), CfgDefaultRatePeriod),
		},
		CORS: CORSConfig{
			Origins:     viperx.Get(v, CfgKeyCORSOrigins.Map(keyMapping...), CfgDefaultCORSOrigins),
			Methods:     viperx.Get(v, CfgKeyCORSMethods.Map(keyMapping...), CfgDefaultCORSMethod),
			Headers:     viperx.Get(v, CfgKeyCORSHeaders.Map(keyMapping...), CfgDefaultCORSHeaders),
			Credentials: viperx.Get(v, CfgKeyCORSCredentials.Map(keyMapping...), false),
		},
		FrontMinVersion: viperx.Get(v, CfgKeyFrontMinVersion.Map(keyMapping...), ""),
	}
}
