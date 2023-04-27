package asttransform

import (
	"fmt"
	"go/ast"
)

// TransformFunc is a function that takes an ast.Node as input,
// applies a transformation, and returns the modified ast.Node or an error.
type TransformFunc func(node ast.Node) (ast.Node, error)

// ApplyTransforms applies a series of TransformFuncs to an ast.Node,
// returning the modified ast.Node or an error if any of the transformations fail.
func ApplyTransforms(node ast.Node, transforms []TransformFunc) (ast.Node, error) {
	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}

	var err error
	for _, transform := range transforms {
		node, err = transform(node)
		if err != nil {
			return nil, fmt.Errorf("error applying transform: %w", err)
		}
	}

	return node, nil
}
