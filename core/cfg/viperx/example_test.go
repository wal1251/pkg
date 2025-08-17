package viperx_test

import (
	"fmt"

	"github.com/wal1251/pkg/core/cfg"
	"github.com/wal1251/pkg/core/cfg/viperx"
)

func ExampleEnvLoader() {
	AppPrefix := "SAMPLE_APP"

	// Разделим неймспейсы конфигов различных сред.
	namespaceProd := "PROD"
	namespaceTest := "TEST"

	// Установим конфиг приложения в env-вах.
	//
	// SAMPLE_APP_ENVIRONMENT=prod
	_ = cfg.EnvVariableSet(cfg.KeyEnvironment, cfg.EnvProd,
		cfg.KeyWithPrefix(AppPrefix),
	)

	// SAMPLE_APP_PROD_FOO=777
	_ = cfg.EnvVariableSet("FOO", "777",
		cfg.KeyWithPrefix(namespaceProd), // переменная для продуктивной среды.
		cfg.KeyWithPrefix(AppPrefix),
	)

	// SAMPLE_APP_TEST_FOO=555
	_ = cfg.EnvVariableSet("FOO", "555",
		cfg.KeyWithPrefix(namespaceTest), // переменная для тестовой среды.
		cfg.KeyWithPrefix(AppPrefix),
	)

	// Создадим загрузчик конфига.
	loader := viperx.EnvLoader(AppPrefix)

	// Мы в продуктивной среде?
	var namespace string
	if viperx.Environment(loader).Is(cfg.EnvProd) {
		namespace = namespaceProd
	} else {
		namespace = namespaceTest
	}

	// Читаем конфиг.
	fmt.Println("NAMESPACE:", namespace)
	fmt.Println("FOO:", viperx.Get(loader, cfg.Key("FOO").Map(cfg.KeyWithPrefix(namespace)), 456))
	fmt.Println("BAR:", viperx.Get(loader, "BAR", 789)) // переменная не установлена.

	// Output:
	// NAMESPACE: PROD
	// FOO: 777
	// BAR: 789
}
