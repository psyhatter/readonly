package readonly_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/psyhatter/readonly"
)

func ExampleSlice_IsNil() {
	s1, s2 := readonly.NewSlice[int](nil), readonly.NewSlice([]int{})
	fmt.Println(s1.IsNil(), s2.IsNil())

	// Output:
	// true false
}

func ExampleSlice_Len() {
	s := readonly.NewSlice([]int{1, 2})

	fmt.Println(s.Len())
	// Output:
	// 2
}

func ExampleSlice_Cap() {
	s := readonly.NewSlice(make([]int, 0, 2))

	fmt.Println(s.Cap())
	// Output:
	// 2
}

func ExampleSlice_Get() {
	s := readonly.NewSlice([]int{1, 2})

	fmt.Println(s.Get(0))
	fmt.Println(s.Get(1))
	// Output:
	// 1
	// 2
}

func ExampleSlice_Range() {
	s := readonly.NewSlice([]int{0, 1})
	s.Range(func(i int, v int) (next bool) {
		fmt.Println(i, v)
		return true
	})
	s.Range(func(i int, v int) (next bool) {
		fmt.Println("this happened") // will be printed only once.
		return false
	})
	// Output:
	// 0 0
	// 1 1
	// this happened
}

func ExampleSlice_CopyTo() {
	src := readonly.NewSlice([]int{0, 1})
	dst := make([]int, src.Len())

	fmt.Println(src.CopyTo(dst), dst)
	// Output:
	// 2 [0 1]
}

func ExampleSlice_Append() {
	src := readonly.NewSlice([]int{0, 1})
	dst := make([]int, 0, src.Len())

	fmt.Println(src.Append(dst))
	// Output:
	// [0 1]
}

func ExampleSlice_AppendInto() {
	src := readonly.NewSlice([]int{0, 1})
	dst := make([]int, 0, src.Len())
	src.AppendInto(&dst)
	src.AppendInto(nil) // don't panic.

	fmt.Println(dst)
	// Output:
	// [0 1]
}

func ExampleSlice_StartAfter() {
	s := readonly.NewSlice([]int{0, 1})

	fmt.Println(s.StartAfter(1))
	// Output:
	// {[1]}
}

func ExampleSlice_EndBefore() {
	s := readonly.NewSlice([]int{0, 1})

	fmt.Println(s.EndBefore(1))

	// Output:
	// {[0]}
}

func ExampleSlice_Slice() {
	a := make([]int, 0, 5)
	a = append(a, 0, 1, 2, 3)
	s := readonly.NewSlice(a)

	fmt.Println(s.Slice(1, 3), s.Slice(s.Cap()-1, s.Cap()))

	// Output:
	// {[1 2]} {[0]}
}

func ExampleSlice_Slice3() {
	a := make([]int, 0, 5)
	a = append(a, 0, 1, 2, 3)
	s := readonly.NewSlice(a)

	sliced := s.Slice3(1, 3, 4)
	fmt.Println(sliced, sliced.Len(), sliced.Cap())

	// Output:
	// {[1 2]} 2 3
}

func ExampleSlice_Copy() {
	s := readonly.NewSlice([]int{1, 2, 3})

	fmt.Println(s.Copy())
	// Output:
	// [1 2 3]
}

const limit = 100000

//nolint:funlen
func BenchmarkSlice_Range(b *testing.B) {
	s := rand.Perm(limit)

	b.Run("built-in", func(b *testing.B) {

		// Fastest way to iterate through all elements of a built-in slice.
		b.Run("fast", func(b *testing.B) {
			s := s
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var count int
				for i := 0; i < len(s); i++ {
					count += i + s[i]
				}
			}
		})

		// Usually a little slower than built-in fast range.
		b.Run("range", func(b *testing.B) {
			s := s
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var count int
				for i, j := range s {
					count += i + j
				}
			}
		})
	})

	b.Run("readonly", func(b *testing.B) {

		// Fastest way to iterate through all the elements of a readonly.Slice.
		// Usually a little slower than built-in fast range.
		b.Run("fast", func(b *testing.B) {
			s := readonly.NewSlice(s)
			s.Get(s.Len() - 1) // avoid boundary check.
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var count int
				for i := 0; i < s.Len(); i++ {
					count += i + s.Get(i)
				}
			}
		})

		// Perhaps a more convenient interface for iterating through all elements
		// of a readonly.Slice.
		// Usually 10 times slower than built-in fast range.
		b.Run("range", func(b *testing.B) {
			s := readonly.NewSlice(s)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var count int
				s.Range(func(i, j int) bool { count += i + j; return true })
			}
		})

		// Iterating over all elements with constant copying is usually 10 times
		// slower than built-in fast range.
		b.Run("range with copying", func(b *testing.B) {
			s := readonly.NewSlice(s)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var count int
				for i, j := range s.Copy() {
					count += i + j
				}
			}
		})
	})
}

func BenchmarkSlice_Copy(b *testing.B) {
	s := rand.Perm(limit)

	// A quick way to get a copy of a built-in slice.
	b.Run("built-in", func(b *testing.B) {
		s := s
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			dst := make([]int, limit)
			copy(dst, s)
		}
	})

	b.Run("readonly", func(b *testing.B) {

		// A quick way to get copy of readonly.Slice.
		// Usually as fast as built-in copy.
		b.Run("copy", func(b *testing.B) {
			s := readonly.NewSlice(s)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				s.Copy()
			}
		})

		// Usually uncritically slower than built-in copy.
		b.Run("copy to", func(b *testing.B) {
			s := readonly.NewSlice(s)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				dst := make([]int, limit)
				s.CopyTo(dst)
			}
		})
	})
}

func BenchmarkSlice_Append(b *testing.B) {
	s := rand.Perm(limit)

	// A quick way to append a built-in slice.
	b.Run("built-in", func(b *testing.B) {
		s := s
		b.ResetTimer()

		dst := make([]int, limit)
		for i := 0; i < b.N; i++ {
			dst = append(dst[:0], s...)
		}
	})

	// A quick way to append readonly.Slice.
	// Usually as fast as built-in append.
	b.Run("readonly", func(b *testing.B) {
		s := readonly.NewSlice(s)
		b.ResetTimer()

		dst := make([]int, limit)
		for i := 0; i < b.N; i++ {
			dst = s.Append(dst[:0])
		}
	})
}
