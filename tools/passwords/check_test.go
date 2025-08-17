package passwords_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/passwords"
)

func TestCheckMatches(t *testing.T) {
	tests := []struct {
		name     string
		password string
		check    passwords.Check
		wantErr  bool
	}{
		{
			name:     "check contains passes",
			password: "a123",
			check:    passwords.CheckContains(`a|\d`, "doesn't contain pattern"),
		},
		{
			name:     "check contains fails",
			password: "bcd",
			check:    passwords.CheckContains(`a|\d`, "doesn't contain pattern"),
			wantErr:  true,
		},
		{
			name:     "check contains empty fails",
			password: "",
			check:    passwords.CheckContains(`a|\d`, "doesn't contain pattern"),
			wantErr:  true,
		},
		{
			name:     "check not contains passes",
			password: "bcd",
			check:    passwords.CheckNotContains(`a|\d`, "contains pattern"),
		},
		{
			name:     "check not contains fails",
			password: "a123",
			check:    passwords.CheckNotContains(`a|\d`, "contains pattern"),
			wantErr:  true,
		},
		{
			name:     "check not contains empty passes",
			password: "",
			check:    passwords.CheckNotContains(`a|\d`, "contains pattern"),
		},
		{
			name:     "check length passes",
			password: "a12345",
			check:    passwords.CheckLength(3, 8),
		},
		{
			name:     "check length fails (min)",
			password: "a",
			check:    passwords.CheckLength(3, 8),
			wantErr:  true,
		},
		{
			name:     "check length fails (max)",
			password: "01234567890",
			check:    passwords.CheckLength(3, 8),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.check(tt.password)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, passwords.ErrValidationFailed)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
