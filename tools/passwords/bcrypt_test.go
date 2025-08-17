package passwords_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/passwords"
)

func TestBcryptEncryptor(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Регулярный кейс",
			password: "fooBarBaz",
		},
		{
			name:     "Пустая строка",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor := passwords.NewBcryptEncryptor()
			encrypted, err := encryptor.Encrypt(tt.password)
			if assert.NoError(t, err) {
				assert.NotEqual(t, tt.password, encrypted)
				assert.True(t, encryptor.Verify(tt.password, encrypted))
			}
		})
	}
}
