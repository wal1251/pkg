package security_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/security"
)

func TestDefaultAuthorizationManager_Authentication(t *testing.T) {
	tests := []struct {
		name string
		auth security.Authentication
	}{
		{
			name: "Basic test",
			auth: security.Authentication{
				User:        security.User{ID: uuid.New(), Name: "foo"},
				Authorities: security.Authorities{"bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.auth.ToContext(context.TODO())
			assert.Equal(t, tt.auth, security.DefaultManager{}.Authorized(ctx))
		})
	}
}

func TestDefaultAuthorizationManager_SecurityRequirements(t *testing.T) {
	tests := []struct {
		name         string
		requirements []string
		want         security.Requirements
	}{
		{
			name:         "Basic test",
			requirements: []string{"foo", "bar"},
			want: security.Requirements{
				IsAuthorizationRequired: true,
				Authorities:             security.Authorities{"foo", "bar"},
			},
		},
		{
			name:         "Permitted test",
			requirements: []string{},
			want: security.Requirements{
				IsAuthorizationRequired: false,
				Authorities:             security.Authorities{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := struct{}{}
			ctx := context.WithValue(context.TODO(), key, tt.requirements)
			assert.Equal(t, tt.want, security.DefaultManager{AuthoritiesContextKey: key}.SecurityRequirements(ctx))
		})
	}
}
