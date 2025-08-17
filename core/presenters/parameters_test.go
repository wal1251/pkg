package presenters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/presenters"
)

func TestJSONHideCredentials(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Базовый тест",
			s: `{"email": "max.ivanov@example.com", "PasSword": "Ku1RjM5iexD",
						"Token": "1234lldnk3333..mnn3434kkknn%%6773$kjd",
						"PasswordCurrent": ""}`,
			want: `{ "email": "max.ivanov@example.com", "PasSword": "{hidden}", "Token": "{hidden}", "PasswordCurrent":"{hidden}"}`,
		},
		{
			name: "Пустая строка",
			s: `{"email": "max.ivanov@example.com", "PasSword": "",
						"Token": "", "PasswordCurrent":""}`,
			want: `{ "email": "max.ivanov@example.com", "PasSword": "{hidden}", "Token": "{hidden}", "PasswordCurrent":"{hidden}"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := presenters.JSONHideCredentials(tt.s, presenters.ViewOptions{
				SecuredKeywords: []string{"password", "token", "passwordCurrent"},
			})
			assert.JSONEq(t, tt.want, result)
		})
	}
}
