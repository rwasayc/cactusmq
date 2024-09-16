package base

import (
	"sync"
	"sync/atomic"
)

// NewSyncMap creates a new concurrent map.
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{internal: &sync.Map{}}
}

// SyncMap is a concurrent map.
type SyncMap[K comparable, V any] struct {
	internal *sync.Map
	size     atomic.Uint32
}

// Store stores a value for a key.
func (m *SyncMap[K, V]) Store(key K, value V) {
	pre, _ := m.internal.Swap(key, value)
	if pre == nil {
		m.size.Add(1)
	}
}

// Load loads a value for a key.
func (m *SyncMap[K, V]) Load(key K) (value V, exist bool) {
	var v any
	v, exist = m.internal.Load(key)
	if !exist {
		return
	}
	value = v.(V)
	return
}

// Delete deletes a value for a key.
func (m *SyncMap[K, V]) Delete(key K) {
	pre, _ := m.internal.LoadAndDelete(key)
	if pre != nil {
		m.size.Add(^uint32(0))
	}
}

// Range ranges over the map.
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.internal.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// Len returns the size of the map.
func (m *SyncMap[K, V]) Len() int {
	return int(m.size.Load())
}
