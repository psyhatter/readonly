package readonly_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"strings"
	"testing"

	"github.com/psyhatter/readonly"
)

//nolint:forcetypeassert
var rf = io.Discard.(io.ReaderFrom)

//nolint:funlen
func BenchmarkReaders(b *testing.B) {
	runBench := func(b *testing.B, size int) {
		data := make([]byte, size)
		_, _ = rand.Read(data)
		b.Run("[]byte", func(b *testing.B) {
			b.Run("readonly.Reader", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := readonly.NewReader(data)
					_, _ = rf.ReadFrom(r)
				}
			})
			b.Run("bytes.Reader", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := bytes.NewReader(data)
					_, _ = rf.ReadFrom(r)
				}
			})
			b.Run("strings.Reader", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := strings.NewReader(string(data))
					_, _ = rf.ReadFrom(r)
				}
			})
			b.Run("io.SectionReader+readonly.ByteSlice", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := io.NewSectionReader(readonly.NewByteSlice(data), 0, int64(len(data)))
					_, _ = rf.ReadFrom(r)
				}
			})
		})
		b.Run("string", func(b *testing.B) {
			data := string(data)
			b.ResetTimer()
			b.Run("readonly.Reader", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := readonly.NewReader(data)
					_, _ = rf.ReadFrom(r)
				}
			})
			b.Run("strings.Reader", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := strings.NewReader(data)
					_, _ = rf.ReadFrom(r)
				}
			})
			b.Run("bytes.Reader", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := bytes.NewReader([]byte(data))
					_, _ = rf.ReadFrom(r)
				}
			})
			b.Run("io.SectionReader+readonly.ByteSlice", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					r := io.NewSectionReader(readonly.NewByteSlice(data), 0, int64(len(data)))
					_, _ = rf.ReadFrom(r)
				}
			})
		})
	}

	b.Run("small", func(b *testing.B) { runBench(b, 1<<8) })
	b.Run("large", func(b *testing.B) { runBench(b, 1<<16) })
}

func BenchmarkWriterTo(b *testing.B) {
	w, data := io.Discard, make([]byte, limit)
	_, _ = rand.Read(data)
	b.Run("[]byte", func(b *testing.B) {
		b.Run("readonly.Reader", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := readonly.NewReader(data)
				_, _ = r.WriteTo(w)
			}
		})
		b.Run("readonly.ByteSlice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := readonly.NewByteSlice(data)
				_, _ = r.WriteTo(w)
			}
		})
		b.Run("bytes.Reader", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader(data)
				_, _ = r.WriteTo(w)
			}
		})
		b.Run("strings.Reader", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := strings.NewReader(string(data))
				_, _ = r.WriteTo(w)
			}
		})
	})
	b.Run("string", func(b *testing.B) {
		data := string(data)
		b.ResetTimer()
		b.Run("readonly.Reader", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := readonly.NewReader(data)
				_, _ = r.WriteTo(w)
			}
		})
		b.Run("readonly.ByteSlice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := readonly.NewByteSlice(data)
				_, _ = r.WriteTo(w)
			}
		})
		b.Run("strings.Reader", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := strings.NewReader(data)
				_, _ = r.WriteTo(w)
			}
		})
		b.Run("bytes.Reader", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader([]byte(data))
				_, _ = r.WriteTo(w)
			}
		})
	})
}
