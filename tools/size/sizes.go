package size

import "fmt"

const (
	k       = 1 << 10
	B  Size = 1      //nolint: varnamelen
	KB      = k * B  //nolint: varnamelen
	MB      = k * KB //nolint: varnamelen
	GB      = k * MB //nolint: varnamelen
	TB      = k * GB //nolint: varnamelen

	unitNameB  = "B"
	unitNameKB = "KB"
	unitNameMB = "MB"
	unitNameGB = "GB"
	unitNameTB = "TB"
)

// Size определяет тип для представления размера в байтах.
type Size int64

// Count возвращает количество единиц измерения в переданном размере.
// Если unit равен 0, возвращает исходное количество байт.
func (s Size) Count(unit Size) int64 {
	if unit == 0 {
		return int64(s)
	}

	return int64(s / unit)
}

// Int конвертирует размер в int.
func (s Size) Int() int {
	return int(s)
}

// Int64 конвертирует размер в int64.
func (s Size) Int64() int64 {
	return int64(s)
}

// String возвращает строковое представление размера, автоматически выбирая подходящую единицу измерения.
func (s Size) String() string {
	switch {
	case s >= TB:
		return sizeToString(s, TB, unitNameTB)
	case s >= GB:
		return sizeToString(s, GB, unitNameGB)
	case s >= MB:
		return sizeToString(s, MB, unitNameMB)
	case s >= KB:
		return sizeToString(s, KB, unitNameKB)
	}

	return sizeToString(s, B, unitNameB)
}

// sizeToString вспомогательная функция для преобразования размера в строку с учетом единицы измерения.
func sizeToString(size, unit Size, unitName string) string {
	if unit == 0 {
		return sizeToString(size, B, unitNameB)
	}

	targetUnits := size / unit
	remainUnits := size % unit

	if remainUnits == 0 {
		return fmt.Sprintf("%d %s", targetUnits, unitName)
	}

	partialUnits := float64(remainUnits) / float64(unit)

	return fmt.Sprintf("%.1f %s", float64(targetUnits)+partialUnits, unitName)
}

// Bytes конвертирует размер из заданной единицы измерения в байты.
func Bytes(size int, unit Size) int {
	return Make(size, unit).Int()
}

// Make создает объект Size из заданного размера и единицы измерения.
func Make(size int, unit Size) Size {
	return Size(size) * unit
}
