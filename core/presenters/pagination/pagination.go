// Package pagination предоставляет интерфейсы и функции для управления постраничным выводом представления коллекций.
package pagination

import (
	"errors"
	"fmt"
)

const DefaultLimit = 15 // Максимальное значение элементов в выдаче по умолчанию.

var (
	ErrInvalidOffset = errors.New("invalid offset") // Задан некорректный параметр смещения.
	ErrInvalidLimit  = errors.New("invalid limit")  // Задан некорректный параметр размера выдачи.
)

var _ PageQuery = (*PageParams)(nil)

type (
	// PageQuery запрос параметров постраничной выдачи представления коллекций.
	PageQuery interface {
		GetOffset() int  // Смещение первого элемента, с которого необходимо отобразить коллекцию.
		GetLimit() int   // Максимальный размер выдачи результата.
		Validate() error // Выполнит проверку заданного параметров.
	}

	// PageParams хранит параметры постраничной выдачи представления коллекций.
	PageParams struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}

	// Page представляет страницу, которая содержит результат выборки элементов коллекции. Поддерживает интерфейс PageQuery.
	Page[T any] struct {
		PageParams
		Content    []*T `json:"content"`
		TotalCount int  `json:"totalCount"`
	}
)

// GetOffset см. PageQuery.GetOffset().
func (p *PageParams) GetOffset() int {
	return p.Offset
}

// GetLimit см. PageQuery.GetLimit().
func (p *PageParams) GetLimit() int {
	return p.Limit
}

// Validate см. PageQuery.Validate().
func (p *PageParams) Validate() error {
	if p.Offset < 0 {
		return fmt.Errorf("%w: is negative", ErrInvalidOffset)
	}

	if p.Limit < 0 {
		return fmt.Errorf("%w: is negative", ErrInvalidLimit)
	}

	return nil
}

// WithOffset устанавливает смещение первого элемента для выборки.
func (p *PageParams) WithOffset(offset *int) *PageParams {
	p.Offset = 0
	if offset != nil {
		p.Offset = *offset
	}

	return p
}

// WithLimit устанавливает максимальное количество элементов, возвращаемых в выборке.
func (p *PageParams) WithLimit(limit *int) *PageParams {
	p.Limit = DefaultLimit
	if limit != nil {
		p.Limit = *limit
	}

	return p
}

// Clone вернет копию объекта.
func (p *PageParams) Clone() PageParams {
	return *p
}

// NewPageParams создает новый экземпляр параметров постраничной выдачи.
func NewPageParams(offset, limit *int) *PageParams {
	return (&PageParams{}).WithLimit(limit).WithOffset(offset)
}
