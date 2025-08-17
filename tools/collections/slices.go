package collections

// Map возвращает копию слайса, преобразованную с помощью функции transform.
func Map[T, R any](slice []T, transform func(T) R) []R {
	destination := make([]R, len(slice))

	for i := 0; i < len(slice); i++ {
		destination[i] = transform(slice[i])
	}

	return destination
}

// FlatMap возвращает преобразованный слайс, преобразованный с помощью функции transform.
// В отличие от Map, функция-преобразования может вернуть несколько значений.
func FlatMap[T, R any](slice []T, transform func(T) []R) []R {
	result := make([]R, 0)

	for _, t := range slice {
		result = append(result, transform(t)...)
	}

	return result
}

// Filter возвращает новый слайс, в который входят элементы s, для которых функция f вернет true.
func Filter[T any](slice []T, condition func(T) bool) []T {
	destination := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if condition(slice[i]) {
			destination = append(destination, slice[i])
		}
	}

	return destination
}

// Split возвращает два слайса, в первый войдут элементы s, для которых f вернет true, во второй - для которых вернет false.
func Split[T any](slice []T, condition func(T) bool) ([]T, []T) {
	slice1 := make([]T, 0)
	slice2 := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if condition(slice[i]) {
			slice1 = append(slice1, slice[i])
		} else {
			slice2 = append(slice2, slice[i])
		}
	}

	return slice1, slice2
}

// Skip возвращает новый слайс, без первых size элементов. Если size больше размера слайса, вернет nil.
func Skip[T any](slice []T, size int) []T {
	if size > len(slice) {
		return nil
	}

	return slice[size:]
}

// Find возвращает первый попавшийся элемент в слайсе, который подходит под условие предиката.
func Find[T any](slice []T, condition func(T) bool) (T, bool) {
	var blank T

	for _, value := range slice {
		if condition(value) {
			return value, true
		}
	}

	return blank, false
}

// MapWithErr возвращает копию слайса, преобразованную с помощью функции f, если f вернет error, то MapWithErr так же
// вернет error.
func MapWithErr[T any, R any](source []T, transform func(T) (R, error)) ([]R, error) {
	destination := make([]R, len(source))

	for index := 0; index < len(source); index++ {
		v, err := transform(source[index])
		if err != nil {
			return nil, err
		}

		destination[index] = v
	}

	return destination, nil
}

// MapChunkedWithErr возвращает преобразованный слайс, слайс преобразуется через ForEachChunkWithErr.
func MapChunkedWithErr[T any, R any](slice []T, size int, transform func([]T) ([]R, error)) ([]R, error) {
	sliceSize := len(slice)
	if sliceSize <= size {
		return transform(slice)
	}

	results := make([]R, 0, sliceSize)

	return results, ForEachChunkWithErr(slice, size, func(chunk []T) error {
		r, err := transform(chunk)
		if err == nil {
			results = append(results, r...)
		}

		return err
	})
}

// ForEach выполняет f для каждого элемента слайса s.
func ForEach[T any](slice []T, act func(T)) {
	for i := 0; i < len(slice); i++ {
		act(slice[i])
	}
}

// ForEachWithError	выполняет f для каждого элемента слайса s,
// если f вернет error, то ForEachWithError так же вернет error.
func ForEachWithError[T any](s []T, f func(T) error) error {
	for i := 0; i < len(s); i++ {
		if err := f(s[i]); err != nil {
			return err
		}
	}

	return nil
}

// ForEachChunk для каждого чанка фиксированного size выполнить f, игнорируя ошибки.
func ForEachChunk[T any](slice []T, size int, act func([]T)) {
	_ = ForEachChunkWithErr(slice, size, func(chunk []T) error {
		act(chunk)

		return nil
	})
}

// ForEachChunkWithErr для каждого чанка фиксированного size выполнить f.
// Если в чанке произошла ошибка, вернуть error.
func ForEachChunkWithErr[T any](slice []T, size int, act func([]T) error) error {
	sliceSize := len(slice)

	for cursor := 0; cursor < sliceSize; {
		end := cursor + size
		if end > sliceSize {
			end = sliceSize
		}

		if err := act(slice[cursor:end]); err != nil {
			return err
		}

		cursor = end
	}

	return nil
}

// NonNil свернет слайс со всеми не nil элементами из values.
func NonNil[T any](values []*T) []T {
	result := make([]T, 0)

	for _, value := range values {
		if value != nil {
			result = append(result, *value)
		}
	}

	return result
}

// ExceptIndexes вернет слайс с элементами s, за исключением тех, индексы которых в indexes.
func ExceptIndexes[T any](src []T, indexes Set[int]) []T {
	if len(indexes) == 0 {
		return src
	}

	dest := make([]T, 0)

	for i, v := range src {
		if !indexes.Contains(i) {
			dest = append(dest, v)
		}
	}

	return dest
}

// KeysOfMap возвращает слайс ключей мапы m.
func KeysOfMap[T comparable, K any](dict map[T]K) []T {
	result := make([]T, 0, len(dict))

	for k := range dict {
		result = append(result, k)
	}

	return result
}

// ValuesOfMap возвращает слайс значений мапы m.
func ValuesOfMap[T comparable, K any](dict map[T]K) []K {
	result := make([]K, 0, len(dict))

	for _, v := range dict {
		result = append(result, v)
	}

	return result
}

// Join возвращает слайс, созданный из последовательного присоединения входных слайсов.
func Join[T any](slices ...[]T) []T {
	c := 0
	for _, slice := range slices {
		c += len(slice)
	}

	result := make([]T, 0, c)
	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}

// Single создает из элемента слайс единичного размера.
func Single[T any](t T) []T {
	return []T{t}
}

// Group преобразует слайс типа V по функции func(V) K в map[K][]V.
func Group[K comparable, V any](list []V, key func(V) K) map[K][]V {
	result := make(map[K][]V)

	for _, item := range list {
		keyValue := key(item)

		if _, ok := result[keyValue]; !ok {
			result[keyValue] = make([]V, 0, 1)
		}

		result[keyValue] = append(result[keyValue], item)
	}

	return result
}

// Dictionary преобразует слайс типа V по функции func(V) K в map[K]V.
func Dictionary[K comparable, V any](list []V, key func(V) K) map[K]V {
	result := make(map[K]V)

	for _, item := range list {
		result[key(item)] = item
	}

	return result
}

// Chunked разбивает слайс на чанки указанного размера.
func Chunked[T any](slice []T, size int) [][]T {
	if size == 0 {
		return nil
	}

	chunks := make([][]T, 0, len(slice)/size+1)

	ForEachChunk(slice, size, func(chunk []T) {
		chunks = append(chunks, chunk)
	})

	return chunks
}
