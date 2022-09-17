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

func ExampleByteSlice_String() {
	fmt.Println(readonly.NewByteSlice("some string").String())
	// Output:
	// some string
}

func ExampleByteSlice_Equal() {
	b := readonly.NewByteSlice("some string")
	fmt.Println(b.Equal([]byte("some string")))
	fmt.Println(b.Equal([]byte("some other string")))
	// Output:
	// true
	// false
}

func TestByteSlice_ReadAt(t *testing.T) {
	readAll := func(r readonly.ByteSlice) ([]byte, error) {
		return io.ReadAll(io.NewSectionReader(r, 0, int64(r.Len())))
	}
	for i, expected := range [][]byte{
		[]byte("short slice"),
		append([]byte("big slice"), make([]byte, 1<<15)...),
	} {
		actual, err := readAll(readonly.NewByteSlice(expected))
		if err != nil {
			t.Fatalf("[%d] unexpected err: %v", i, err)
		}
		if !bytes.Equal(expected, actual) {
			t.Fatalf("expected %q, got %q", expected, actual)
		}
	}

	r := readonly.NewByteSlice("123")
	_, err := r.ReadAt(nil, int64(r.Len()))
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected %q, got %q", io.EOF, err)
	}
}

type writer func([]byte) (int, error)

func (w writer) Write(p []byte) (n int, err error) { return w(p) }

func TestByteSlice_WriteTo(t *testing.T) {
	var (
		buf      bytes.Buffer
		expected = []byte("some text")
		r        = readonly.NewByteSlice(expected)
	)

	_, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("expected %q, got %q", expected, buf.Bytes())
	}

	var mustNotBeCalled writer = func(p []byte) (int, error) {
		return 0, fmt.Errorf("unexpected function call with slice: %s", p)
	}

	for i, f := range []func(io.Writer) (int64, error){
		readonly.ByteSlice{}.WriteTo,
		readonly.NewByteSlice("").WriteTo,
		readonly.NewByteSlice([]byte{}).WriteTo,
		readonly.NewByteSlice([]byte(nil)).WriteTo,
		readonly.NewByteSlice[[]byte](nil).WriteTo,
	} {
		if _, err = f(mustNotBeCalled); err != nil {
			t.Fatalf("[%d]: %v", i, err)
		}
	}

	var dontWrite, returnsErr writer = func([]byte) (int, error) { return 0, nil },
		func([]byte) (int, error) { return 0, errors.New("some error") }

	_, err = r.WriteTo(dontWrite)
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("expected %q, got %q", io.ErrShortWrite, err)
	}
	_, err = r.WriteTo(returnsErr)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
