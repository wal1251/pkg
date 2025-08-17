package size_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/size"
)

func TestCount(t *testing.T) {
	tests := []struct {
		size     size.Size
		unit     size.Size
		expected int64
	}{
		{1024, size.B, 1024},
		{1024, size.KB, 1},
		{1536, size.KB, 1},
		{1048576, size.MB, 1},
		{1073741824, size.GB, 1},
		{1099511627776, size.TB, 1},
		{0, size.KB, 0}, // Edge case: size is 0
		{1024, 0, 1024}, // Edge case: unit is 0
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%dB in %s = %d", tt.size, tt.unit, tt.expected)
		t.Run(testName, func(t *testing.T) {
			result := tt.size.Count(tt.unit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntAndInt64(t *testing.T) {
	tests := []struct {
		size          size.Size
		expectedInt   int
		expectedInt64 int64
	}{
		{1024, 1024, 1024},
		{1073741824, 1073741824, 1073741824},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%d to Int() and Int64()", tt.size)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tt.expectedInt, tt.size.Int())
			assert.Equal(t, tt.expectedInt64, tt.size.Int64())
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		size     size.Size
		expected string
	}{
		{1, "1 B"},
		{1024, "1 KB"},
		{1536, "1.5 KB"},
		{524288, "512 KB"},
		{1048576, "1 MB"},
		{1073741824, "1 GB"},
		{1099511627776, "1 TB"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%d to String() = %s", tt.size, tt.expected)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.size.String())
		})
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		size     int
		unit     size.Size
		expected int
	}{
		{1, size.KB, 1024},
		{1, size.MB, 1048576},
		{1, size.GB, 1073741824},
		{1, size.TB, 1099511627776},
		{0, size.GB, 0}, // Edge case: size is 0
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%d %s to Bytes() = %d", tt.size, tt.unit, tt.expected)
		t.Run(testName, func(t *testing.T) {
			result := size.Bytes(tt.size, tt.unit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMake(t *testing.T) {
	tests := []struct {
		size     int
		unit     size.Size
		expected size.Size
	}{
		{1, size.KB, 1024},
		{1, size.MB, 1048576},
		{1, size.GB, 1073741824},
		{1, size.TB, 1099511627776},
		{0, size.GB, 0}, // Edge case: size is 0
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("Make(%d, %s) = %d", tt.size, tt.unit, tt.expected)
		t.Run(testName, func(t *testing.T) {
			result := size.Make(tt.size, tt.unit)
			assert.Equal(t, tt.expected, result)
		})
	}
}
