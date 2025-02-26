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

type StringPool struct {
	Storage map[string]string
	Counter int
}

func (pool *StringPool) RequestAddress(str string) string {
	// Check if the string has been saved previously
	address, ok := pool.Storage[str]

	if ok {
		return address
	}

	// Create a new address for this string
	address = fmt.Sprintf("__str_x%d", pool.Counter)
	pool.Storage[str] = address
	pool.Counter++

	return address
}
