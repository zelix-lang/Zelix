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

// StringPool is a structure that holds a pool of strings and their associated counters.
type StringPool struct {
	Storage map[string]string // Storage maps strings to their addresses.
	Counter map[int]int       // Counter keeps track of the number of strings associated with each id.
	Prefix  string            // Prefix is a prefix that is appended to all memory spaces.
}

// AddNewId adds a new id to the pool if it does not already exist.
func (pool *StringPool) AddNewId(id int) {
	// Check if the id has been added previously
	_, ok := pool.Counter[id]

	if !ok {
		pool.Counter[id] = 0
	}
}

// RemoveId removes an id from the pool if it exists.
func (pool *StringPool) RemoveId(id int) {
	// Check if the id has been added previously
	_, ok := pool.Counter[id]

	if ok {
		delete(pool.Counter, id)
	}
}

// RequestAddress returns the address of a string associated with a given id.
// If the string is not already in the pool, it creates a new address for it.
func (pool *StringPool) RequestAddress(id int, str string) string {
	// Check if the string has been saved previously
	address, ok := pool.Storage[str]

	if ok {
		return address
	}

	// Get the counter for the specified id
	counter := pool.Counter[id]

	// Create a new address for this string
	address = fmt.Sprintf("%sf%d_x%d", pool.Prefix, id, counter)
	pool.Storage[str] = address
	pool.Counter[id]++

	return address
}
