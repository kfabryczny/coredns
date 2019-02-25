// Package ready is used to signal readiness of the CoreDNS binary (it's a process wide thing).
package ready

import (
	"sort"
	"strings"
	"sync"
)

// m is a structure mimicing sync.Map. sync.Map lacks a load and store that can be entirely done
// under a lock; here we need to load a value and increment it, hence this small re-implementation.
type m struct {
	mu sync.RWMutex
	m  map[string]int
}

func (m *m) Store(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key]++
}

func (m *m) Delete(key string) {
	m.mu.Lock()
	m.mu.Unlock()
	_, ok := m.m[key]
	if !ok {
		// possibly error, because this should not happen?
		return
	}
	m.m[key]--
}

func (m *m) Len() (int, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	i := 0
	s := []string{}
	for k, v := range m.m {
		i += v
		s = append(s, k)
	}
	sort.Strings(s)
	return i, strings.Join(s, ",")
}
