package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

// ManagerCfgFromViper загрузка конфига ManagerConfig с помощью viper.
func ManagerCfgFromViper(v *viper.Viper, keyMapping ...cfg.KeyMap) *ManagerConfig {
	singingMethod := getSingingMethodByName(v)

	return &ManagerConfig{
		SigningMethod:   singingMethod,
		Secret:          []byte(viperx.Get(v, CfgKeySecret.Map(keyMapping...), "")),
		PrivateKeyPath:  viperx.Get(v, CfgKeyPrivateKeyPath.Map(keyMapping...), ""),
		PublicKeyPath:   viperx.Get(v, CfgKeyPublicKeyPath.Map(keyMapping...), ""),
		Audience:        viperx.Get(v, CfgAudience.Map(keyMapping...), ""),
		Issuer:          viperx.Get(v, CfgIssuer.Map(keyMapping...), ""),
		AccessTokenTTL:  viperx.Get(v, CfgKeyAccessTokenTTL.Map(keyMapping...), CfgDefaultAccessTokenTTL),
		RefreshTokenTTL: viperx.Get(v, CfgKeyRefreshTokenTTL.Map(keyMapping...), CfgDefaultRefreshTokenTTL),
	}
}

func getSingingMethodByName(v *viper.Viper) jwt.SigningMethod {
	singingMethodName := v.GetString(CfgSigningMethodName.String())
	switch singingMethodName {
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS512":
		return jwt.SigningMethodRS512
	case "RS384":
		return jwt.SigningMethodRS384
	case "HS512":
		return jwt.SigningMethodHS512
	case "ES256":
		return jwt.SigningMethodES256
	case "ES512":
		return jwt.SigningMethodES512
	default:
		return jwt.SigningMethodHS256
	}
}
