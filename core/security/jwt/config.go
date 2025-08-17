package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/wal1251/pkg/core/cfg"
)

const (
	CfgKeySecret          cfg.Key = "JWT_SECRET"            // Секрет для подписывания токенов JWT (string).
	CfgKeyPublicKeyPath   cfg.Key = "JWT_PUBLIC_KEY_PATH"   // Путь до публичного RSA ключа (используется для валидации JWT токенов).
	CfgKeyPrivateKeyPath  cfg.Key = "JWT_PRIVATE_KEY_PATH"  // Путь до приватного RSA ключа (используется для подписания JWT токенов).
	CfgKeyAccessTokenTTL  cfg.Key = "JWT_ACCESS_TOKEN_TTL"  // Время жизни токена доступа (duration).
	CfgKeyRefreshTokenTTL cfg.Key = "JWT_REFRESH_TOKEN_TTL" //nolint: gosec // Время жизни токена продления (duration).
	CfgAudience           cfg.Key = "JWT_CLAIM_AUDIENCE"    // Значение поля claim audience
	CfgIssuer             cfg.Key = "JWT_CLAIM_ISSUER"      // Значение поля claim issuer
	CfgSigningMethodName  cfg.Key = "JWT_SIGNING_METHOD_NAME"

	CfgDefaultAccessTokenTTL  = 7 * 24 * time.Hour           // Время жизни токена доступа по умолчанию.
	CfgDefaultRefreshTokenTTL = 4 * CfgDefaultAccessTokenTTL // Время жизни токена продления по умолчанию.
)

// ManagerConfig конфигурация менеджера JWT токенов.
type ManagerConfig struct {
	SigningMethod  jwt.SigningMethod
	PrivateKeyPath string // Путь до приватного ключа, которым подписываются JWT токены.
	PublicKeyPath  string // Путь до публичного ключа, для проверки JWT токенов.
	Secret         []byte // Секрет для подписывания токенов JWT.

	Audience        string        // Claim Audience - для кого токен выпущен
	Issuer          string        // Claim Issuer - для кого токен выпущен
	AccessTokenTTL  time.Duration // Время жизни токена доступа.
	RefreshTokenTTL time.Duration // Время жизни токена продления.
}
