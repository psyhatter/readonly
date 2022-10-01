# readonly

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/psyhatter/readonly)
[![Go Report Card](https://goreportcard.com/badge/github.com/psyhatter/readonly)](https://goreportcard.com/report/github.com/psyhatter/readonly)

Package readonly provides an interface to some built-in
container types (such as slices, maps, and channels) that
allow them to be read-only. The package uses generics so
that they can be used for any types. Also, byte sequences
(strings or slices of bytes) are additionally processed
through the interfaces of standard libraries (for example,
io.Reader, io.WriterTo and others).

## Types

### type [ByteSlice](https://github.com/psyhatter/readonly/blob/main/byteslice.go#L20)

`type ByteSlice struct{ ... }`

ByteSlice wrapper over []byte that limits the interface to read-only.

#### func [NewByteSlice](https://github.com/psyhatter/readonly/blob/main/byteslice.go#L11)

`func NewByteSlice[T ~string | ~[]byte](src T) (b ByteSlice)`

NewByteSlice constructor for ByteSlice.
Accepts a string or slice of bytes as input, avoiding allocations.

#### func (ByteSlice) [Equal](https://github.com/psyhatter/readonly/blob/main/byteslice.go#L26)

`func (b ByteSlice) Equal(bb []byte) bool`

Equal s1 == s2 equivalent for string.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	b := readonly.NewByteSlice("some string")
	fmt.Println(b.Equal([]byte("some string")))
	fmt.Println(b.Equal([]byte("some other string")))
}

```

 Output:

```
true
false
```

#### func (ByteSlice) [ReadAt](https://github.com/psyhatter/readonly/blob/main/byteslice.go#L32)

`func (b ByteSlice) ReadAt(p []byte, off int64) (int, error)`

ReadAt implements io.ReaderAt.

#### func (ByteSlice) [String](https://github.com/psyhatter/readonly/blob/main/byteslice.go#L23)

`func (b ByteSlice) String() string`

String string(b) equivalent for byte slice.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	fmt.Println(readonly.NewByteSlice("some string").String())
}

```

 Output:

```
some string
```

#### func (ByteSlice) [WriteTo](https://github.com/psyhatter/readonly/blob/main/byteslice.go#L41)

`func (b ByteSlice) WriteTo(w io.Writer) (n int64, err error)`

WriteTo implements io.WriterTo.
w must not modify the slice data, even temporarily, see io.Writer.

### type [Chan](https://github.com/psyhatter/readonly/blob/main/chan.go#L8)

`type Chan[T any] <-chan T`

Chan wrapper over a built-in chan that limits the interface
to read-only.

#### func [NewChan](https://github.com/psyhatter/readonly/blob/main/chan.go#L4)

`func NewChan[T any](ch <-chan T) Chan[T]`

NewChan returns a chan interface limited to read-only methods.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	ch := make(chan int, 3)
	rch := readonly.NewChan(ch)

	// can't send to channel
	// rch <- 1
	ch <- 1
	ch <- 2
	ch <- 3

	// can't close the channel
	// close(rch)
	close(ch)

	// but will allow reading from the channel
	i1 := <-rch
	fmt.Println(i1)

	i2, ok := <-rch
	fmt.Println(i2, ok)

	for i3 := range rch {
		fmt.Println(i3)
	}

	fmt.Println(len(rch), cap(rch))
}

