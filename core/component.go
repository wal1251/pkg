package core

var _ Component = (*DefaultComponent)(nil)

type (
	// Component интерфейс компонента приложения.
	Component interface {
		// Name возвращает имя компонента, компоненты одного класса имеют одинаковое имя.
		Name() string
		// Label возвращает отличительную метку экземпляра компонента, например разные экземпляры одного класса
		// компоненты могут иметь одинаковое имя, но разные метки.
		Label() string
	}

	// DefaultComponent простейшая реализация интерфейса Component предоставляет базовую функциональность идентификации
	// компонент приложения.
	DefaultComponent struct {
		name  string
		label string
	}
)

func (c *DefaultComponent) Name() string {
	return c.name
}

func (c *DefaultComponent) Label() string {
	return c.label
}

// NewDefaultComponent возвращает реализацию компоненты по-умолчанию.
func NewDefaultComponent(name, label string) *DefaultComponent {
	return &DefaultComponent{
		name:  name,
		label: label,
	}
}
