package collections

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

func ExampleSyncList() {
	sortedPrint := func(list *SyncList[int]) {
		cp := list.Copy()
		slices.Sort(cp)

		fmt.Println(cp)
	}

	// Создаем экземпляр SyncList.
	sl := NewList[int]()

	// Выполняем операции добавления в список из нескольких горутин.
	wg := sync.WaitGroup{}
	count := 10
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sl.Add(i)
		}(i)
	}

	wg.Wait()

	// Используем сортировку, т.к. в другом случае порядок вывода
	// не будет детерминированным.
	sortedPrint(sl)

	// Записываем в ячейку с индексом i значение i.
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sl.Set(i, i)
		}(i)
	}
	wg.Wait()

	// Получение элемента по индексу.
	sl.Get(3)
	fmt.Println(3)

	// Создание слайса из SyncList.
	slice := sl.Copy()
	fmt.Println(slice)

	// Output:
	// [0 1 2 3 4 5 6 7 8 9]
	// 3
	// [0 1 2 3 4 5 6 7 8 9]
}

func ExampleMultiMap() {
	// Создаем экземпляр MultiMap.
	m := make(MultiMap[string, int])

	// Добавляем список значений для ключа "a".
	m.Append("a", 1, 2, 3)

	fmt.Println(m["a"])

	// Output:
	// [1 2 3]
}

func ExampleSet() {
	// Используем сортировку, чтобы вывод был детерминированным,
	// иначе порядок элементов не будет гарантированно соответствовать
	// порядку добавления.
	sortedPrint := func(values Set[int]) {
		s := values.ToSlice()
		slices.Sort(s)

		// Вывод аналогичный методу String() для типа Set.
		fmt.Printf("[%s]\n", strings.Join(Map(s, func(t int) string { return fmt.Sprint(t) }), ", "))
	}

	// Создание множества.
	set := NewSet(1, 2, 3)
	sortedPrint(set)

	// Добавление нового элемента.
	set.Add(4)
	sortedPrint(set)

	// Добавление нескольких элементов.
	set.Add(5, 6, 7, 8)
	sortedPrint(set)

	// Удаление элемента.
	set.Remove(4)
	sortedPrint(set)

	// Удаление нескольких элементов.
	set.Remove(5, 6, 7)
	sortedPrint(set)

	// Проверка наличия элемента.
	fmt.Println(set.Contains(1))

	// Проверка отсутствия элемента.
	fmt.Println(set.NotContains(10))

	// Проверка содержит ли множество хотя бы один элемент из перечисленных.
	fmt.Println(set.ContainsAny(1, 10))

	// Получение количество элементов.
	fmt.Println(set.Len())

	// Создание слайса из множества.
	set.ToSlice()

	// Output:
	// [1, 2, 3]
	// [1, 2, 3, 4]
	// [1, 2, 3, 4, 5, 6, 7, 8]
	// [1, 2, 3, 5, 6, 7, 8]
	// [1, 2, 3, 8]
	// true
	// 10 true
	// true
	// 4
}

func ExampleMap() {
	// Функция преобразования.
	transform := func(i int) int {
		return i * i
	}

	// Создаем исходный слайc.
	src := []int{1, 2, 3}

	// Применяем преобразование к исходному слайсу.
	result := Map(src, transform)
	fmt.Println(result)

	// Output:
	// [1 4 9]
}

func ExampleFlatMap() {
	// Функция преобразования.
	transform := func(i int) []int {
		return []int{i, i}
	}

	// Создаем исходный слайс.
	src := []int{1, 2, 3}

	// Применяем преобразование к исходному слайсу.
	result := FlatMap(src, transform)
	fmt.Println(result)

	// Output:
	// [1 1 2 2 3 3]
}

func ExampleFilter() {
	// Функция, которая отвечает за фильтрацию.
	condition := func(i int) bool {
		if i%2 == 0 {
			return true
		}

		return false
	}

	// Создаем исходный слайс.
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Применяем фильтрацию к исходному слайсу.
	result := Filter(src, condition)
	fmt.Println(result)

	// Output:
	// [2 4 6 8]
}

func ExampleSplit() {
	// Функция, которая отвечает за разделение исходного слайса на два
	// в первый войдут элементы исходного слайса, для которых condition вернет true, во второй - для которых вернет false.
	condition := func(i int) bool {
		if i%2 == 0 {
			return true
		}

		return false
	}

	// Создаем исходный слайс.
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Применяем функцию-разделение на два слайса по условию.
	result1, result2 := Split(src, condition)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// [2 4 6 8]
	// [1 3 5 7 9]
}

func ExampleSkip() {
	// Создаем исходный слайс.
	src := []int{1, 2, 3, 4, 5}

	// Получаем новый слайс, в котором пропущены первые 3 элемента из исходного.
	result := Skip(src, 3)

	fmt.Println(result)

	// Output:
	// [4 5]
}

func ExampleFind() {
	// Функция-предикат.
	condition := func(i int) bool {
		if i%3 == 0 {
			return true
		}

		return false
	}

	// Создаем исходный слайс.
	src := []int{1, 2, 3, 4, 5}

	// Ищем элемент в исходном слайсе.
	value, ok := Find(src, condition)

	fmt.Println(ok)
	fmt.Println(value)

	// Output:
	// true
	// 3
}

