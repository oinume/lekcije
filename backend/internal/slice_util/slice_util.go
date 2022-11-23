package slice_util

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func Sort[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}
