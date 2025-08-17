package errs_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/errs"
)

func TestError_Equals(t *testing.T) {
	tests := []struct {
		name   string
		reason errs.Error
		Code   string
		Type   errs.Type
		want   bool
	}{
		{
			name:   "Equals with same code and type",
			reason: errs.Error{Code: "foo", Type: errs.Type("bar")},
			Code:   "foo",
			Type:   "bar",
			want:   true,
		},
		{
			name:   "Equals with same code but different types",
			reason: errs.Error{Code: "foo", Type: errs.Type("bar")},
			Code:   "foo",
			Type:   "baz",
			want:   true,
		},
		{
			name:   "Not equals with different codes",
			reason: errs.Error{Code: "foo", Type: errs.Type("bar")},
			Code:   "baz",
			Type:   "bar",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.reason.Equals(errs.Error{
				Code: tt.Code,
				Type: tt.Type,
			}))
		})
	}
}

func TestError_Is(t *testing.T) {
	tests := []struct {
		name   string
		reason errs.Error
		err    error
		want   bool
	}{
		{
			name:   "Reason is same reason",
			reason: errs.Error{Code: "foo"},
			err:    errs.Error{Code: "foo"},
			want:   true,
		},
		{
			name:   "Reason is same reason (pointer)",
			reason: errs.Error{Code: "foo"},
			err:    &errs.Error{Code: "foo"},
			want:   true,
		},
		{
			name:   "Reason is not same on different reason",
			reason: errs.Error{Code: "foo"},
			err:    errs.Error{Code: "bar"},
			want:   false,
		},
		{
			name:   "Reason is not same on different error",
			reason: errs.Error{Code: "foo"},
			err:    errors.New("foo"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.reason.Is(tt.err))
		})
	}
}
