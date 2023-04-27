package asttransform

import (
	"go/ast"
)

// TransformFunc is a function that takes an ast.Node as input,
// applies a transformation, and returns the modified ast.Node or an error.
type TransformFunc func(node ast.Node) (ast.Node, error)

// ApplyTransforms applies a series of TransformFuncs to an ast.Node,
// returning the modified ast.Node or an error if any of the transformations fail.
func ApplyTransforms(node ast.Node, transforms []TransformFunc) (ast.Node, error) {
	// ...
	return nil, nil
}
