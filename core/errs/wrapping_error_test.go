package errs_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/errs"
)

func TestWith(t *testing.T) {
	tests := []struct {
		name      string
		reason    errs.Error
		err       error
		want      bool
		wantError string
	}{
		{
			name:      "Make error that will be the same by reason",
			reason:    errs.Error{Code: "foo"},
			err:       errs.With(errs.Error{Code: "foo"}, errors.New("fake")),
			want:      true,
			wantError: "foo: fake",
		},
		{
			name:      "Make error that will be not the same by reason",
			reason:    errs.Error{Code: "foo"},
			err:       errs.With(errs.Error{Code: "bar"}, errors.New("fake")),
			want:      false,
			wantError: "bar: fake",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if assert.Equalf(t, tt.want, errors.Is(tt.err, tt.reason), "equality check is not passed") {
				assert.Equalf(t, tt.wantError, tt.err.Error(), "error presentation is not valid")
			}
		})
	}
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		name      string
		reason    errs.Error
		err       error
		want      bool
		wantError string
	}{
		{
			name:      "Make error that will be the same by reason",
			reason:    errs.Error{Code: "foo"},
			err:       errs.Wrapf(errs.Error{Code: "foo"}, "fake: %d", 123),
			want:      true,
			wantError: "foo: fake: 123",
		},
		{
			name:      "Make error that will be not the same by reason",
			reason:    errs.Error{Code: "foo"},
			err:       errs.Wrapf(errs.Error{Code: "bar"}, "fake: %d", 123),
			want:      false,
			wantError: "bar: fake: 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if assert.Equalf(t, tt.want, errors.Is(tt.err, tt.reason), "equality check is not passed") {
				assert.Equalf(t, tt.wantError, tt.err.Error(), "error presentation is not valid")
			}
		})
	}
}
