package collections_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wal1251/pkg/tools/collections"
)

func TestNewSorter_String(t *testing.T) {
	listString := []string{"foo", "bar", "1aa"}

	sorterString := collections.NewSorter(listString, func(a, b string) bool {
		return a < b
	})
	sort.Sort(sorterString)

	assert.Equal(t, []string{"1aa", "bar", "foo"}, listString)
}

func TestNewSorter_Int(t *testing.T) {
	listInt := []int{12, -10, 0}

	sorterInt := collections.NewSorter(listInt, func(a, b int) bool {
		return a < b
	})
	sort.Sort(sorterInt)

	assert.Equal(t, []int{-10, 0, 12}, listInt)
}

func TestSorter_SortWith(t *testing.T) {
	listString := []string{"foo00", "bar111", "1"}

	collections.SortWith(listString, func(elem string) int {
		return len(elem)
	})

	assert.Equal(t, []string{"1", "foo00", "bar111"}, listString)
}
