package pagination

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wal1251/pkg/tools/collections"
)

const (
	SortAsc  SortDirection = "asc"  // Сортировка элементов по возрастанию.
	SortDesc SortDirection = "desc" // Сортировка элементов по убыванию.

	DefaultSortDirection = SortAsc // Порядок сортировки по умолчанию.
)

var (
	ErrInvalidSortDirection = errors.New("invalid sorting Direction")   // Некорректный параметр направления сортировки.
	ErrInvalidSortField     = errors.New("invalid sorting mappedField") // Некорректное поле сортировки.
)

var _ SortQuery = (*SortParams)(nil)

type (
	// SortDirection порядок сортировки выдачи.
	SortDirection string

	// SortQuery запрос порядка выдачи элементов в представлении коллекции.
	SortQuery interface {
		GetDirection() SortDirection // Вернет порядок сортировки элементов.
		GetFields() []string         // Вернет поля элемента коллекции, по которым необходимо выполнить упорядочивание.
		Validate() error             // Выполнит проверку заданного параметров.
	}

	// SortParams хранит параметры сортировки выдачи представления коллекции. Поддерживает интерфейс SortQuery.
	SortParams struct {
		Direction        SortDirection
		Fields           []string
		FieldsAcceptable collections.Set[string]
		FieldsMapping    map[string]string
	}
)

// Validate вернет ошибку, если значение порядка сортировки является некорректным значением.
func (d SortDirection) Validate() error {
	if d == SortAsc || d == SortDesc {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrInvalidSortDirection, d)
}

// GetDirection см. SortQuery.GetDirection().
func (s *SortParams) GetDirection() SortDirection {
	return s.Direction
}

// GetFields см. SortQuery.GetFields().
func (s *SortParams) GetFields() []string {
	return collections.Map(s.Fields, s.mappedField)
}

// Validate см. SortQuery.Validate().
func (s *SortParams) Validate() error {
	if err := s.Direction.Validate(); err != nil {
		return err
	}

	if !s.FieldsAcceptable.Contains(s.Fields...) {
		return ErrInvalidSortField
	}

	return nil
}

// WithMapping устанавливает маппинг полей.
func (s *SortParams) WithMapping(m map[string]string) *SortParams {
	s.FieldsMapping = m

	return s
}

// WithAcceptableFields добавляет допустимые поля.
func (s *SortParams) WithAcceptableFields(fields ...string) *SortParams {
	s.FieldsAcceptable.Add(fields...)

	return s
}

// WithDefaultFields добавить поля сортировки, если они ранее не были добавлены.
func (s *SortParams) WithDefaultFields(fields ...string) *SortParams {
	if len(s.Fields) == 0 {
		s.Fields = append(s.Fields, fields...)
	}

	return s
}

func (s *SortParams) mappedField(field string) string {
	if m, ok := s.FieldsMapping[field]; ok {
		return m
	}

	return field
}

// NewSorting создает новый объект с параметрами сортировки.
func NewSorting(fields []string, direction *string) *SortParams {
	sortDirection := DefaultSortDirection
	if direction != nil {
		sortDirection = SortDirection(strings.ToLower(*direction))
	}

	return &SortParams{
		Direction:        sortDirection,
		Fields:           fields,
		FieldsAcceptable: collections.NewSet[string](),
		FieldsMapping:    make(map[string]string),
	}
}
