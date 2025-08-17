package cfg_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/cfg"
)

func TestPropertyProvider_TypeMatches(t *testing.T) {
	pString := cfg.PropertyProvider[string](func(string) string { return "" })
	assert.True(t, pString.TypeMatches("foo"))
	assert.False(t, pString.TypeMatches(time.Second))
	assert.False(t, pString.TypeMatches(true))
	assert.False(t, pString.TypeMatches(1))

	pDuration := cfg.PropertyProvider[time.Duration](func(string) time.Duration { return time.Duration(0) })
	assert.False(t, pDuration.TypeMatches("foo"))
	assert.True(t, pDuration.TypeMatches(time.Second))
	assert.False(t, pDuration.TypeMatches(true))
	assert.False(t, pDuration.TypeMatches(1))

	pBool := cfg.PropertyProvider[bool](func(string) bool { return false })
	assert.False(t, pBool.TypeMatches("foo"))
	assert.False(t, pBool.TypeMatches(time.Second))
	assert.True(t, pBool.TypeMatches(true))
	assert.False(t, pBool.TypeMatches(1))

	pInt := cfg.PropertyProvider[int](func(string) int { return 0 })
	assert.False(t, pInt.TypeMatches("foo"))
	assert.False(t, pInt.TypeMatches(time.Second))
	assert.False(t, pInt.TypeMatches(true))
	assert.True(t, pInt.TypeMatches(1))
}

func TestPropertyProviderAdapter(t *testing.T) {
	p1 := cfg.PropertyProvider[string](func(string) string { return "foo" })
	p2 := cfg.PropertyProviderAdapter[string, string](p1)

	assert.Equal(t, "foo", p2.Get(""))
}