func ExampleMapWithErr() {
	// Функция преобразования, которая всегда возвращает ошибку.
	transformErr := func(i int) (int, error) {
		return 0, fmt.Errorf("some error")
	}
	// Функция преобразования, которая не возвращает ошибок.
	transform := func(i int) (int, error) {
		return i * i, nil
	}

	// Создаем исходный слайс.
	src := []int{1, 2, 3}

	// Применяем преобразование (возвращающее ошибку) к исходному слайсу.
	result, err := MapWithErr(src, transformErr)

	fmt.Println(result)
	fmt.Println(err)

	// Применяем преобразование (не возвращающее ошибки) к исходному слайсу.
	result1, err := MapWithErr(src, transform)

	fmt.Println(result1)
	fmt.Println(err)

	// Output:
	// []
	// some error
	// [1 4 9]
	// <nil>
}

func ExampleMapChunkedWithErr() {
	// Функция преобразования, которая всегда возвращает ошибку.
	transformErr := func([]int) ([]int, error) {
		return nil, fmt.Errorf("some error")
	}

	// Функция преобразования, которая не возвращает ошибок.
	transform := func(src []int) ([]int, error) {
		transformed := make([]int, len(src))
		for i, v := range src {
			transformed[i] = v * v
		}

		return transformed, nil
	}

	// Создаем исходный слайс.
	src := []int{1, 2, 3}

	// Применяем преобразование (возвращающее ошибку) к исходному слайсу.
	result, err := MapChunkedWithErr(src, 2, transformErr)

	fmt.Println(result)
	fmt.Println(err)

	// Применяем преобразование (не возвращающее ошибки) к исходному слайсу.
	result1, err := MapChunkedWithErr(src, 2, transform)

	fmt.Println(result1)
	fmt.Println(err)

	// Output:
	// []
	// some error
	// [1 4 9]
	// <nil>
}

func ExampleForEach() {
	// Функция преобразования.
	act := func(i *int) {
		v := *i
		*i = v * v
	}

	// Исходный слайс.
	src := make([]*int, 3)
	for i := 0; i < len(src); i++ {
		v := i + 1
		src[i] = &v
	}

	// Применяем act для каждого элемента слайса src.
	ForEach(src, act)

	for i := 0; i < len(src); i++ {
		fmt.Println(*src[i])
	}

	// Output:
	// 1
	// 4
	// 9
}

func ExampleForEachWithError() {
	// Функция преобразования, которая всегда возвращает ошибку.
	actErr := func(*int) error {
		return fmt.Errorf("some error")
	}

	// Функция преобразования, которая не возвращает ошибок.
	act := func(i *int) error {
		v := *i
		*i = v * v

		return nil
	}

	// Исходный слайс.
	src := make([]*int, 3)
	for i := 0; i < len(src); i++ {
		v := i + 1
		src[i] = &v
	}

	// Применяем преобразование (возвращающее ошибку) к исходному слайсу.
	err := ForEachWithError(src, actErr)
	fmt.Println(err)

	// Применяем преобразование (не возвращающее ошибки) к исходному слайсу.
	err = ForEachWithError(src, act)
	fmt.Println(err)

	for i := 0; i < len(src); i++ {
		fmt.Println(*src[i])
	}

	// Output:
	// some error
	// <nil>
	// 1
	// 4
	// 9
}

func ExampleForEachChunk() {
	// Функция преобразования.
	act := func(values []*int) {
		for i := 0; i < len(values); i++ {
			v := *values[i]
			*values[i] = v * v
		}
	}

	// Исходный слайс.
	src := make([]*int, 5)
	for i := 0; i < len(src); i++ {
		v := i + 1
		src[i] = &v
	}

	// Применяем преобразование к исходному слайсу по чанкам.
	ForEachChunk(src, 2, act)

	for i := 0; i < len(src); i++ {
		fmt.Println(*src[i])
	}

	// Output:
	// 1
	// 4
	// 9
	// 16
	// 25
}

func ExampleForEachChunkWithErr() {
	// Функция преобразования, которая всегда возвращает ошибку.
	actErr := func(values []*int) error {
		return fmt.Errorf("some error")
	}

	// Функция преобразования, которая не возвращает ошибок.
	act := func(values []*int) error {
		for i := 0; i < len(values); i++ {
			v := *values[i]
			*values[i] = v * v
		}

		return nil
	}

	// Исходный слайс.
	src := make([]*int, 5)
	for i := 0; i < len(src); i++ {
		v := i + 1
		src[i] = &v
	}

	// Применяем преобразование (возвращающее ошибку) к исходному слайсу.
	err := ForEachChunkWithErr(src, 2, actErr)
	fmt.Println(err)

	// Применяем преобразование (не возвращающее ошибки) к исходному слайсу.
	err = ForEachChunkWithErr(src, 2, act)
	fmt.Println(err)

	for i := 0; i < len(src); i++ {
		fmt.Println(*src[i])
	}

	// Output:
	// some error
	// <nil>
	// 1
	// 4
	// 9
	// 16
	// 25
}

