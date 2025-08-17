package httpx_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/errs"
	"github.com/wal1251/pkg/httpx"
)

func TestHTTPStatus(t *testing.T) {
	errMap := httpx.NewErrorToStatusMapper(httpx.DefaultErrorToStatusMapping())

	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "Unrecognized reason causes 500",
			err:  errors.New("fake"),
			want: errMap.SystemFailure,
		},
		{
			name: "Wrapped reason",
			err:  fmt.Errorf("%w", errs.ErrNotFound),
			want: errMap.Status(errs.ErrNotFound.Type),
		},
		{
			name: "Unclassified reason",
			err:  errs.Error{Code: "FAKE"},
			want: errMap.Default,
		},
		{
			name: "Wrapped described reason",
			err:  fmt.Errorf("%w", errs.Wrapf(errs.ErrNotFound, "object not found")),
			want: errMap.Status(errs.ErrNotFound.Type),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, errMap.Status(errs.AsReason(tt.err).Type))
		})
	}
}
