//nolint:exhaustivestruct
package readonly_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/psyhatter/readonly"
)

func ExampleNewReader() {
	readonly.NewReader("123")
	readonly.NewReader([]byte("123"))
	readonly.NewReader(readonly.NewByteSlice("123"))
}

func ExampleReader_Len() {
	fmt.Println(readonly.NewReader("123").Len())
	// Output:
	// 3
}

func ExampleReader_ReadByte() {
	b, err := readonly.NewReader("abc").ReadByte()
	fmt.Println(string(b), err)
	// Output:
	// a <nil>
}

func ExampleReader_ReadRune() {
	r, size, err := readonly.NewReader("фыва").ReadRune()
	fmt.Println(string(r), size, err)
	// Output:
	// ф 2 <nil>
}

func ExampleReader_ResetString() {
	r := readonly.NewReader("abc")
	fmt.Println(r.Len())
	r.ResetString("a")
	fmt.Println(r.Len())
	// Output:
	// 3
	// 1
}

func ExampleReader_ResetBytes() {
	r := readonly.NewReader("abc")
	fmt.Println(r.Len())
	r.ResetBytes([]byte("a"))
	fmt.Println(r.Len())
	// Output:
	// 3
	// 1
}

func ExampleReader_ResetByteSlice() {
	r := readonly.NewReader("abc")
	fmt.Println(r.Len())
	r.ResetByteSlice(readonly.NewByteSlice("a"))
	fmt.Println(r.Len())
	// Output:
	// 3
	// 1
}

func TestReader_ReadByte(t *testing.T) {
	var (
		expected = []byte("123")
		actual   = make([]byte, 0, len(expected))
		r        = readonly.NewReader(expected)
	)

	for b, err := r.ReadByte(); !errors.Is(err, io.EOF); b, err = r.ReadByte() {
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		actual = append(actual, b)
	}

	if !bytes.Equal(expected, actual) {
		t.Fatalf("expected %q, got %q", expected, actual)
	}

	for i, f := range []func() (byte, error){
		(&readonly.Reader{}).ReadByte,
		readonly.NewReader("").ReadByte,
		readonly.NewReader([]byte{}).ReadByte,
		readonly.NewReader([]byte(nil)).ReadByte,
		readonly.NewReader[[]byte](nil).ReadByte,
		readonly.NewReader(readonly.ByteSlice{}).ReadByte,
	} {
		if _, err := f(); !errors.Is(err, io.EOF) {
			t.Fatalf("[%d] expected %q, got %q", i, io.EOF, err)
		}
	}
}

func TestReader_ReadRune(t *testing.T) {
	var (
		expected = []rune("ABCэюя")
		actual   = make([]rune, 0, len(expected))
		r        = readonly.NewReader(string(expected))
	)

	for c, _, err := r.ReadRune(); !errors.Is(err, io.EOF); c, _, err = r.ReadRune() {
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		actual = append(actual, c)
	}

	if string(expected) != string(actual) {
		t.Fatalf("expected %q, got %q", expected, actual)
	}

	for i, f := range []func() (rune, int, error){
		(&readonly.Reader{}).ReadRune,
		readonly.NewReader("").ReadRune,
		readonly.NewReader([]byte{}).ReadRune,
		readonly.NewReader([]byte(nil)).ReadRune,
		readonly.NewReader[[]byte](nil).ReadRune,
		readonly.NewReader(readonly.ByteSlice{}).ReadRune,
	} {
		if _, _, err := f(); !errors.Is(err, io.EOF) {
			t.Fatalf("[%d] expected %q, got %q", i, io.EOF, err)
		}
	}
}

func TestReader_WriteTo(t *testing.T) {
	var (
		buf      bytes.Buffer
		expected = []byte("some text")
		r        = readonly.NewReader(expected)
	)

	_, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("expected %q, got %q", expected, buf.Bytes())
	}

	var mustNotBeCalled writer = func(p []byte) (int, error) {
		return 0, fmt.Errorf("unexpected function call with slice: %q", p)
	}

	for i, f := range []func(io.Writer) (int64, error){
		readonly.Reader{}.WriteTo,
		readonly.NewReader("").WriteTo,
		readonly.NewReader([]byte{}).WriteTo,
		readonly.NewReader([]byte(nil)).WriteTo,
		readonly.NewReader[[]byte](nil).WriteTo,
		readonly.NewReader(readonly.ByteSlice{}).WriteTo,
	} {
		if _, err = f(mustNotBeCalled); err != nil {
			t.Fatalf("[%d]: %v", i, err)
		}
	}

	var dontWrite, returnsErr writer = func(p []byte) (int, error) { return 0, nil },
		func(p []byte) (int, error) { return 0, errors.New("some error") }

	_, err = r.WriteTo(dontWrite)
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("expected %q, got %q", io.ErrShortWrite, err)
	}
	_, err = r.WriteTo(returnsErr)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestReader_Read(t *testing.T) {
	var (
		expected = []byte("123")
		actual   = make([]byte, len(expected))
		r        = readonly.NewReader(expected)
	)

	_, err := r.Read(actual)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if !bytes.Equal(expected, actual) {
		t.Fatalf("expected %q, got %q", expected, actual)
	}

	for i, r := range []*readonly.Reader{
		{},
		readonly.NewReader(""),
		readonly.NewReader([]byte{}),
		readonly.NewReader([]byte(nil)),
		readonly.NewReader[[]byte](nil),
		readonly.NewReader(readonly.ByteSlice{}),
	} {
		if _, err := r.Read(nil); !errors.Is(err, io.EOF) {
			t.Fatalf("[%d] expected %q, got %q", i, io.ErrShortWrite, err)
		}
	}
}

func TestReader_Reset(t *testing.T) {
	var r readonly.Reader

	first := "123"
	r.ResetString(first)
	b, err := io.ReadAll(&r)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if first != string(b) {
		t.Fatalf("expected %q, got %q", first, b)
	}

	second := []byte("456")
	r.ResetBytes(second)
	b, err = io.ReadAll(&r)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if !bytes.Equal(second, b) {
		t.Fatalf("expected %q, got %q", first, b)
	}

	third := readonly.NewByteSlice("789")
	r.ResetByteSlice(third)
	b, err = io.ReadAll(&r)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if third.String() != string(b) {
		t.Fatalf("expected %q, got %q", first, b)
	}
}
