package checks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/checks"
)

func TestIsEmptyString(t *testing.T) {
	assert.True(t, checks.IsEmpty(""), "Empty string must return true")
	assert.False(t, checks.IsEmpty("Go"), "Non-empty string must return false")
}

func TestIsEmptyInt(t *testing.T) {
	assert.True(t, checks.IsEmpty(0), "Default int value must return true")
	assert.False(t, checks.IsEmpty(42), "Non-default int value must return false")
}

func TestIsEmptyBool(t *testing.T) {
	assert.True(t, checks.IsEmpty(false), "false (as a default value of bool) must return true")
	assert.False(t, checks.IsEmpty(true), "true must return false")
}

func TestIsEmptyPointer(t *testing.T) {
	var strPtr *string
	assert.True(t, checks.IsEmpty(strPtr), "Nil pointer must return true")

	str := "Go"
	strPtr = &str
	assert.False(t, checks.IsEmpty(strPtr), "Non-nil pointer must return false")
}

func TestIsEmptyAny(t *testing.T) {
	tests := []struct {
		name string
		arg  any
		want bool
	}{
		{
			name: "Empty string",
			arg:  "",
			want: true,
		},
		{
			name: "Non-empty string",
			arg:  "Hello",
			want: false,
		},
		{
			name: "Empty slice",
			arg:  []int{},
			want: true,
		},
		{
			name: "Non-empty slice",
			arg:  []int{1, 2, 3},
			want: false,
		},
		{
			name: "Empty map",
			arg:  map[string]int{},
			want: true,
		},
		{
			name: "Non-empty map",
			arg:  map[string]int{"a": 1},
			want: false,
		},
		{
			name: "Default bool (false)",
			arg:  false,
			want: true,
		},
		{
			name: "true bool",
			arg:  true,
			want: false,
		},
		{
			name: "Default int (0)",
			arg:  0,
			want: true,
		},
		{
			name: "Non-default int",
			arg:  42,
			want: false,
		},
		{
			name: "nil pointer",
			arg:  (*int)(nil),
			want: true,
		},
		{
			name: "Non-nil pointer",
			arg:  new(int),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checks.IsEmptyAny(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMustBePositive(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		paramName string
		wantPanic bool
	}{
		{
			name:      "Positive value",
			value:     1,
			paramName: "age",
			wantPanic: false,
		},
		{
			name:      "Zero value",
			value:     0,
			paramName: "age",
			wantPanic: true,
		},
		{
			name:      "Negative value",
			value:     -1,
			paramName: "age",
			wantPanic: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tc.wantPanic {
					t.Errorf("MustBePositive() panic = %v, wantPanic %v", r, tc.wantPanic)
				}
			}()
			checks.MustBePositive(tc.value, tc.paramName)
		})
	}
}
