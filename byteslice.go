package readonly

import (
	"io"
	"reflect"
	"unsafe"
)

// NewByteSlice constructor for ByteSlice.
// Accepts a string or slice of bytes as input, avoiding allocations.
func NewByteSlice[T ~string | ~[]byte](src T) (b ByteSlice) {
	// Here it is guaranteed that the byte slice will be unchanged,
	// since the interface is read-only.
	b.s = *(*[]byte)(unsafe.Pointer(&src))
	(*reflect.SliceHeader)(unsafe.Pointer(&b.s)).Cap = len(src)
	return b
}

// ByteSlice wrapper over []byte that limits the interface to read-only.
type ByteSlice struct{ Slice[byte] }

// String string(b) equivalent for byte slice.
func (b ByteSlice) String() string { return string(b.s) }

// Equal s1 == s2 equivalent for string.
func (b ByteSlice) Equal(bb []byte) bool {
	// Neither cmd/compile nor gccgo allocates for these string conversions.
	return string(b.s) == string(bb)
}

// ReadAt implements io.ReaderAt.
func (b ByteSlice) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(b.s)) {
		return 0, io.EOF
	}
	return copy(p, b.s[off:]), nil
}

// WriteTo implements io.WriterTo.
// w must not modify the slice data, even temporarily, see io.Writer.
func (b ByteSlice) WriteTo(w io.Writer) (n int64, err error) {
	if len(b.s) == 0 {
		return 0, nil
	}

	m, err := w.Write(b.s)
	switch n = int64(m); {
	case err != nil:
		return n, err
	case m != len(b.s):
		return n, io.ErrShortWrite
	}
	return n, nil
}
