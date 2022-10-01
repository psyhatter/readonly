//nolint:ireturn
package readonly

// NewSlice returns a slice interface limited to read-only methods.
func NewSlice[T any](s []T) Slice[T] { return Slice[T]{s: s} }

// Slice wrapper over a built-in slice that limits the interface
// to read-only.
type Slice[T any] struct{ s []T }

// IsNil equivalent to s == nil.
func (s Slice[T]) IsNil() bool { return s.s == nil }

// Len equivalent to len(s).
func (s Slice[T]) Len() int { return len(s.s) }

// Cap equivalent to cap(s).
func (s Slice[T]) Cap() int { return cap(s.s) }

// Get equivalent to v := s[index].
func (s Slice[T]) Get(index int) (v T) { return s.s[index] }

// Range equivalent to read-only for range loop.
// Does nothing if f == nil.
// Breaks the loop if next == false.
// An order of magnitude slower than the built-in slice loop, for
// optimizations, you can use slice index access (see benchmarks).
func (s Slice[T]) Range(f func(index int, val T) (next bool)) {
	if f != nil {
		for index := range s.s {
			if !f(index, s.s[index]) {
				return
			}
		}
	}
}

// Copy returns a new copy of the built-in slice.
// As fast as inline copy to new slice, but faster than copying to a new
// slice with (s Slice) CopyTo.
func (s Slice[T]) Copy() []T { return append([]T(nil), s.s...) }

// CopyTo copies elements from a source slice into a destination slice.
// The source and destination may overlap. Copy returns the number of
// elements copied, which will be the minimum of (Slice) Len() and len(dst).
func (s Slice[T]) CopyTo(dst []T) int { return copy(dst, s.s) }

// Append appends elements to the end of dst and returns the updated slice.
func (s Slice[T]) Append(dst []T) []T { return append(dst, s.s...) }

// AppendInto adds elements to the end of the slice located at the dst
// pointer and places the new slice at the dst pointer.
// Does nothing if dst == nil.
func (s Slice[T]) AppendInto(to *[]T) {
	if to != nil {
		*to = append(*to, s.s...)
	}
}

// StartAfter equivalent to s[i:].
func (s Slice[T]) StartAfter(i int) Slice[T] { return Slice[T]{s.s[i:]} }

// EndBefore equivalent to s[:i].
func (s Slice[T]) EndBefore(i int) Slice[T] { return Slice[T]{s.s[:i]} }

// Slice equivalent to s[start:end].
func (s Slice[T]) Slice(start, end int) Slice[T] { return Slice[T]{s.s[start:end]} }

// Slice3 equivalent to s[i:j:k].
func (s Slice[T]) Slice3(i, j, k int) Slice[T] { return Slice[T]{s.s[i:j:k]} }
