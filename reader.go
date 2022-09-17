package readonly

import (
	"io"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

// NewReader returns a new Reader reading from src.
// It is similar to strings.NewReader or bytes.NewReader but more
// efficient.
func NewReader[T ~string | ~[]byte | ByteSlice](src T) *Reader {
	return &Reader{s: *(*string)(unsafe.Pointer(&src))}
}

// Reader implements the io.Reader, io.ByteReader, io.RuneReader and
// io.WriterTo interfaces by reading from a string.
// The zero value for Reader operates like a Reader of an empty string,
// nil byte slice or an empty byte slice.
type Reader struct{ s string }

// Len returns the number of bytes of the unread portion of the
// string.
func (r *Reader) Len() int { return len(r.s) }

// Read implements the io.Reader interface.
func (r *Reader) Read(p []byte) (n int, err error) {
	if len(r.s) == 0 {
		return 0, io.EOF
	}
	n = copy(p, r.s)
	r.s = r.s[n:]
	return n, nil
}

// ReadByte implements the io.ByteReader interface.
func (r *Reader) ReadByte() (b byte, err error) {
	if len(r.s) == 0 {
		return 0, io.EOF
	}
	b, r.s = r.s[0], r.s[1:]
	return b, nil
}

// ReadRune implements the io.RuneReader interface.
func (r *Reader) ReadRune() (ch rune, size int, err error) {
	if len(r.s) == 0 {
		return 0, 0, io.EOF
	}
	if r.s[0] < utf8.RuneSelf {
		ch, size, r.s = rune(r.s[0]), 1, r.s[1:]
		return ch, size, nil
	}
	ch, size = utf8.DecodeRuneInString(r.s)
	r.s = r.s[size:]
	return ch, size, nil
}

// WriteTo implements the io.WriterTo interface.
// w must not modify the slice data, even temporarily, see io.Writer.
func (r Reader) WriteTo(w io.Writer) (n int64, err error) {
	if len(r.s) == 0 {
		return 0, nil
	}

	// io.Writer has to guarantee that the slice of bytes will be
	// unchanged.
	b := *(*[]byte)(unsafe.Pointer(&r.s))
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = len(r.s)
	m, err := w.Write(b)
	switch n = int64(m); {
	case err != nil:
		return n, err
	case m != len(r.s):
		return n, io.ErrShortWrite
	}
	return n, nil
}

// ResetString resets the Reader to be reading from s.
func (r *Reader) ResetString(s string) { r.s = s }

// ResetBytes resets the Reader to be reading from b.
func (r *Reader) ResetBytes(b []byte) { r.s = *(*string)(unsafe.Pointer(&b)) }

// ResetByteSlice resets the Reader to be reading from b.
func (r *Reader) ResetByteSlice(b ByteSlice) { r.s = *(*string)(unsafe.Pointer(&b)) }
