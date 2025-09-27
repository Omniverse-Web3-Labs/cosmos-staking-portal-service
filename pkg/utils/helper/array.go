package helper

import (
	"math/rand"
	"strings"
)

type Addable interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | string
}

func InArray[T comparable](a []T, b T) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// Traverses the  array
func ArrayMap[T any, K any](a []T, f func(s T) K) []K {
	b := make([]K, 0)
	for i := range a {
		b = append(b, f(a[i]))
	}
	return b
}

// Remove item
func ArrayRemove[T comparable](a []T, s T) []T {
	b := make([]T, 0)
	for _, item := range a {
		if item != s {
			b = append(b, item)
		}
	}
	return b
}

func ArrayFilter[T comparable](a []T, f func(s T) bool) []T {
	b := make([]T, 0)
	for i := range a {
		if f(a[i]) {
			b = append(b, a[i])
		}
	}
	return b
}

func ArraySum[T Addable](a []T) T {
	var sum T
	for _, v := range a {
		sum = v + sum
	}
	return sum
}

func ArrayWalk[T any](a []T, f func(item T)) {
	for i := range a {
		f(a[i])
	}
}

func ArrayToMap[T any, K comparable](a []T, key func(item T) K) map[K]T {
	h := make(map[K]T)
	for i := range a {
		h[key(a[i])] = a[i]
	}
	return h
}

func ArrayColumns[T any, K any](a []T, column func(item T) K) []K {
	b := make([]K, 0)
	for i := range a {
		b = append(b, column(a[i]))
	}
	return b
}

func ArrayRandom[T any](a []T) T {
	if len(a) == 0 {
		var t T
		return t
	}
	n := len(a)
	return a[rand.Intn(n)]
}

func ArraySlice[T any](a []T, start, length int) []T {
	n := len(a)
	if n == 0 {
		return a
	}
	if start >= n {
		return []T{}
	} else if start < 0 {
		start = n + start
		if start < 0 {
			start = 0
		}
	}
	if start+length > n {
		length = n - start
	}
	return a[start : start+length]
}

func Shuffle[T any](a []T) []T {
	length := len(a)
	for i := range a {
		j := rand.Intn(length)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func Join[T any](a []T, sep string) string {
	var b []string
	for i := range a {
		b = append(b, StringVal(a[i]))
	}
	return strings.Join(b, sep)
}

func ArrayUnique[T comparable](a []T) []T {
	b := make([]T, 0, len(a))
	m := make(map[T]bool)

	for _, item := range a {
		if _, ok := m[item]; !ok {
			b = append(b, item)
			m[item] = true
		}
	}

	return b
}
