package concurrent

import (
	"sync"
)

// TypedConcurrentSlice is a thread-safe implementation of a slice
type TypedConcurrentSlice[T any] struct {
	mu    sync.RWMutex
	slice []T
}

// NewTypedConcurrentSlice initializes and returns a new TypedConcurrentSlice
func NewTypedConcurrentSlice[T any]() *TypedConcurrentSlice[T] {
	return &TypedConcurrentSlice[T]{}
}

// Append adds an element to the slice
func (cs *TypedConcurrentSlice[T]) Append(value *T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.slice = append(cs.slice, *value)
}

// Get retrieves the element at a given index
func (cs *TypedConcurrentSlice[T]) Get(index int) (T, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	if index < 0 || index >= len(cs.slice) {
		var zeroValue T
		return zeroValue, false
	}
	return cs.slice[index], true
}

// Remove removes the element at a given index
func (cs *TypedConcurrentSlice[T]) Remove(index int) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if index < 0 || index >= len(cs.slice) {
		return false
	}
	cs.slice = append(cs.slice[:index], cs.slice[index+1:]...)
	return true
}

// Length returns the current length of the slice
func (cs *TypedConcurrentSlice[T]) Length() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return len(cs.slice)
}

// Snapshot returns a copy of the current slice
func (cs *TypedConcurrentSlice[T]) Snapshot() []T {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	copySlice := make([]T, len(cs.slice))
	copy(copySlice, cs.slice)
	return copySlice
}
