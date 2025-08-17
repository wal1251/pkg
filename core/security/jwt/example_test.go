package jwt_test

import (
	"context"
	"fmt"
	"log"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/wal1251/pkg/core/gen"

	"github.com/google/uuid"

	"github.com/wal1251/pkg/core/security"
	"github.com/wal1251/pkg/core/security/jwt"
	"github.com/wal1251/pkg/tools/collections"
)

func ExampleHS256Manager_accessToken() {
	// Новый менеджер JWT, который выпускает практически бессрочные токены.
	manager, err := jwt.NewTokenManager(
		&jwt.ManagerConfig{
			AccessTokenTTL: 1<<63 - 1,
			Secret:         []byte("jwt-secret"),
			SigningMethod:  jwtlib.SigningMethodHS256,
		},
		func() jwt.RichToken { return &jwt.DefaultRichToken{Claims: &jwt.DefaultClaims{}} },
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Новый токен доступа для пользователя с указанными полномочиями.
	token, err := manager.CreateAccessToken(security.Authentication{
		User:        security.User{ID: uuid.MustParse("d00ba9a3-d8bc-431c-ab36-37cd72ae825c")},
		Authorities: collections.Single[security.Authority]("user"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Парсим выпущенный токен.
	jwtToken, err := manager.ParseToken(token.GetBearerToken())
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("user id:", jwtToken.GetSubjectID())
	fmt.Println("authorities:", jwtToken.GetAuthorities())

	// Output:
	// user id: d00ba9a3-d8bc-431c-ab36-37cd72ae825c
	// authorities: [user]
}

func ExampleHS256Manager_refreshToken() {
	manager, err := jwt.NewTokenManager(
		&jwt.ManagerConfig{
			RefreshTokenTTL: 1<<63 - 1,
			Secret:          []byte("jwt-secret"),
			SigningMethod:   jwtlib.SigningMethodHS512,
		},
		func() jwt.RichToken { return &jwt.DefaultRichToken{Claims: &jwt.DefaultClaims{}} },
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	token, err := manager.CreateRefreshToken(uuid.MustParse("09b5fdac-ca0d-4fa2-bbcb-b9a9958d9904"))
	if err != nil {
		log.Fatal(err)
		return
	}

	jwtToken, err := manager.ParseToken(token.GetBearerToken())
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("subject id:", jwtToken.GetSubjectID())

	// Output:
	// subject id: 09b5fdac-ca0d-4fa2-bbcb-b9a9958d9904
}

func ExampleRS256Manager_accessToken() {
	// Новый менеджер JWT, который выпускает практически бессрочные токены.
	manager, err := jwt.NewTokenManager(&jwt.ManagerConfig{
		AccessTokenTTL: 1<<63 - 1,
		PrivateKeyPath: "testdata/private_key",
		PublicKeyPath:  "testdata/public_key",
		SigningMethod:  jwtlib.SigningMethodRS256,
	},
		func() jwt.RichToken { return &jwt.DefaultRichToken{Claims: &jwt.DefaultClaims{}} },
	)
	if err != nil {
		log.Fatal(err)
	}

	// Новый токен доступа для пользователя с указанными полномочиями.
	token, err := manager.CreateAccessToken(security.Authentication{
		User:        security.User{ID: uuid.MustParse("f6c82fb6-0f41-44d8-ae00-90066d4f39c8")},
		Authorities: collections.Single[security.Authority]("user"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Парсим выпущенный токен.
	jwtToken, err := manager.ParseToken(token.GetBearerToken())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user id:", jwtToken.GetSubjectID())
	fmt.Println("authorities:", jwtToken.GetAuthorities())

	// Output:
	// user id: f6c82fb6-0f41-44d8-ae00-90066d4f39c8
	// authorities: [user]
}

func ExampleRS256Manager_refreshToken() {
	// Новый менеджер JWT, который выпускает практически бессрочные токены.
	manager, err := jwt.NewTokenManager(&jwt.ManagerConfig{
		RefreshTokenTTL: 1<<63 - 1,
		PrivateKeyPath:  "testdata/private_key",
		PublicKeyPath:   "testdata/public_key",
		SigningMethod:   jwtlib.SigningMethodRS256,
	},
		func() jwt.RichToken { return &jwt.DefaultRichToken{Claims: &jwt.DefaultClaims{}} },
	)
	if err != nil {
		log.Fatal(err)
	}

	// Новый токен продления.
	token, err := manager.CreateRefreshToken(uuid.MustParse("4bbc07bc-d340-4148-ba0b-a32763155a61"))
	if err != nil {
		log.Fatal(err)
	}

	// Парсим выпущенный токен.
	jwtToken, err := manager.ParseToken(token.GetBearerToken())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("subject id:", jwtToken.GetSubjectID())

	// Output:
	// subject id: 4bbc07bc-d340-4148-ba0b-a32763155a61
}

type AppClaims struct {
	jwtlib.RegisteredClaims
	NewField    string               `json:"new_field"`
	Name        string               `json:"name"`
	Authorities security.Authorities `json:"roles"`
}

func (c *AppClaims) GetID() string {
	return c.ID
}

type AppRichToken struct {
	jwt.DefaultRichToken
	Claims *AppClaims
}

func (r *AppRichToken) GetNewField() string {
	r.GetClaims()
	return r.Claims.NewField
}

func (r *AppRichToken) GetAuthorities() []security.Authority {
	return r.Claims.Authorities
}

func (r *AppRichToken) GetName() string {
	return r.Claims.Name
}

func (r *AppRichToken) GetClaims() jwtlib.Claims { return r.Claims }

func (r *AppRichToken) ForAccess(auth security.Authentication, ttl time.Duration, audience, issuer string, extra ...any) {
	r.PopulateRegisteredClaims(auth.User.ID, ttl, audience, issuer)
	r.Claims = &AppClaims{
		RegisteredClaims: r.DefaultRichToken.Claims.RegisteredClaims,
		Name:             extra[0].(string),
		NewField:         extra[1].(string),
		Authorities:      auth.Authorities,
	}
}

func ExampleManagerWithCustomClaims() {
	// Пример при использовании полей jwt, отличных от дефолтных, и дополнительных полей

	// Менеджер JWT: выпускает и парсит токены доступа.
	jwtManager, _ := jwt.NewTokenManager(
		&jwt.ManagerConfig{
			Secret:         []byte("jwt-secret"),
			AccessTokenTTL: 1<<63 - 1,
		},
		func() jwt.RichToken { return &AppRichToken{Claims: &AppClaims{}} },
	)

	// Создадим токен с полномочиями user.
	t, err := jwtManager.CreateAccessToken(security.Authentication{
		User:        security.User{ID: gen.UUID().Next(), Name: "alex@black.com"},
		Authorities: security.Authorities{"active"},
	}, "robot", "it works")
	if err != nil {
		log.Fatal(err)
	}

	r, err := jwtManager.ParseToken(t.GetBearerToken())
	if err != nil {
		log.Fatal(err)
	}
	rich := r.(*AppRichToken)

	fmt.Println(rich.GetName())
	fmt.Println(rich.GetAuthorities())
	fmt.Println(rich.GetNewField())

	// Output:
	// robot
	// [active]
	// it works
}

func ExampleAuthenticationProvider() {
	manager, err := jwt.NewTokenManager(
		&jwt.ManagerConfig{AccessTokenTTL: 1<<63 - 1, Secret: []byte("jwt-secret")},
		func() jwt.RichToken { return &jwt.DefaultRichToken{Claims: &jwt.DefaultClaims{}} },
	)
	if err != nil {
		log.Fatal(err)
	}

	token, err := manager.CreateAccessToken(security.Authentication{
		User: security.User{
			ID:   uuid.MustParse("d00ba9a3-d8bc-431c-ab36-37cd72ae825c"),
			Name: "d.zhupanov@redmadrobot.com",
		},
		Authorities: collections.Single[security.Authority]("user"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	authentication, err := jwt.NewAuthenticationProvider(manager).Authenticate(context.TODO(), token.GetBearerToken())
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("user id:", authentication.User.ID)
	fmt.Println("user name:", authentication.User.Name)
	fmt.Println("authorities:", authentication.Authorities)

	// Output:
	// user id: d00ba9a3-d8bc-431c-ab36-37cd72ae825c
	// user name: d.zhupanov@redmadrobot.com
	// authorities: [user]
}
