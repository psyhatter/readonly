// Package readonly provides an interface to some built-in
// container types (such as slices, maps, and channels) that
// allow them to be read-only. The package uses generics so
// that they can be used for any types. Also, byte sequences
// (strings or slices of bytes) are additionally processed
// through the interfaces of standard libraries (for example,
// io.Reader, io.WriterTo and others).
package readonly
