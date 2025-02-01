/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package types

import (
	"fluent/filecode/trace"
)

// TypeWrapper represents a type with additional metadata.
type TypeWrapper struct {
	PointerCount int             // Number of pointers to the base type.
	ArrayCount   int             // Number of arrays of the base type.
	Children     *[]*TypeWrapper // Nested types within this type.
	BaseType     string          // The base type as a string.
	Trace        trace.Trace     // Trace information for the type.
	IsPrimitive  bool            // Indicates whether the type is a primitive type.
}

type pairQueueElement struct {
	w     *TypeWrapper
	other *TypeWrapper
}

// Compare compares the current TypeWrapper with another TypeWrapper.
// It returns true if they are equivalent, otherwise false.
//
// Parameters:
//   - other: The other TypeWrapper to compare with.
//
// Returns:
//   - bool: True if the TypeWrappers are equivalent, otherwise false.
func (w *TypeWrapper) Compare(other TypeWrapper) bool {
	result := true

	// Use a queue to compare the children
	queue := make([]pairQueueElement, 0)
	queue = append(queue, pairQueueElement{
		w:     w,
		other: &other,
	})

	for len(queue) > 0 {
		// Pop the first element from the queue
		element := queue[0]
		queue = queue[1:]

		// Compare the base type
		if element.w.BaseType != element.other.BaseType {
			result = false
			break
		}

		// Compare the pointer count
		if element.w.PointerCount != element.other.PointerCount {
			result = false
			break
		}

		// Compare the array count
		if element.w.ArrayCount != element.other.ArrayCount {
			result = false
			break
		}

		// Compare the number of children
		if len(*element.w.Children) != len(*element.other.Children) {
			result = false
			break
		}

		// Add the children to the queue
		for i := 0; i < len(*element.w.Children); i++ {
			queue = append(queue, pairQueueElement{
				w:     (*element.w.Children)[i],
				other: (*element.other.Children)[i],
			})
		}
	}

	return result
}
