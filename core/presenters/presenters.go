// Package presenters предоставляет функции формирования представлений объектов.
// Например, возможность скрытия чувствительной информации в логах.
package presenters

import (
	"fmt"
	"io"
)

const (
	ViewIdentity ViewType = ""             // Представление "как есть".
	ViewPrivate  ViewType = "VIEW_PRIVATE" // Представление для внутреннего использования.
	ViewLogs     ViewType = "VIEW_LOGS"    // Представление для логирования.
	ViewPublic   ViewType = "VIEW_PUBLIC"  // Представление "для всех".

	DefaultCredentialsPlaceholder = "{hidden}"
)

type (
	// ViewType тип представления, не все можно показывать всем. Влияет на то в каком виде объект должен быть
	// представлен пользователю.
	ViewType string

	// ViewOptions дополнительные параметры формирования представления объекта.
	ViewOptions struct {
		SecuredKeywords []string // Атрибуты с этими полями необходимо скрывать от пользователей.
		MaxStringLength int      // Максимальная длина представления объекта. Накладно выводить большие объекты целиком.
	}

	// StringViewer объект реализующий интерфейс осуществляет поддержку механизма формирования строчных представлений
	// объекта.
	// Может быть полезно, например, если хочется скрыть чувствительную информацию в логах.
	StringViewer interface {
		StringView(ViewType, ViewOptions) string
	}

	// StringPresenterWriter предназначен для записи представления с последующим получением результата записи в виде
	// строки. Например:
	//
	//	var withoutPassword StringPresenterWriter = credentialsHidingBuf()
	//	err := json.NewEncoder(withoutPassword).Encode(struct{ password string }{password: "secret"})
	//	// ...
	//	fmt.Println(withoutPassword)
	//
	// См. Cropper.
	StringPresenterWriter interface {
		io.Writer
		fmt.Stringer
	}
)

func NewViewOptions(cfg *Config) ViewOptions {
	return ViewOptions{
		SecuredKeywords: cfg.SecuredKeywords,
		MaxStringLength: cfg.MaxStringLength,
	}
}

func DefaultViewOptions() ViewOptions {
	return ViewOptions{}
}
