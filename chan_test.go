package readonly_test

import (
	"fmt"
	"testing"

	"github.com/psyhatter/readonly"
)

func ExampleNewChan() {
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
	// Output:
	// 1
	// 2 true
	// 3
	// 0 3
}

func BenchmarkChan(b *testing.B) {
	getChan := func() chan int {
		ch := make(chan int, 100)
		for i := 0; i < cap(ch)-1; i++ {
			ch <- i
		}
		close(ch)
		return ch
	}
	b.Run("built-in", func(b *testing.B) {
		ch := getChan()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			<-ch
		}
	})
	b.Run("readonly", func(b *testing.B) {
		ch := readonly.NewChan(getChan())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			<-ch
		}
	})
}
