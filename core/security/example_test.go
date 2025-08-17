package security_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/wal1251/pkg/core/security"
	"github.com/wal1251/pkg/core/security/jwt"
	"github.com/wal1251/pkg/httpx"
	"github.com/wal1251/pkg/httpx/mw"
)

func ExampleManager() {
	// Менеджер авторизации.
	manager := security.DefaultManager{AuthoritiesContextKey: "BearerAuthScopes"}

	// Менеджер JWT: выпускает и парсит токены доступа.
	jwtManager, _ := jwt.NewTokenManager(
		&jwt.ManagerConfig{
			Secret:         []byte("jwt-secret"),
			SigningMethod:  jwtlib.SigningMethodHS512,
			AccessTokenTTL: 1<<63 - 1,
			Audience:       "test-app",
		},
		func() jwt.RichToken { return &jwt.DefaultRichToken{Claims: &jwt.DefaultClaims{}} },
	)

	// Аутентификация на базе токенов JWT. Простейшая реализация: если токен валиден, то аутентификация пользователя
	// считается пройденной.
	authManager := jwt.NewAuthenticationProvider(jwtManager)

	// Посредник: устанавливает требования доступа к обработчику next.
	require := func(next http.Handler, authorities []string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Заявляем необходимые полномочия. В реальном приложении полномочия могут быть установлены отдельно для
			// каждого эндпоинта.
			ctx := context.WithValue(r.Context(), manager.AuthoritiesContextKey, authorities)

			// Передаем управление дальше.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	// Посредник: выполнит next, только если у пользователя есть полномочия.
	authorized := mw.Authorizer[security.BearerToken](
		manager,
		httpx.BearerTokenExtract,
		authManager)

	// Бизнес логика.
	greeting := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получим авторизованного пользователя.
		auth := manager.Authorized(r.Context())

		w.WriteHeader(http.StatusOK)

		// Поприветствуем пользователя.
		if _, err := w.Write([]byte(fmt.Sprintf("Hello, %s!", auth.User.Name))); err != nil {
			panic("can't write response")
		}
	})

	// Запустим сервер с защищенным эндпоинтом, только пользователь с полномочиями admin имеет доступ.
	server := httptest.NewServer(require(authorized(greeting), []string{"admin"}))
	defer server.Close()

	// Функция чтения запроса.
	readResponse := func(resp *http.Response) {
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(resp.StatusCode, " ", string(body))
	}

	// Запрос без данных для авторизации.
	respNoCreds, err := http.Get(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	// Этот запрос будет отклонен со статусом 401.
	readResponse(respNoCreds)

	// Создадим токен с полномочиями user.
	tokenUser, err := jwtManager.CreateAccessToken(security.Authentication{
		User:        security.User{Name: "zhupanovdm@github.com", ID: uuid.New()},
		Authorities: security.Authorities{"user"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Создадим запрос к серверу с данными для авторизации user.
	requestUser, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Отправим данные авторизации в заголовках запроса.
	httpx.BearerTokenSet(requestUser, tokenUser.GetBearerToken())

	respUser, err := http.DefaultClient.Do(requestUser)
	if err != nil {
		log.Fatal(err)
	}

	// Этот запрос будет отклонен со статусом 403, предоставленных полномочий недостаточно.
	readResponse(respUser)

	// Создадим токен с полномочиями user.
	tokenAdmin, err := jwtManager.CreateAccessToken(security.Authentication{
		User:        security.User{Name: "zhupanovdm@github.com", ID: uuid.New()},
		Authorities: security.Authorities{"admin"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Создадим запрос к серверу с данными для авторизации user.
	requestAdmin, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Отправим данные авторизации в заголовках запроса.
	httpx.BearerTokenSet(requestAdmin, tokenAdmin.GetBearerToken())

	respAdmin, err := http.DefaultClient.Do(requestAdmin)
	if err != nil {
		log.Fatal(err)
	}

	// Этот запрос выполнится успешно - авторизация успешна.
	readResponse(respAdmin)

	// Output:
	// 401 {"code":"AUTH_FAILURE","message":"AUTH_FAILURE: can't parse bearer token: token is malformed: token contains an invalid number of segments"}
	// 403 {"code":"FORBIDDEN","message":"FORBIDDEN: requested action is not permitted for current user"}
	// 200 Hello, zhupanovdm@github.com!
}
