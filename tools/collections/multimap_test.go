package collections_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wal1251/pkg/tools/collections"
)

func TestMultiMap_Append(t *testing.T) {
	multiMap := make(collections.MultiMap[string, int])

	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			multiMap.Append("0", i)
		} else if i%3 == 1 {
			multiMap.Append("1", i)
		} else {
			multiMap.Append("2", i)
		}
	}

	require.ElementsMatch(t, []int{0, 3, 6, 9}, multiMap["0"])
	require.ElementsMatch(t, []int{1, 4, 7}, multiMap["1"])
	require.ElementsMatch(t, []int{2, 5, 8}, multiMap["2"])
}