```

 Output:

```
1
2 true
3
0 3
```

### type [Map](https://github.com/psyhatter/readonly/blob/main/map.go#L8)

`type Map[k comparable, v any] struct { ... }`

Map wrapper over a built-in map that limits the interface
to read-only.

#### func [NewMap](https://github.com/psyhatter/readonly/blob/main/map.go#L4)

`func NewMap[k comparable, v any](m map[k]v) Map[k, v]`

NewMap returns a map interface limited to read-only methods.

#### func (Map[k, v]) [Get](https://github.com/psyhatter/readonly/blob/main/map.go#L17)

`func (m Map[k, v]) Get(key k) v`

Get equivalent to v := m[key].

#### func (Map[k, v]) [Get2](https://github.com/psyhatter/readonly/blob/main/map.go#L23)

`func (m Map[k, v]) Get2(key k) (v, bool)`

Get2 equivalent to v, ok := m[key].

#### func (Map[k, v]) [Has](https://github.com/psyhatter/readonly/blob/main/map.go#L20)

`func (m Map[k, v]) Has(key k) bool`

Has equivalent to _, ok := m[key].

#### func (Map[k, v]) [IsNil](https://github.com/psyhatter/readonly/blob/main/map.go#L11)

`func (m Map[k, v]) IsNil() bool`

IsNil equivalent to s == nil.

#### func (Map[k, v]) [Len](https://github.com/psyhatter/readonly/blob/main/map.go#L14)

`func (m Map[k, v]) Len() int`

Len equivalent to len(m).

#### func (Map[k, v]) [Range](https://github.com/psyhatter/readonly/blob/main/map.go#L28)

`func (m Map[k, v]) Range(f func(key k, val v) (next bool))`

Range equivalent to read-only for range loop.
Does nothing if f == nil.
Breaks the loop if next == false.

### type [Reader](https://github.com/psyhatter/readonly/blob/main/reader.go#L21)

`type Reader struct { ... }`

Reader implements the io.Reader, io.ByteReader, io.RuneReader and
io.WriterTo interfaces by reading from a string.
The zero value for Reader operates like a Reader of an empty string,
nil byte slice or an empty byte slice.

#### func [NewReader](https://github.com/psyhatter/readonly/blob/main/reader.go#L13)

`func NewReader[T ~string | ~[]byte | ByteSlice](src T) *Reader`

NewReader returns a new Reader reading from src.
It is similar to strings.NewReader or bytes.NewReader but more
efficient.

```golang
package main

import (
	"github.com/psyhatter/readonly"
)

func main() {
	readonly.NewReader("123")
	readonly.NewReader([]byte("123"))
	readonly.NewReader(readonly.NewByteSlice("123"))
}

```

#### func (*Reader) [Len](https://github.com/psyhatter/readonly/blob/main/reader.go#L25)

`func (r *Reader) Len() int`

Len returns the number of bytes of the unread portion of the
string.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	fmt.Println(readonly.NewReader("123").Len())
}

```

 Output:

```
3
```

#### func (*Reader) [Read](https://github.com/psyhatter/readonly/blob/main/reader.go#L28)

`func (r *Reader) Read(p []byte) (n int, err error)`

Read implements the io.Reader interface.

#### func (*Reader) [ReadByte](https://github.com/psyhatter/readonly/blob/main/reader.go#L38)

`func (r *Reader) ReadByte() (b byte, err error)`

ReadByte implements the io.ByteReader interface.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	b, err := readonly.NewReader("abc").ReadByte()
	fmt.Println(string(b), err)
}

```

 Output:

```
a <nil>
```

#### func (*Reader) [ReadRune](https://github.com/psyhatter/readonly/blob/main/reader.go#L47)

`func (r *Reader) ReadRune() (ch rune, size int, err error)`

ReadRune implements the io.RuneReader interface.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	r, size, err := readonly.NewReader("фыва").ReadRune()
	fmt.Println(string(r), size, err)
}

```

 Output:

```
ф 2 <nil>
```

#### func (*Reader) [ResetByteSlice](https://github.com/psyhatter/readonly/blob/main/reader.go#L88)

`func (r *Reader) ResetByteSlice(b ByteSlice)`

ResetByteSlice resets the Reader to be reading from b.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	r := readonly.NewReader("abc")
	fmt.Println(r.Len())
	r.ResetByteSlice(readonly.NewByteSlice("a"))
	fmt.Println(r.Len())
}

```

 Output:

```
3
1
```

#### func (*Reader) [ResetBytes](https://github.com/psyhatter/readonly/blob/main/reader.go#L85)

`func (r *Reader) ResetBytes(b []byte)`

ResetBytes resets the Reader to be reading from b.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	r := readonly.NewReader("abc")
	fmt.Println(r.Len())
	r.ResetBytes([]byte("a"))
	fmt.Println(r.Len())
}

```

 Output:

```
3
1
```

#### func (*Reader) [ResetString](https://github.com/psyhatter/readonly/blob/main/reader.go#L82)

`func (r *Reader) ResetString(s string)`

ResetString resets the Reader to be reading from s.

