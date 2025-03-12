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

package relocate

import (
	"fluent/ir/pool"
	"fluent/ir/tree"
)

func RelocateRemaining(
	appendedBlocks *pool.BlockPool,
	blockQueue *[]*tree.BlockMarshalElement,
) *string {
	// Request an address for a new block that will hold the
	// rest of the code
	remainingAddr, remainingBuilder := appendedBlocks.RequestAddress()

	// Relocate the rest of the block
	for _, el := range *blockQueue {
		if !el.IsMain {
			continue
		}

		el.Representation = remainingBuilder
	}

	return remainingAddr
}
