package presenters_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/rs/zerolog"

	"github.com/wal1251/pkg/core/presenters"
	"github.com/wal1251/pkg/httpx/mw"
	"github.com/wal1251/pkg/proxy/hooks"
)

func ExampleCropper() {
	cropped := presenters.NewCropper(10)
	src := bytes.NewBufferString("hello hello - is very lo-o-o-ong string!")
	if _, err := io.Copy(cropped, src); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cropped)

	// Output:
	// hello hell...
}

// ExampleNewViewOptions_httpMiddleware пример скрытия чувствительных данных в логах,
// которые создаются автоматически при использовании mw.Logger().
func ExampleNewViewOptions_httpMiddleware() {
	cfg := &presenters.Config{SecuredKeywords: []string{"password"}}
	options := presenters.NewViewOptions(cfg)
	// Создаем middleware для логирования запросов, который скрывает чувствительные данные (поле 'password').
	loggerMiddleware := mw.Logger(presenters.ViewLogs, options)

	var loggerBuffer bytes.Buffer
	logger := zerolog.New(&loggerBuffer)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Добавляем логгер в контекст запроса.
		// Внимание! Это необходимо лишь в этом примере, т.к. httptest не содержит возможности указать базовый контекст.
		// В реальном приложении это происходит автоматически при вызове функции httpx.StartServer().
		loggerMiddleware(handler).ServeHTTP(w, r.WithContext(logger.WithContext(r.Context())))
	}))
	defer server.Close()

	_, err := http.Post(server.URL, "application/json", bytes.NewBufferString(`{"name": "John", "password": "123456"}`))
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем, что пароль был скрыт в логах.
	if strings.Contains(loggerBuffer.String(), `"request body: {\"name\": \"John\", \"password\": \"{hidden}\"}"`) {
		fmt.Println("Password was hidden in the logs.")
	}

	// Output:
	// Password was hidden in the logs.
}

// ExampleNewViewOptions_loggerHook пример скрытия чувствительных данных в логах, которые создаются автоматически
// при использовании хуков логирования.
func ExampleNewViewOptions_loggerHook() {
	cfg := &presenters.Config{SecuredKeywords: []string{"password"}}
	options := presenters.NewViewOptions(cfg)

	var loggerBuffer bytes.Buffer
	logger := zerolog.New(&loggerBuffer)

	// Создаем хук логирования, который скрывает чувствительные данные (поле 'password').
	hook := hooks.LogBeforeCall(presenters.ViewLogs, options)
	hook(logger.WithContext(context.Background()), nil, "method", []interface{}{
		struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "John",
			Password: "123456",
		},
	})

	// Проверяем, что пароль был скрыт в логах.
	if strings.Contains(loggerBuffer.String(), `"args: [{\"name\":\"John\",\"password\":\"{hidden}\"}]"`) {
		fmt.Println("Password was hidden in the logs.")
	}

	// Output:
	// Password was hidden in the logs.
}