```golang
package main

import (
	"fmt"
	"github.com/psyhatter/readonly"
)

func main() {
	r := readonly.NewReader("abc")
	fmt.Println(r.Len())
	r.ResetString("a")
	fmt.Println(r.Len())
}

```

 Output:

```
3
1
```

#### func (Reader) [WriteTo](https://github.com/psyhatter/readonly/blob/main/reader.go#L62)

`func (r Reader) WriteTo(w io.Writer) (n int64, err error)`

WriteTo implements the io.WriterTo interface.
w must not modify the slice data, even temporarily, see io.Writer.

### type [Slice](https://github.com/psyhatter/readonly/blob/main/slice.go#L9)

`type Slice[T any] struct { ... }`

Slice wrapper over a built-in slice that limits the interface
to read-only.

#### func [NewSlice](https://github.com/psyhatter/readonly/blob/main/slice.go#L5)

`func NewSlice[T any](s []T) Slice[T]`

NewSlice returns a slice interface limited to read-only methods.

#### func (Slice[T]) [Append](https://github.com/psyhatter/readonly/blob/main/slice.go#L49)

`func (s Slice[T]) Append(dst []T) []T`

Append appends elements to the end of dst and returns the updated slice.

#### func (Slice[T]) [AppendInto](https://github.com/psyhatter/readonly/blob/main/slice.go#L54)

`func (s Slice[T]) AppendInto(to *[]T)`

AppendInto adds elements to the end of the slice located at the dst
pointer and places the new slice at the dst pointer.
Does nothing if dst == nil.

#### func (Slice[T]) [Cap](https://github.com/psyhatter/readonly/blob/main/slice.go#L18)

`func (s Slice[T]) Cap() int`

Cap equivalent to cap(s).

#### func (Slice[T]) [Copy](https://github.com/psyhatter/readonly/blob/main/slice.go#L41)

`func (s Slice[T]) Copy() []T`

Copy returns a new copy of the built-in slice.
As fast as inline copy to new slice, but faster than copying to a new
slice with (s Slice) CopyTo.

#### func (Slice[T]) [CopyTo](https://github.com/psyhatter/readonly/blob/main/slice.go#L46)

`func (s Slice[T]) CopyTo(dst []T) int`

CopyTo copies elements from a source slice into a destination slice.
The source and destination may overlap. Copy returns the number of
elements copied, which will be the minimum of (Slice) Len() and len(dst).

#### func (Slice[T]) [EndBefore](https://github.com/psyhatter/readonly/blob/main/slice.go#L64)

`func (s Slice[T]) EndBefore(i int) Slice[T]`

EndBefore equivalent to s[:i].

#### func (Slice[T]) [Get](https://github.com/psyhatter/readonly/blob/main/slice.go#L21)

`func (s Slice[T]) Get(index int) (v T)`

Get equivalent to v := s[index].

#### func (Slice[T]) [IsNil](https://github.com/psyhatter/readonly/blob/main/slice.go#L12)

`func (s Slice[T]) IsNil() bool`

IsNil equivalent to s == nil.

#### func (Slice[T]) [Len](https://github.com/psyhatter/readonly/blob/main/slice.go#L15)

`func (s Slice[T]) Len() int`

Len equivalent to len(s).

#### func (Slice[T]) [Range](https://github.com/psyhatter/readonly/blob/main/slice.go#L28)

`func (s Slice[T]) Range(f func(index int, val T) (next bool))`

Range equivalent to read-only for range loop.
Does nothing if f == nil.
Breaks the loop if next == false.
An order of magnitude slower than the built-in slice loop, for
optimizations, you can use slice index access (see benchmarks).

#### func (Slice[T]) [Slice](https://github.com/psyhatter/readonly/blob/main/slice.go#L67)

`func (s Slice[T]) Slice(start, end int) Slice[T]`

Slice equivalent to s[start:end].

#### func (Slice[T]) [Slice3](https://github.com/psyhatter/readonly/blob/main/slice.go#L70)

`func (s Slice[T]) Slice3(i, j, k int) Slice[T]`

Slice3 equivalent to s[i:j:k].

#### func (Slice[T]) [StartAfter](https://github.com/psyhatter/readonly/blob/main/slice.go#L61)

`func (s Slice[T]) StartAfter(i int) Slice[T]`

StartAfter equivalent to s[i:].

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
