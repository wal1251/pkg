package anyobj

import (
	"encoding/json"
	"fmt"
)

// SafeCopy осуществляет безопасное копирование данных из src в dest.
// Эта функция сначала преобразует src в JSON, а затем обратно из JSON в dest.
// Это гарантирует, что данные в dest являются копией src и не содержат ссылок на оригинальные данные src.
// Однако важно учитывать, что SafeCopy не копирует неэкспортируемые поля структур, поскольку стандартная библиотека
// encoding/json не обрабатывает неэкспортируемые поля.
//
// Параметры:
// src - исходный объект, который должен быть копирован.
// dest - целевой объект, куда будет производиться копирование. Должен быть указателем на тип.
//
// Возвращаемое значение:
// Возвращает ошибку, если произошла ошибка при сериализации src в JSON или при десериализации JSON в dest.
// Если копирование прошло успешно, возвращается nil.
func SafeCopy(src, dest any) error {
	raw, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("can't marshal source: %w", err)
	}

	if err = json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("can't marshal to destination: %w", err)
	}

	return nil
}
