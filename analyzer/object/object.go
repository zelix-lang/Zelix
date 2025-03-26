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

package object

import (
	"fluent/filecode/types/wrapper"
)

// Object represents a generic object in the Fluent programming language.
// It contains information about whether the object is allocated on the heap,
// the value of the object, and the type of the object.
type Object struct {
	// IsHeap indicates if the object is allocated on the heap.
	IsHeap bool
	// Value holds the actual value of the object.
	Value interface{}
	// Type represents the type information of the object.
	Type wrapper.TypeWrapper
}
