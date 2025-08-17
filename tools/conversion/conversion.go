package conversion

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/wal1251/pkg/tools/collections"
)

var ErrStrToNumberConversion = errors.New("unable to convert string to number")

// TableRowStringToPointerString This function takes a column name and a table row object as input and returns
// a pointer to a string value. If column name is not found, it returns nil and an error.
func TableRowStringToPointerString(col string, row *collections.TableRow[string]) (*string, error) {
	if v, ok := row.Column(col); ok {
		return &v, nil
	}

	return nil, fmt.Errorf("%w: '%s'", collections.ErrUnknownColumnName, col)
}

// TableRowStringToPointerInt This function takes a column name and a table row object as input and
// returns a pointer to an integer value. If the value cannot be converted to an integer or column name
// is not found, it returns an error.
func TableRowStringToPointerInt(col string, row *collections.TableRow[string]) (*int, error) {
	if raw, ok := row.Column(col); ok {
		if i, err := strconv.Atoi(raw); err == nil {
			return &i, nil
		}

		return nil, fmt.Errorf("%w: '%s'", ErrStrToNumberConversion, raw)
	}

	return nil, fmt.Errorf("%w: '%s'", collections.ErrUnknownColumnName, col)
}

// PointerStrToEmptyStrOrValue This function takes a pointer to a string value and returns either an
// empty string or the value itself.
func PointerStrToEmptyStrOrValue(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

// PointerBoolToBool This function takes a pointer to a boolean value and returns either false or the value itself.
func PointerBoolToBool(v *bool) bool {
	if v == nil {
		return false
	}

	return *v
}

// PointerIntToEmptyStrOrValue This function takes a pointer to an integer value
// and returns either an empty string or the value itself.
func PointerIntToEmptyStrOrValue(v *int) string {
	if v == nil {
		return ""
	}

	return strconv.Itoa(*v)
}

// StringFrom This function takes any value as input and returns its string representation.
// If the value is nil, it returns an empty string.
func StringFrom(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprint(v)
}

// StringFromNonZ This function takes a comparable value as input and returns its string representation.
// If the value is zero, it returns an empty string.
func StringFromNonZ[T comparable](value T) string {
	var z T
	if value == z {
		return ""
	}

	return StringFrom(value)
}

// BoolToString This function takes a boolean value as input and returns its string representation.
func BoolToString(v bool) string {
	if v {
		return "true"
	}

	return "false"
}

// StringToBool This function takes a string value as input and returns its boolean representation.
// If the string is not "true", it returns false. Overall, these functions can be useful for handling data conversions
// and formatting in various applications.
func StringToBool(v string) bool {
	return v == "true"
}

// Ptr This function takes a value as input and returns a pointer to that value.
// It is useful when you do not want to create a new variable for a value to pass it as a pointer, for example if you want to pass a literal as a pointer.
func Ptr[T any](v T) *T {
	return &v
}
