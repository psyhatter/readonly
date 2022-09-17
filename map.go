package readonly

// NewMap returns a map interface limited to read-only methods.
func NewMap[k comparable, v any](m map[k]v) Map[k, v] { return Map[k, v]{m: m} }

// Map wrapper over a built-in map that limits the interface
// to read-only.
type Map[k comparable, v any] struct{ m map[k]v }

// Len equivalent to len(m).
func (m Map[k, v]) Len() int { return len(m.m) }

// Get equivalent to v := m[key].
func (m Map[k, v]) Get(key k) v { return m.m[key] }

// Has equivalent to _, ok := m[key].
func (m Map[k, v]) Has(key k) bool { _, ok := m.m[key]; return ok }

// Get2 equivalent to v, ok := m[key].
func (m Map[k, v]) Get2(key k) (v, bool) { val, ok := m.m[key]; return val, ok }

// Range equivalent to read-only for range loop.
// Does nothing if f == nil.
// Breaks the loop if next == false.
func (m Map[k, v]) Range(f func(key k, val v) (next bool)) {
	if f != nil {
		for key, val := range m.m {
			if !f(key, val) {
				return
			}
		}
	}
}
