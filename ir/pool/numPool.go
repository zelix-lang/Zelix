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

package pool

import "fmt"

// NumPool is a structure that holds a pool of numbers and their associated counters.
type NumPool struct {
	Storage map[int]string // Storage maps numbers to their addresses.
	Counter map[int]int    // Counter keeps track of the number of strings associated with each id.
}

// NumPool adds a new id to the pool if it does not already exist.
func (pool *StringPool) NumPool(id int) {
	// Check if the id has been added previously
	_, ok := pool.Counter[id]

	if !ok {
		pool.Counter[id] = 0
	}
}

// RemoveId removes an id from the pool if it exists.
func (pool *NumPool) RemoveId(id int) {
	// Check if the id has been added previously
	_, ok := pool.Counter[id]

	if ok {
		delete(pool.Counter, id)
	}
}

// RequestAddress returns the address of a number associated with a given id.
// If the string is not already in the pool, it creates a new address for it.
func (pool *NumPool) RequestAddress(id int, val int) string {
	// Check if the string has been saved previously
	address, ok := pool.Storage[val]

	if ok {
		return address
	}

	// Get the counter for the specified id
	counter := pool.Counter[id]

	// Create a new address for this string
	address = fmt.Sprintf("__trace_magic__f%d_x%d", id, counter)
	pool.Storage[val] = address
	pool.Counter[id]++

	return address
}
