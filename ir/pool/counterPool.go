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

// CounterPool represents a pool of counters with exclusions.
type CounterPool struct {
	Ref        int          // Ref is the current reference counter.
	Exclusions map[int]bool // Exclusions is a map of excluded references.
}

// RequestSuitable finds and returns the next suitable reference
// that is not in the exclusions map.
func (pool *CounterPool) RequestSuitable() int {
	// Find a suitable reference
	for pool.Exclusions[pool.Ref] {
		pool.Ref++
	}

	// Store the current ref and return the next suitable value
	returnVal := pool.Ref
	pool.Ref++
	return returnVal
}
