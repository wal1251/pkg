package indices_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/providers/es/indices"
)

func TestIsValidCorrectType(t *testing.T) {
	tests := []struct {
		name string
		in   indices.Dynamic
		out  bool
	}{
		{
			name: "true",
			in:   indices.DynamicTrue,
			out:  true,
		},
		{
			name: "runtime",
			in:   indices.DynamicRuntime,
			out:  true,
		},
		{
			name: "false",
			in:   indices.DynamicFalse,
			out:  true,
		},
		{
			name: "strict",
			in:   indices.DynamicStrict,
			out:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.IsValid())
		})
	}
}

func TestIsValidInCorrectType(t *testing.T) {
	tests := []struct {
		name string
		in   indices.Dynamic
		out  bool
	}{
		{
			name: "not_true",
			in:   indices.Dynamic("not_true"),
			out:  false,
		},
		{
			name: "not_runtime",
			in:   indices.Dynamic("not_runtime"),
			out:  false,
		},
		{
			name: "not_false",
			in:   indices.Dynamic("not_false"),
			out:  false,
		},
		{
			name: "not_strict",
			in:   indices.Dynamic("not_strict"),
			out:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.IsValid())
		})
	}
}