func ExampleNonNil() {
	// Объявляем и заполняем исходный слайс.
	src := make([]*int, 9)
	for i := 0; i < len(src); i++ {
		v := i + 1
		if v%2 == 0 {
			src[i] = &v
		} else {
			src[i] = nil
		}
	}

	// Применяем фильтрацию, которая возвращает новый слайc содержащий только не nil элементы.
	result := NonNil(src)
	fmt.Println(result)

	// Output:
	// [2 4 6 8]
}

func ExampleExceptIndexes() {
	// Объявляем исходный слайс.
	src := []int{1, 2, 3, 4, 5, 6, 7}

	// Объявляем множество индексов, которые нужно исключить.
	indexes := NewSet[int](0, 2, 4)

	// Получаем новый слайс, который не содержит элементы с индексами из indexes.
	result := ExceptIndexes(src, indexes)
	fmt.Println(result)

	// Output:
	// [2 4 6 7]
}

func ExampleKeysOfMap() {
	// Объявляем мапу.
	m := map[int]string{1: "a", 2: "b", 3: "c"}

	// Получаем слайс ключей мапы.
	keys := KeysOfMap(m)

	slices.Sort(keys)
	fmt.Println(keys)

	// Output:
	// [1 2 3]
}

func ExampleValuesOfMap() {
	// Объявляем мапу.
	m := map[int]string{1: "a", 2: "b", 3: "c"}

	// Получаем слайс значений мапы.
	values := ValuesOfMap(m)

	slices.Sort(values)
	fmt.Println(values)

	// Output:
	// [a b c]
}

func ExampleJoin() {
	// Объявляем двумерный слайс и заполняем его.
	s := make([][]int, 3)
	for i := 0; i < len(s); i++ {
		si := make([]int, 3)
		for j := 0; j < len(si); j++ {
			si[j] = i
		}

		s[i] = si
	}
	fmt.Println(s)

	joined := Join(s...)

	fmt.Println(joined)

	// Output:
	// [[0 0 0] [1 1 1] [2 2 2]]
	// [0 0 0 1 1 1 2 2 2]
}

func ExampleSingle() {
	single := Single(42)

	fmt.Println(single)

	// Output:
	// [42]
}

func ExampleGroup() {
	// Функция, которая определяет ключ для элемента мапы по значению.
	key := func(i int) string {
		if i%2 == 0 {
			return "even"
		} else {
			return "odd"
		}
	}

	// Объявляем исходный слайс.
	src := []int{1, 2, 3, 4, 5}

	// Преобразуем исходный слайс в мапу используя группировку.
	groups := Group(src, key)
	fmt.Println(groups)

	// Output:
	// map[even:[2 4] odd:[1 3 5]]
}

func ExampleDictionary() {
	// Функция, которая определяет ключ для элемента мапы по значению.
	key := func(i int) int {
		return i
	}

	// Объявляем исходный слайс.
	src := []int{1, 2, 3, 4, 5}

	// Преобразуем исходный слайс в мапу.
	result := Dictionary(src, key)
	fmt.Println(result)

	// Output:
	// map[1:1 2:2 3:3 4:4 5:5]
}

func ExampleChunked() {
	// Объявляем исходный слайс.
	src := []int{1, 2, 3, 4, 5, 6, 7}

	// Разбиваем исходный слайс на чанки.
	result := Chunked(src, 3)
	fmt.Println(result)

	// Output:
	// [[1 2 3] [4 5 6] [7]]
}

func ExampleNewSorter() {
	// Объявляем слайс для сортировки.
	listString := []string{"foo123", "bar", "1"}

	// Сортируем слайс по длине строки.
	sorterString := NewSorter(listString, func(a, b string) bool {
		return len(a) < len(b)
	})
	sort.Sort(sorterString)

	fmt.Println(listString)

	// Output:
	// [1 bar foo123]
}

func ExampleTable() {
	// Создаем новую таблицу с 2 столбцами и колонками типа string.
	t := NewTable[string]("name", "type")

	// Получим список столбцов.
	columns := t.Columns()
	fmt.Println(columns)

	// Добавим новый столбец.
	t.AddColumns("weight")
	fmt.Println(t.Columns())

	// Добавим 2 строки.
	t.AddRow("orange", "fruit", "200g")
	t.AddRow("carrot", "vegetable", "350g")

	// Получим размер таблицы (количество строк).
	size := t.Size()
	fmt.Println(size)

	// Получим строку по индексу.
	row := t.Get(0)
	// Получим все значения в строке.
	values := row.GetValues()
	fmt.Println(values)

	// Установим новые значения для строки.
	row.SetValues("apple", "fruit", "100g")
	fmt.Println(row.GetValues())

	// Получим значение столбца 'name' для строки.
	name, ok := row.Column("name")
	fmt.Println(name, ok)

	// Получим индекс строки.
	index := row.Index()
	fmt.Println(index)

	// Output:
	// [name type]
	// [name type weight]
	// 2
	// [orange fruit 200g]
	// [apple fruit 100g]
	// apple true
	// 0
}
