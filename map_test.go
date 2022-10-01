package readonly_test

import (
	"fmt"
	"testing"

	"github.com/psyhatter/readonly"
)

func ExampleMap_IsNil() {
	m1, m2 := readonly.NewMap[string, int](nil), readonly.NewMap(map[string]int{})
	fmt.Println(m1.IsNil(), m2.IsNil())

	// Output:
	// true false
}

func ExampleMap_Len() {
	m := readonly.NewMap(map[string]int{"1": 1, "2": 2})
	fmt.Println(m.Len())

	// Output:
	// 2
}

func ExampleMap_Has() {
	m := readonly.NewMap(map[string]int{"1": 1, "2": 2})

	fmt.Println(m.Has("1"))
	fmt.Println(m.Has("2"))
	fmt.Println(m.Has("3"))
	// Output:
	// true
	// true
	// false
}

func ExampleMap_Get() {
	m := readonly.NewMap(map[string]int{"1": 1, "2": 2})

	fmt.Println(m.Get("1"))
	fmt.Println(m.Get("2"))
	fmt.Println(m.Get("3"))
	// Output:
	// 1
	// 2
	// 0
}

func ExampleMap_Get2() {
	m := readonly.NewMap(map[string]int{"1": 1, "2": 2})

	fmt.Println(m.Get2("1"))
	fmt.Println(m.Get2("2"))
	fmt.Println(m.Get2("3"))
	// Output:
	// 1 true
	// 2 true
	// 0 false
}

func ExampleMap_Range() {
	m := readonly.NewMap(map[string]int{"1": 1, "2": 2})
	m.Range(func(k string, v int) (next bool) {
		fmt.Println(k, v)
		return true
	})
	m.Range(func(k string, v int) (next bool) {
		fmt.Println("this happened") // will be printed only once.
		return false
	})
	m.Range(nil) // do nothing.

	// Unordered output:
	// 1 1
	// 2 2
	// this happened
}

var m = func() map[int]int {
	m := make(map[int]int, limit)
	for i := 0; i < limit; i++ {
		m[i] = i
	}
	return m
}()

func BenchmarkMap_Range(b *testing.B) {
	b.Run("built-in", func(b *testing.B) {
		m := m
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var count int
			for key, val := range m {
				count += key + val
			}
		}
	})
	b.Run("readonly", func(b *testing.B) {
		m := readonly.NewMap(m)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var count int
			m.Range(func(key, val int) (next bool) {
				count += key + val
				return true
			})
		}
	})
}
