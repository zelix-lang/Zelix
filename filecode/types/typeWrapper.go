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
	"strings"
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
func (t *TypeWrapper) Compare(other TypeWrapper) bool {
	result := true

	// Use a queue to compare the children
	queue := make([]pairQueueElement, 0)
	queue = append(queue, pairQueueElement{
		w:     t,
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

		if element.w.Children == nil && element.other.Children != nil {
			result = false
			break
		}

		if element.w.Children != nil && element.other.Children == nil {
			result = false
			break
		}

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

// Marshal serializes the TypeWrapper into a string representation.
// It constructs the result using a string builder and processes nodes iteratively.
//
// Returns:
//   - string: The serialized string representation of the TypeWrapper.
func (t *TypeWrapper) Marshal() string {
	// Use a string builder to construct the result
	var builder strings.Builder

	// Define a queue to process nodes iteratively
	type queueItem struct {
		node  *TypeWrapper
		depth int
	}

	queue := []queueItem{{node: t, depth: 0}}

	// Stack to store output fragments
	var outputStack []string

	for len(queue) > 0 {
		// Dequeue
		item := queue[0]
		queue = queue[1:]

		node := item.node

		// Process pointer symbols
		prefix := strings.Repeat("&", node.PointerCount)

		// Process base type
		content := prefix + node.BaseType

		// Handle children
		if node.Children != nil && len(*node.Children) > 0 {
			childStrings := make([]string, len(*node.Children))
			for i, child := range *node.Children {
				queue = append(queue, queueItem{node: child, depth: item.depth + 1})
				childStrings[i] = child.BaseType
			}

			content += "<" + strings.Join(childStrings, ", ") + ">"
		}

		// Process array count
		content += strings.Repeat("[]", node.ArrayCount)

		// Store the processed part
		outputStack = append(outputStack, content)
	}

	// Construct the final result
	for _, part := range outputStack {
		builder.WriteString(part)
	}

	return builder.String()
}
