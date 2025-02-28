/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package util

import (
	"container/list"
)

// OrderedMap is a map that maintains the insertion order of keys.
type OrderedMap[K comparable, V any] struct {
	data map[K]*list.Element
	list *list.List
}

// KeyValue represents a key-value pair.
type KeyValue[K comparable, V any] struct {
	key   K
	value V
}

// NewOrderedMap creates a new OrderedMap.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		data: make(map[K]*list.Element),
		list: list.New(),
	}
}

// Set inserts or updates the value associated with the key.
func (om *OrderedMap[K, V]) Set(key K, value V) {
	if elem, exists := om.data[key]; exists {
		// Update existing value
		elem.Value.(*KeyValue[K, V]).value = value
	} else {
		// Insert new key-value pair
		kv := &KeyValue[K, V]{key, value}
		om.data[key] = om.list.PushBack(kv)
	}
}

// Get retrieves the value associated with the key.
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	if elem, exists := om.data[key]; exists {
		return elem.Value.(*KeyValue[K, V]).value, true
	}
	var zero V
	return zero, false
}

// Delete removes the key-value pair associated with the key.
func (om *OrderedMap[K, V]) Delete(key K) {
	if elem, exists := om.data[key]; exists {
		om.list.Remove(elem)
		delete(om.data, key)
	}
}

// Iterate iterates over all key-value pairs in insertion order.
func (om *OrderedMap[K, V]) Iterate(predicate func(K, V)) {
	for elem := om.list.Front(); elem != nil; elem = elem.Next() {
		kv := elem.Value.(*KeyValue[K, V])
		predicate(kv.key, kv.value)
	}
}
