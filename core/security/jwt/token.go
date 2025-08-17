package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/gen"
	"github.com/wal1251/pkg/core/security"
)

// DefaultClaims стандартные клеймы, которые помимо основных используют кастомный клейм с именем и полномочиями.
type DefaultClaims struct {
	Name        string               `json:"prl"`
	Authorities security.Authorities `json:"aus"`
	jwt.RegisteredClaims
}

// DefaultRichToken обогащенный метаинформацией токен. Знает, как обрабатывать дополнительные клеймы.
type DefaultRichToken struct {
	ID        uuid.UUID // Идентификатор токена.
	SubjectID uuid.UUID // Предмет токена (пользователь, другой токен).
	ExpiresAt time.Time // Срок жизни токена. Годен до.
	IssuedAt  time.Time // Время выпуска токена.
	Issuer    string    // Сервис/приложение, выпустившее токен.
	Audience  string    // Сервис/приложение, для которого токен предназначен.

	Claims      *DefaultClaims
	BearerToken security.BearerToken // Строка токена на предъявителя.
}

func (r *DefaultRichToken) GetPhoneNumber() string {
	// реализовано для имплементации интерфейса
	return r.Claims.Name
}

func (r *DefaultRichToken) GetSessionID() string {
	// реализовано для имплементации интерфейса
	return r.SubjectID.String()
}

func (c *DefaultClaims) GetID() string {
	return c.ID
}

func (r *DefaultRichToken) SetBearerToken(token string) {
	r.BearerToken = security.BearerToken(token)
}

func (r *DefaultRichToken) GetBearerToken() security.BearerToken {
	return r.BearerToken
}

func (r *DefaultRichToken) GetID() uuid.UUID {
	return r.ID
}

func (r *DefaultRichToken) GetSubjectID() uuid.UUID {
	return r.SubjectID
}

func (r *DefaultRichToken) GetExpiresAt() time.Time {
	return r.ExpiresAt
}

func (r *DefaultRichToken) GetIssuedAt() time.Time {
	return r.IssuedAt
}

func (r *DefaultRichToken) GetIssuer() string {
	return r.Issuer
}

func (r *DefaultRichToken) GetAudience() string {
	return r.Audience
}

// GetName возвращает нестандартный клейм Name на основе ClaimsImpl.
func (r *DefaultRichToken) GetName() string {
	return r.Claims.Name
}

// GetAuthorities возвращает нестандартный клейм Authorities на основе ClaimsImpl.
func (r *DefaultRichToken) GetAuthorities() []security.Authority {
	return r.Claims.Authorities
}

func (r *DefaultRichToken) GetClaims() jwt.Claims { return r.Claims }

func (r *DefaultRichToken) SetFromBearer(token jwt.Token) error {
	r.BearerToken = security.BearerToken(token.Raw)
	claims, ok := token.Claims.(ClaimsWithID)
	if !ok {
		return errors.New("invalid token claims, your CustomClaims must implement ClaimsWithID") //nolint:goerr113
	}

	if s, err := claims.GetSubject(); err == nil {
		sj, err := uuid.Parse(s)
		if err != nil {
			return ErrInvalidTokenClaims
		}
		r.SubjectID = sj
	} else {
		return ErrInvalidTokenClaims
	}

	if t, err := claims.GetExpirationTime(); err == nil {
		r.ExpiresAt = t.Time
	} else {
		return ErrInvalidTokenClaims
	}
	if t, err := claims.GetIssuedAt(); err == nil {
		r.IssuedAt = t.Time
	} else {
		return ErrInvalidTokenClaims
	}
	if s, err := claims.GetIssuer(); err == nil {
		r.Issuer = s
	} else {
		return ErrInvalidTokenClaims
	}
	if s, err := claims.GetAudience(); err == nil {
		m, err := s.MarshalJSON()
		if err != nil {
			return ErrInvalidTokenClaims
		}
		r.Audience = string(m)
	} else {
		return ErrInvalidTokenClaims
	}

	ID, err := uuid.Parse(claims.GetID())
	if err != nil {
		return ErrInvalidTokenClaims
	}
	r.ID = ID

	return nil
}

func (r *DefaultRichToken) ForAccess(auth security.Authentication, ttl time.Duration, audience, issuer string, _ ...any) {
	r.PopulateRegisteredClaims(auth.User.ID, ttl, audience, issuer)
	r.Claims = &DefaultClaims{
		RegisteredClaims: r.Claims.RegisteredClaims,
		Name:             auth.User.Name,
		Authorities:      auth.Authorities,
	}
}

func (r *DefaultRichToken) ForRefresh(tokenID uuid.UUID, ttl time.Duration, audience, issuer string) {
	r.PopulateRegisteredClaims(tokenID, ttl, audience, issuer)
}

func (r *DefaultRichToken) PopulateRegisteredClaims(subjectID uuid.UUID, ttl time.Duration, audience, issuer string) {
	id := gen.UUID().Next()
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(ttl)
	claims := jwt.RegisteredClaims{
		ID:        id.String(),
		Subject:   subjectID.String(),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	if len(audience) > 0 {
		claims.Audience = jwt.ClaimStrings{audience}
	}
	if len(issuer) > 0 {
		claims.Issuer = issuer
	}

	r.ID = id
	r.SubjectID = subjectID
	r.ExpiresAt = expiresAt
	r.IssuedAt = issuedAt
	r.Audience = audience
	r.Issuer = issuer
	r.Claims = &DefaultClaims{
		RegisteredClaims: claims,
	}
}
