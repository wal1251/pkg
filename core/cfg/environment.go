package cfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MarkerTest = "TEST" // Маркер тестового объекта или среды (используется как префикс или постфикс в наименованиях свойств).

	KeyEnvironment Key = "ENVIRONMENT" // Ключ: среда выполнения приложения.
	KeyPrefix      Key = "PREFIX"      // Ключ: префикс конфигурации приложения. Определяет пространство имен конфига.

	EnvDev   Environment = "dev"   // Среда разработки.
	EnvStage Environment = "stage" // Тестовая среда.
	EnvProd  Environment = "prod"  // Продуктивная среда.
)

// Environment тип среды выполнения приложения.
type Environment string

func (e Environment) String() string {
	return string(e)
}

// Is вернет true, если среда эквивалентна одной из указанных. Если ни одной целевой среды не указано, вернет false.
func (e Environment) Is(envs ...Environment) bool {
	for _, env := range envs {
		if e == env {
			return true
		}
	}

	return false
}

// EnvVariableSet устанавливает переменную окружения ОС.
func EnvVariableSet(key Key, value any, mapping ...KeyMap) error {
	if value == nil {
		return nil
	}

	valueString := strings.TrimSpace(fmt.Sprint(value))
	if valueString == "" {
		return nil
	}

	keyString := key.Map(mapping...).String()

	if err := os.Setenv(keyString, valueString); err != nil {
		return fmt.Errorf("can't set OS env %s=%s: %w", keyString, valueString, err)
	}

	return nil
}

// IsTestRuntime получает информацию о том является ли runtime запуском тестов.
// Эта информация хранится в переменной среды %s_TEST, где %s - значение переменной PREFIX, устанавливаемой viper.
// Переменная должна быть установлена до запуска, потому что импорты срабатывают раньше, чем выполняются какие-либо
// инструкции, которые могут установить переменные среды, как это работает в случае с viper и как это победить - я не
// понял.
func IsTestRuntime() bool {
	var prefix string

	if value, ok := os.LookupEnv(KeyPrefix.String()); ok {
		prefix = value
	}

	if value, ok := os.LookupEnv(Key(MarkerTest).Map(KeyWithPrefix(prefix)).String()); ok {
		result, _ := strconv.ParseBool(value)

		return result
	}

	return false
}
