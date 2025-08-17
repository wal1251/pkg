package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/security"
)

var _ Manager = (*TokenManager)(nil)

// TokenManager представляет собой структуру, реализующую Manager и
// используемую для управления созданием и проверкой JWT-токенов.
type TokenManager struct {
	SigningMethod jwt.SigningMethod // Конкретная версия алгоритма подписи
	Secret        []byte            // Секрет для подписывания токенов методами HS.
	PrivateKey    *rsa.PrivateKey   // Секрет для подписывания токенов методами RS.
	PublicKey     *rsa.PublicKey    // Секрет для подписывания токенов методами RS.

	Audience        string        // Claim Audience - для кого токен выпущен
	Issuer          string        // Claim Issuer - для кого токен выпущен
	AccessTokenTTL  time.Duration // Время жизни токенов доступа.
	RefreshTokenTTL time.Duration // Время жизни токенов продления.

	NewRichToken func() RichToken
}

// GetSigningSecret получает секрет для подписи в зависимости от алгоритма из конфигурации.
func (m *TokenManager) GetSigningSecret() any {
	if _, ok := m.SigningMethod.(*jwt.SigningMethodHMAC); ok {
		return m.Secret
	}

	return m.PrivateKey
}

// GetParsingSecret получает секрет для проверки подписи в зависимости от алгоритма из конфигурации.
func (m *TokenManager) GetParsingSecret() any {
	if _, ok := m.SigningMethod.(*jwt.SigningMethodHMAC); ok {
		return m.Secret
	}

	return m.PublicKey
}

// CreateAccessToken см. Manager.CreateAccessToken().
func (m *TokenManager) CreateAccessToken(auth security.Authentication, extra ...any) (RichToken, error) {
	richToken := m.NewRichToken()
	richToken.ForAccess(auth, m.AccessTokenTTL, m.Audience, m.Issuer, extra...)

	token := jwt.NewWithClaims(m.SigningMethod, richToken.GetClaims())
	jwtToken, err := token.SignedString(m.GetSigningSecret())
	if err != nil {
		return nil, fmt.Errorf("can't create jwt: %w", err)
	}
	richToken.SetBearerToken(jwtToken)

	return richToken, nil
}

// CreateRefreshToken см. Manager.CreateRefreshToken().
func (m *TokenManager) CreateRefreshToken(tokenID uuid.UUID) (RichToken, error) {
	richToken := m.NewRichToken()
	richToken.ForRefresh(tokenID, m.RefreshTokenTTL, m.Audience, m.Issuer)

	token := jwt.NewWithClaims(m.SigningMethod, richToken.GetClaims())
	jwtToken, err := token.SignedString(m.GetSigningSecret())
	if err != nil {
		return nil, fmt.Errorf("can't create jwt: %w", err)
	}
	richToken.SetBearerToken(jwtToken)

	return richToken, nil
}

// ParseToken см. Manager.ParseToken().
func (m *TokenManager) ParseToken(jwtToken security.BearerToken) (RichToken, error) {
	richToken := m.NewRichToken()

	token, err := jwt.ParseWithClaims(string(jwtToken), richToken.GetClaims(), func(*jwt.Token) (any, error) { return m.GetParsingSecret(), nil })
	if err != nil {
		return nil, fmt.Errorf("can't parse bearer token: %w", err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	err = richToken.SetFromBearer(*token)
	if err != nil {
		return nil, err
	}

	return richToken, nil
}

// NewTokenManager создает новый экземпляр TokenManager с заданной конфигурацией.
func NewTokenManager(cfg *ManagerConfig, newRichToken func() RichToken) (*TokenManager, error) {
	manager := TokenManager{
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
		Secret:          cfg.Secret,
		SigningMethod:   jwt.SigningMethodHS256,
		Audience:        cfg.Audience,
		Issuer:          cfg.Issuer,
		NewRichToken:    newRichToken,
	}
	if cfg.SigningMethod != nil {
		manager.SigningMethod = cfg.SigningMethod
	} else {
		manager.SigningMethod = jwt.SigningMethodHS256
	}

	if len(cfg.PublicKeyPath) > 0 {
		publicKeyBytes, err := os.ReadFile(cfg.PublicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read public key from path: %s. %w", cfg.PublicKeyPath, err)
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}
		manager.PublicKey = publicKey
	}
	if len(cfg.PrivateKeyPath) > 0 {
		privateKeyBytes, err := os.ReadFile(cfg.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key from path: %s. %w", cfg.PublicKeyPath, err)
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		manager.PrivateKey = privateKey
	}
	if _, ok := manager.SigningMethod.(*jwt.SigningMethodHMAC); ok && len(manager.Secret) == 0 {
		return nil, errors.New("HMAC secret empty") //nolint:goerr113
	}

	return &manager, nil
}
