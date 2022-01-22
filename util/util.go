package util

import (
	"constraints"
	"sort"
)

func SortResourcesById[T any](r []T, idGetter func(T) string) []T {
	sort.Slice(
		r, func(a, b int) bool {
			return idGetter(r[a]) < idGetter(r[b])
		},
	)
	return r
}

func GetSortedMapValues[T any](r map[string]T) []T {
	var keys []string
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var result []T
	for _, k := range keys {
		result = append(result, r[k])
	}
	return result
}

func MaxBy[K constraints.Ordered, V any, P constraints.Ordered](r map[K]V, selector func(V) P) K {
	var currentMax K
	for k, v := range r {
		if _, ok := r[currentMax]; !ok {
			currentMax = k
		}
		if selector(r[currentMax]) < selector(v) || (selector(r[currentMax]) == selector(v) && currentMax < k) {
			currentMax = k
		}
	}
	return currentMax
}

func Contains[T comparable](list []T, a T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Keys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
