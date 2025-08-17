package mw

import (
	"context"
	"net/http"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/core/logs"
	"github.com/wal1251/pkg/core/security"
	"github.com/wal1251/pkg/httpx"
)

const UserInfoKey = "userInfo"

func Authorizer[T security.RequestCredentials](
	authManager security.Manager,
	credentialsProvider security.HTTPCredentialsProvider[T],
	authProvider security.AuthenticationProvider[T],
) httpx.Middleware {
	errResponse := httpx.ServerErrorResponses(
		httpx.MakeServerError,
		httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping()),
	)

	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		ctx := request.Context()
		logger := logs.FromContext(ctx)

		// Проверяем, требуется ли авторизация
		requirements := authManager.SecurityRequirements(ctx)

		logger.Info().Msg("authorizing user request")

		// Аутентифицируем пользователя с помощью дефолтного или своего провайдера
		authentication, err := authProvider.Authenticate(ctx, credentialsProvider(request))
		if err != nil {
			if !requirements.IsAuthorizationRequired {
				next.ServeHTTP(response, request)

				return
			}
			logger.Warn().Msg("invalid access token credentials")
			httpx.SendResponse(ctx, response, errResponse(err))

			return
		}

		// Добавляем данные аутентификации в логи
		logger.UpdateContext(logs.Options(
			logs.UserIDTag.Option(authentication.User.ID),
			logs.UserAuthorityTag.Option(logs.TagStringArray(authentication.Authorities.ToStrings())),
			logs.TokenIDTag.Option(authentication.TokenID),
		))

		// Проверяем права доступа
		if authentication.Authorities.Meets(requirements.Authorities) {
			logger.Info().Msg("user authority granted")

			// Добавляем информацию о пользователе в контекст
			userInfo := map[string]interface{}{
				"user_id":      authentication.User.ID,
				"phone_number": authentication.User.PhoneNumber,
				"name":         authentication.User.Name,
				"session_id":   authentication.SessionID,
			}
			// Используем ключ `UserInfoKey` для хранения данных в контексте
			ctx = context.WithValue(ctx, UserInfoKey, userInfo) //nolint:revive,staticcheck //FIXME

			// Обновляем контекст логгера и передаем управление дальше
			ctx = authentication.ToContext(ctx)
			ctx = logs.ToContext(ctx, logger)

			next.ServeHTTP(response, request.WithContext(ctx))

			return
		}

		logger.Warn().Msg("requested operation is forbidden")

		httpx.SendResponse(ctx, response, errResponse(errs.Wrapf(errs.ErrForbidden, "requested action is not permitted for current user")))
	}).Middleware()
}

func SecureAuthorizer[T security.RequestCredentials](
	authManager security.Manager,
	credentialsProvider security.HTTPCredentialsProvider[T],
	authProvider security.AuthenticationProvider[T],
) httpx.Middleware {
	errResponse := httpx.ServerErrorResponses(
		httpx.MakeSecureServerError,
		httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping()),
	)

	return httpx.MiddlewareFn(func(response http.ResponseWriter, request *http.Request, next http.Handler) {
		ctx := request.Context()
		logger := logs.FromContext(ctx)

		// Проверяем, требуется ли авторизация
		requirements := authManager.SecurityRequirements(ctx)

		logger.Info().Msg("authorizing user request")

		// Аутентифицируем пользователя с помощью дефолтного или своего провайдера
		authentication, err := authProvider.Authenticate(ctx, credentialsProvider(request))
		if err != nil {
			if !requirements.IsAuthorizationRequired {
				next.ServeHTTP(response, request)

				return
			}
			logger.Warn().Msg("invalid access token credentials")
			httpx.SendResponse(ctx, response, errResponse(errs.Reasons(err.Error(), errs.TypeAuthFailure, "1.1")))

			return
		}

		// Добавляем данные аутентификации в логи
		logger.UpdateContext(logs.Options(
			logs.UserIDTag.Option(authentication.User.ID),
			logs.UserAuthorityTag.Option(logs.TagStringArray(authentication.Authorities.ToStrings())),
			logs.TokenIDTag.Option(authentication.TokenID),
		))

		// Проверяем права доступа
		if authentication.Authorities.Meets(requirements.Authorities) {
			logger.Info().Msg("user authority granted")

			// Добавляем информацию о пользователе в контекст
			userInfo := map[string]interface{}{
				"user_id":      authentication.User.ID,
				"phone_number": authentication.User.PhoneNumber,
				"name":         authentication.User.Name,
				"session_id":   authentication.SessionID,
			}
			// Используем ключ `UserInfoKey` для хранения данных в контексте
			ctx = context.WithValue(ctx, UserInfoKey, userInfo) //nolint:revive,staticcheck //FIXME

			// Обновляем контекст логгера и передаем управление дальше
			ctx = authentication.ToContext(ctx)
			ctx = logs.ToContext(ctx, logger)

			next.ServeHTTP(response, request.WithContext(ctx))

			return
		}

		logger.Warn().Msg("requested operation is forbidden")

		httpx.SendResponse(ctx, response, errResponse(errs.Reasons("requested action is not permitted for current user", errs.TypeForbidden, "1.3")))
	}).Middleware()
}
