package viperx_test

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/core/cfg/viperx"
)

func TestViper(t *testing.T) {
	v := viper.New()
	v.Set("sample_string", "foo")
	v.Set("sample_bool", "true")
	v.Set("sample_int", "777")
	v.Set("sample_float", "3.14")
	v.Set("sample_duration", "3s")

	assert.Equal(t, "foo", viperx.Get(v, "sample_string", ""))
	assert.Equal(t, true, viperx.Get(v, "sample_bool", false))
	assert.Equal(t, 777, viperx.Get(v, "sample_int", 0))
	assert.Equal(t, 3.14, viperx.Get(v, "sample_float", 0.0))
	assert.Equal(t, 3*time.Second, viperx.Get(v, "sample_duration", time.Duration(0)))
	assert.Equal(t, "default", viperx.Get(v, "sample_not_set", "default"))
}
