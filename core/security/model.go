package security

import (
	"context"

	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/ctxs"
	"github.com/wal1251/pkg/core/presenters"
	"github.com/wal1251/pkg/tools/collections"
)

const displayTokenMaxLength = 5

var (
	_ presenters.StringViewer = Password("")
	_ presenters.StringViewer = BearerToken("")
)

type (
	// Requirements требования безопасности: нужна ли авторизация, полномочия.
	Requirements struct {
		IsAuthorizationRequired bool        // Нужна авторизация?
		Authorities             Authorities // Требуемые полномочия.
	}

	BearerToken string // Токен на предъявителя.
	Password    string // Пароль.

	// Credentials учетные данные для аутентификации.
	Credentials struct {
		Login    string   // Логин (имя учетной записи).
		Password Password // Пароль.
	}

	// User авторизованное лицо.
	User struct { //
		ID          uuid.UUID `json:"id"`           // Идентификатор записи.
		Name        string    `json:"name"`         // Имя учетной записи.
		PhoneNumber string    `json:"phone_number"` //nolint:tagliatelle // Номер телефона
	}

	// Authority полномочие пользователя на выполнение операции или группы операций.
	Authority string

	// Authorities набор полномочий.
	Authorities []Authority

	// Authentication данные об аутентификации.
	Authentication struct {
		User        User        // Авторизованное лицо.
		TokenID     uuid.UUID   // Идентификатор токена доступа.
		SessionID   string      // Идентификатор сессии.
		Authorities Authorities `json:"authorities,omitempty"` // Полномочия авторизованного пользователя.
	}
)

func (p Password) StringView(view presenters.ViewType, _ presenters.ViewOptions) string {
	if view == presenters.ViewLogs {
		return presenters.DefaultCredentialsPlaceholder
	}

	return string(p)
}

func (t BearerToken) StringView(view presenters.ViewType, _ presenters.ViewOptions) string {
	if view == presenters.ViewLogs {
		return presenters.StringTail(string(t), displayTokenMaxLength)
	}

	return string(t)
}

// Meets возвращает true, полномочия a удовлетворяют запрошенным полномочиям required.
func (a Authorities) Meets(required Authorities) bool {
	return len(required) == 0 || collections.NewSet[Authority](a...).ContainsAny(required...)
}

// ToStrings возвращает набор полномочий преобразованный в слайс строк.
func (a Authorities) ToStrings() []string {
	return collections.Map(a, func(t Authority) string { return string(t) })
}

// ToContext поместить полномочия в контекст.
func (a Authentication) ToContext(ctx context.Context) context.Context {
	return ctxs.ValuePut(ctx, AuthorizedContextKey, a)
}

// AuthoritiesFromString возвращает набор полномочий из слайса строк.
func AuthoritiesFromString(s []string) Authorities {
	return collections.Map(s, func(t string) Authority { return Authority(t) })
}
