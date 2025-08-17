package api

import "io"

type (
	Factory[T any]       func() (T, error)
	FactoryOption[T any] func(*T) error
)

func (f Factory[T]) JSONReader() io.Reader {
	if f == nil {
		return nil
	}

	value, err := f.Produce()
	if err != nil {
		return NewErrReader(err)
	}

	return NewJSONReader(value)
}

func (f Factory[T]) Produce() (T, error) {
	if f == nil {
		var blank T

		return blank, nil
	}

	value, err := (func() (T, error))(f)()
	if err != nil {
		var blank T

		return blank, nil //nolint:nilerr
	}

	return value, nil
}

func (f Factory[T]) With(opts ...FactoryOption[T]) Factory[T] {
	return func() (T, error) {
		value, err := f.Produce()
		if err != nil {
			var blank T

			return blank, err
		}

		for _, opt := range opts {
			if err = opt(&value); err != nil {
				var blank T

				return blank, err
			}
		}

		return value, nil
	}
}
