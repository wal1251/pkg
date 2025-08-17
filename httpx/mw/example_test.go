package mw_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/security"
	"github.com/wal1251/pkg/core/security/jwt"
	"github.com/wal1251/pkg/httpx"
	"github.com/wal1251/pkg/httpx/mw"
	"github.com/wal1251/pkg/tools/serial"
)

func ExampleAuthorizer_jwt() {
	key := struct{}{}
	jwtManager, err := jwt.NewTokenManager(
		&jwt.ManagerConfig{AccessTokenTTL: time.Hour, Secret: []byte("secret")},
		func() jwt.RichToken { return &jwt.DefaultRichToken{Claims: &jwt.DefaultClaims{}} },
	)
	if err != nil {
		log.Fatal(err)
	}
	manager := security.DefaultManager{AuthoritiesContextKey: key}
	authorizer := mw.Authorizer[security.BearerToken](
		manager,
		httpx.BearerTokenExtract,
		jwt.NewAuthenticationProvider(jwtManager),
	)

	// Обработчик ответит JSON'ом с авторизованным пользователем.
	endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.SendResponse(r.Context(), w,
			httpx.NewServerResponse[security.Authentication]().
				WithContentTypeJSON().
				WithValue(manager.Authorized(r.Context())),
		)
	})

	// Создадим запрос с данными для авторизации.
	token, err := jwtManager.CreateAccessToken(security.Authentication{
		User: security.User{
			ID:   uuid.MustParse("d00ba9a3-d8bc-431c-ab36-37cd72ae825c"),
			Name: "d.zhupanov@redmadrobot.com",
		},
		Authorities: security.Authorities{"exec"},
	})
	if err != nil {
		log.Fatalf("failed to create access token: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.GetBearerToken()))

	// Обработаем запрос клиента.
	wr := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Эмулируем установку полномочий, требуемых для доступа к методу.
		ctx := context.WithValue(r.Context(), key, []string{"exec"})

		// Защитим наш обработчик с помощью посредника авторизации.
		authorizer(endpointHandler).ServeHTTP(w, r.WithContext(ctx))
	}).ServeHTTP(wr, r)

	// Прочитаем ответ сервера.
	a, err := serial.JSONDecode[security.Authentication](wr.Body)
	if err != nil {
		log.Fatalf("failed to decode server response: %v", err)
	}

	fmt.Println("user id:", a.User.ID)
	fmt.Println("user name:", a.User.Name)
	fmt.Println("authorities:", a.Authorities)

	// Output:
	// user id: d00ba9a3-d8bc-431c-ab36-37cd72ae825c
	// user name: d.zhupanov@redmadrobot.com
	// authorities: [exec]
}
