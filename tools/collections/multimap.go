package collections

// MultiMap представляет собой структуру обобщенной Map, где значением является список значений.
type MultiMap[K comparable, V any] map[K][]V

// Append добавить по ключу, список значений.
func (m MultiMap[K, V]) Append(key K, values ...V) {
	if _, ok := m[key]; !ok {
		m[key] = make([]V, 0, len(values))
	}
	m[key] = append(m[key], values...)
}

// MergeMapsValue Функция для объединения значений всех ключей из неограниченного количества карт.
func MergeMapsValue(maps ...map[string]string) []string {
	var mergedValues []string

	for _, mp := range maps {
		for _, value := range mp {
			mergedValues = append(mergedValues, value)
		}
	}

	return mergedValues
}
