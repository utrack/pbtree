package vmap

// Map is just a versioned map :) Any change to it
// increments its version by 1.
//
// Not thread-safe.
type Map struct {
	m map[string]string
	v int
}

func New(m map[string]string) *Map {
	if m == nil {
		m = map[string]string{}
	}
	return &Map{m: m}
}

func (m Map) Get(k string) (string, bool) {
	ret, ok := m.m[k]
	return ret, ok
}

func (m *Map) Put(k, v string) {
	m.m[k] = v
	m.v++
}

func (m *Map) Version() int {
	return m.v
}

func (m *Map) Values() map[string]string {
	return m.m
}
