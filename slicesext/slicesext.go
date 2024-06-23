package slicesext

import (
	"cmp"
	"slices"
)

// ToStructMap converts the slice to a map with struct{} values.
func ToStructMap[T comparable](s []T) map[T]struct{} {
	m := make(map[T]struct{}, len(s))
	for _, e := range s {
		m[e] = struct{}{}
	}
	return m
}

// MapKeysToSortedSlice converts the map's keys to a sorted slice.
func MapKeysToSortedSlice[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	s := MapKeysToSlice(m)
	slices.Sort(s)
	return s
}

// MapKeysToSlice converts the map's keys to a slice.
func MapKeysToSlice[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}
