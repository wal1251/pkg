// Package security предоставляет интерфейсы и базовые реализации объектов для обеспечения функций безопасности.
//
// Аспекты: аутентификация, авторизация разделение прав доступа и т.д...
//
// В общем случае для реализации авторизации пользователя алгоритм может выглядеть так:
//
// 1. В точке входа приложения (при получении запроса API, получении сообщения по подписке и т.д.), предположим, некоторый
// посредник middleware должен записать в контекст требуемые критерии безопасности (нужна ли авторизация, какими
// полномочиями должен обладать пользователь для исполнения операции), этот контекст далее передается для последующей
// обработки в нижележащие слои приложения.
//
// Например, так может выглядеть HTTP обработчик, устанавливающий параметры доступа к операции:
//
//	// FavouriteGroupGetList operation middleware
//	func (s *ServerInterfaceWrapper) GroupGetList(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		ctx = context.WithValue(ctx, BearerAuthScopes, []string{"admin"})
//
//		var handler = func(w http.ResponseWriter, r *http.Request) {
//			s.Handler.GroupGetList(w, r)
//		}
//
//		for _, middleware := range siw.HandlerMiddlewares {
//			handler = middleware(handler)
//		}
//
//		handler(w, r.WithContext(ctx))
//	}
//
// 2. Далее должен следовать слой (или посредник), отвечающий за аутентификацию и авторизацию, он вызывает реализацию
// интерфейса Manager, которая знает как интерпретировать параметры ранее установленные вышестоящим слоем,
// с помощью вызова Manager.SecurityRequirements().
//
// 3. Если требуется авторизация пользователя, тогда аутентификацию пользователя перед авторизацией нужно выполнить с
// помощью вызова реализации интерфейса AuthenticationProvider методом AuthenticationProvider.Authenticate(). В метод
// передаем извлеченные из запроса (или сообщения) данные для авторизации (например, извлечь из заголовка Authorization
// токен в случае HTTP запроса, для этого можно воспользоваться реализацией функции HTTPCredentialsProvider, см. httpx
// пакет).
//
// 4. После успешной аутентификации и проверки полномочий пользователя компонентом, отвечающим за авторизацию необходимо
// аутентификационные данные Authentication поместить в контекст приложения, чтобы они были доступны нижележащим слоям
// приложения с помощью той-же реализации Manager для получения данных об авторизованном пользователе
// вызовом метода Manager.Authorized().
//
// Конкретные реализации интерфейсов аутентификации и авторизации сильно зависят от протокола, поэтому реализации
// держим в пакетах с протоколом, например как для mw.Authorizer.
package security

import (
	"context"
	"net/http"

	"github.com/wal1251/pkg/core/ctxs"
)

const AuthorizedContextKey = "AUTH-Authorized" // Ключ хранения данных об авторизации.

var _ Manager = DefaultManager{}

type (
	// RequestCredentials формирует скоуп аутентификационных запросов.
	RequestCredentials interface {
		BearerToken | Credentials
	}

	// Manager менеджер авторизации. Ответит на вопросы: какие права нужны для доступа и кто сейчас авторизован.
	Manager interface {
		Authorized(context.Context) Authentication         // Возвращает данные аутентификации авторизованного пользователя.
		SecurityRequirements(context.Context) Requirements // Декларирует заявленные для операции полномочия.
	}

	// AuthenticationProvider провайдер аутентификации позволяет обслуживать запросы аутентификации пользователей.
	AuthenticationProvider[CRED RequestCredentials] interface {
		// Authenticate выполняет аутентификацию и идентификацию пользователя согласно запросу на аутентификацию.
		// Возвращает данные аутентификации, если аутентификация прошла успешно. В противном случае возвращает ошибку
		// типа errs.TypeAuthFailure.
		Authenticate(context.Context, CRED) (Authentication, error)
	}

	// HTTPCredentialsProvider извлекает аутентификационный запрос из HTTP запроса.
	HTTPCredentialsProvider[CRED RequestCredentials] func(r *http.Request) CRED

	// DefaultManager дефолтная реализация Manager.
	DefaultManager struct {
		AuthoritiesContextKey any
	}
)

// Authorized см. Manager.Authorized().
func (m DefaultManager) Authorized(ctx context.Context) Authentication {
	return ctxs.ValueGet[Authentication](ctx, AuthorizedContextKey)
}

// SecurityRequirements см. Manager.SecurityRequirements().
func (m DefaultManager) SecurityRequirements(ctx context.Context) Requirements {
	scopes := ctx.Value(m.AuthoritiesContextKey)
	if scopes == nil {
		return Requirements{}
	}

	switch requiredAuthorities := scopes.(type) {
	case []string:
		return Requirements{
			IsAuthorizationRequired: len(requiredAuthorities) != 0,
			Authorities:             AuthoritiesFromString(requiredAuthorities),
		}
	case Authorities:
		return Requirements{
			IsAuthorizationRequired: len(requiredAuthorities) != 0,
			Authorities:             requiredAuthorities,
		}
	default:
		return Requirements{}
	}
}
