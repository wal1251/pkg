package presenters

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/wal1251/pkg/tools/serial"
)

// ParameterView возвращает строку представление параметра (функции, запроса...).
func ParameterView(param any, view ViewType, options ViewOptions) string {
	var builder strings.Builder
	if param == nil {
		builder.WriteString("nil")

		return builder.String()
	}

	switch value := param.(type) {
	case StringViewer:
		builder.WriteRune('"')
		builder.WriteString(value.StringView(view, options))
		builder.WriteRune('"')

		return builder.String()
	case context.Context:
		builder.WriteString("context.Context")

		return builder.String()
	case error:
		if value != nil {
			builder.WriteString("{\"error\": \"")
			builder.WriteString(value.Error())
			builder.WriteString("\"}")

			return builder.String()
		}
	}

	var buffer StringPresenterWriter
	if options.MaxStringLength == 0 {
		buffer = bytes.NewBufferString("")
	} else {
		buffer = NewCropper(options.MaxStringLength)
	}

	if err := serial.JSONEncode(buffer, param); err != nil {
		typ := reflect.TypeOf(param)
		if typ == nil {
			builder.WriteString("***FAILED TO REPRESENT***")
		} else {
			builder.WriteString(typ.String())
		}

		return builder.String()
	}

	builder.WriteString(JSONString(strings.TrimSuffix(buffer.String(), "\n"), view, options))

	return builder.String()
}

// ParameterListView возвращает строку представление списка параметра (функции, запроса...).
func ParameterListView(params []any, view ViewType, options ViewOptions) string {
	var builder strings.Builder

	builder.WriteRune('[')

	for i, param := range params {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(ParameterView(param, view, options))
	}

	builder.WriteRune(']')

	return builder.String()
}

func JSONString(value string, view ViewType, options ViewOptions) string {
	if view == ViewLogs || view == ViewPublic {
		return JSONHideCredentials(value, options)
	}

	return value
}

func JSONHideCredentials(value string, options ViewOptions) string {
	for _, keyword := range options.SecuredKeywords {
		re := regexp.MustCompile(fmt.Sprintf(`(?mi)("%s"\s*:\s*").*?(")`, keyword))
		value = re.ReplaceAllString(value, fmt.Sprintf(`$1%s$2`, DefaultCredentialsPlaceholder))
	}

	return value
}
