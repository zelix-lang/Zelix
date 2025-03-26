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

import (
	"fmt"
	"strings"
)

// BlockPool represents a pool of blocks with a storage map and a counter.
type BlockPool struct {
	Storage map[string]*strings.Builder // Storage holds the block data.
	Counter int                         // Counter keeps track of the number of blocks.
}

// RequestAddress creates a new block address and builder, stores them in the pool, and returns them.
func (pool *BlockPool) RequestAddress() (*string, *strings.Builder) {
	// Create a new builder
	builder := strings.Builder{}

	// Create a new address
	address := fmt.Sprintf("__block_x%d__", pool.Counter)

	// Write the changes
	pool.Storage[address] = &builder
	pool.Counter++

	return &address, &builder
}
