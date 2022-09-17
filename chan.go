package readonly

// NewChan returns a chan interface limited to read-only methods.
func NewChan[T any](ch <-chan T) Chan[T] { return ch }

// Chan wrapper over a built-in chan that limits the interface
// to read-only.
type Chan[T any] <-chan T
