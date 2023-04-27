package asttransform

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
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

// ParseGoFile reads a Go source file at the given path and returns its AST.
func ParseGoFile(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Print takes an *ast.File and pretty-prints it as a Go source code string.
// It returns an error if the pretty-printing process fails.
func Print(file *ast.File) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}

	var buf bytes.Buffer
	err := format.Node(&buf, token.NewFileSet(), file)
	if err != nil {
		return "", fmt.Errorf("error formatting AST: %w", err)
	}

	return buf.String(), nil
}
