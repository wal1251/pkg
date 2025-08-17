// Package jwt предоставляет стратегию аутентификации JWT. Содержит функции для генерации и парсинга токенов JWT.
package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/core/security"
)

var _ security.AuthenticationProvider[security.BearerToken] = (*AuthenticationProvider)(nil)

var (
	ErrInvalidToken       = errors.New("invalid token")        // Некорректный токен.
	ErrInvalidTokenClaims = errors.New("invalid token claims") // Некорректный токен.
)

type (
	// ClaimsWithID интерфейс-обертка над jwt.Claims для получения ID токена.
	ClaimsWithID interface {
		jwt.Claims
		GetID() string
	}

	// RichToken основной интерфейс токена, который будет передаваться в приложение.
	// Содержит метаинформацию о токене, данные аутентификации и полномочия. Расширяется дополнительными методами для кастомных клеймов.
	RichToken interface {
		ForAccess(auth security.Authentication, ttl time.Duration, audience, issuer string, extra ...any) // Заполнение полей и клеймов для создания access
		ForRefresh(tokenID uuid.UUID, ttl time.Duration, audience, issuer string)                         // Заполнение полей и клеймов для создания refresh

		SetFromBearer(jwt.Token) error // Наполнение данными из jwt по интерфейсу jwt.Claims.
		SetBearerToken(string)         // Установка bearer при создании.

		GetID() uuid.UUID        // Идентификатор токена.
		GetSubjectID() uuid.UUID // Предмет токена (пользователь, другой токен).
		GetExpiresAt() time.Time // Срок жизни токена. Годен до.
		GetIssuedAt() time.Time  // Время выпуска токена.
		GetIssuer() string       // Сервис/приложение, выпустившее токен.
		GetAudience() string     // Сервис/приложение, для которого токен предназначен.
		GetName() string         // Имя владельца токена.
		GetPhoneNumber() string  // Телефон владельца токена.

		GetSessionID() string // Идентификатор сессии.

		GetAuthorities() []security.Authority // Данные о полномочиях.
		GetBearerToken() security.BearerToken // Строка токена на предъявителя.
		GetClaims() jwt.Claims
	}

	// TokenIssuer отвечает за выпуск токенов.
	TokenIssuer interface {
		GetSigningSecret() any                                                           // Получение секрета для подписи токена
		CreateAccessToken(auth security.Authentication, extra ...any) (RichToken, error) // Выпускает токен доступа.
		CreateRefreshToken(tokenID uuid.UUID) (RichToken, error)                         // Выпускает токен продления для токена доступа.
	}

	// TokenParser отвечает за парсинг и проверку токенов.
	TokenParser interface {
		GetParsingSecret() any                                    // Получение секрета для подписи токена
		ParseToken(token security.BearerToken) (RichToken, error) // Парсит и возвращает метаинформацию токена доступа.
	}

	// Manager позволяет выпускать новые токены и парсить ранее выпущенные.
	Manager interface {
		TokenIssuer
		TokenParser
	}

	// AuthenticationProvider обертка для использования TokenParser в качестве провайдера аутентификации.
	AuthenticationProvider struct {
		TokenParser
	}
)

// Authenticate см. AuthenticationProvider.Authenticate().
func (p *AuthenticationProvider) Authenticate(_ context.Context, t security.BearerToken) (security.Authentication, error) {
	token, err := p.TokenParser.ParseToken(t)
	if err != nil {
		return security.Authentication{}, errs.With(errs.ErrAuthFailure, err)
	}

	if token.GetName() == "" || token.GetSubjectID() == uuid.Nil {
		return security.Authentication{}, errs.Wrapf(errs.ErrAuthFailure, "token has no authentication data")
	}
	auth := security.Authentication{
		User: security.User{
			ID:          token.GetSubjectID(),
			Name:        token.GetName(),
			PhoneNumber: token.GetPhoneNumber(),
		},
		TokenID:     token.GetID(),
		Authorities: token.GetAuthorities(),
		SessionID:   token.GetSessionID(),
	}

	return auth, nil
}

// NewAuthenticationProvider провайдер аутентификации на базе Manager.
func NewAuthenticationProvider(m Manager) *AuthenticationProvider {
	return &AuthenticationProvider{m}
}
