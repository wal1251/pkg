// Package viperx является расширением библиотеки viper: функции для легкого подключения конфигураций к приложению.
package viperx

import (
	"time"

	"github.com/spf13/viper"

	"github.com/wal1251/pkg/core/cfg"
)

var _ cfg.Property[string] = (*Property[string])(nil)

// Property предоставляет функции загрузки свойства с помощью библиотеки viper.
type Property[T cfg.ValueType] struct {
	key   string
	viper *viper.Viper
}

// Get см. Property.
func (v Property[T]) Get() T {
	var value T

	switch {
	case cfg.PropertyProvider[string](v.viper.GetString).TypeMatches(value):
		value = cfg.PropertyProviderAdapter[string, T](v.viper.GetString).Get(v.key)
	case cfg.PropertyProvider[bool](v.viper.GetBool).TypeMatches(value):
		value = cfg.PropertyProviderAdapter[bool, T](v.viper.GetBool).Get(v.key)
	case cfg.PropertyProvider[int](v.viper.GetInt).TypeMatches(value):
		value = cfg.PropertyProviderAdapter[int, T](v.viper.GetInt).Get(v.key)
	case cfg.PropertyProvider[float64](v.viper.GetFloat64).TypeMatches(value):
		value = cfg.PropertyProviderAdapter[float64, T](v.viper.GetFloat64).Get(v.key)
	case cfg.PropertyProvider[time.Duration](v.viper.GetDuration).TypeMatches(value):
		value = cfg.PropertyProviderAdapter[time.Duration, T](v.viper.GetDuration).Get(v.key)
	default:
	}

	return value
}

// Get загружает значение заданного свойства конфигурации, с помощью библиотеки viper.
func Get[T cfg.ValueType](v *viper.Viper, key cfg.Key, defaultValue T) T {
	return NewProperty(v, key, defaultValue).Get()
}

// NewProperty возвращает заданное свойство конфигурации, с помощью библиотеки viper.
func NewProperty[T cfg.ValueType](loader *viper.Viper, key cfg.Key, defaultValue T) *Property[T] {
	var z T
	if defaultValue != z {
		loader.SetDefault(string(key), defaultValue)
	}

	return &Property[T]{
		key:   string(key),
		viper: loader,
	}
}

// EnvLoader создает и возвращает экземпляр viper загрузчика конфигурации из переменных среды ОС.
func EnvLoader(prefix string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()

	return v
}

// Environment возвращает значение среды выполнения приложения с помощью viper.
func Environment(v *viper.Viper) cfg.Environment {
	return cfg.Environment(Get(v, cfg.KeyEnvironment, cfg.EnvDev.String()))
}
