package asttransform_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/albertywu/goshift/pkg/asttransform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sampleTransform is a simple transform function that renames all functions
// with the name "oldName" to "newName".
func sampleTransform(oldName, newName string) asttransform.TransformFunc {
	return func(node ast.Node) (ast.Node, error) {
		ast.Inspect(node, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if ok && fn.Name.Name == oldName {
				fn.Name.Name = newName
			}
			return true
		})
		return node, nil
	}
}

func TestApplyTransforms(t *testing.T) {
	src := `
package main

func oldName() {
	println("Hello, world!")
}

func main() {
	oldName()
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	require.NoError(t, err)

	transforms := []asttransform.TransformFunc{
		sampleTransform("oldName", "newName"),
	}

	modifiedNode, err := asttransform.ApplyTransforms(file, transforms)
	require.NoError(t, err)

	modifiedFile, ok := modifiedNode.(*ast.File)
	require.True(t, ok, "Expected modifiedNode to be of type *ast.File")

	// Check if the function was renamed
	found := false
	for _, decl := range modifiedFile.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == "newName" {
			found = true
			break
		}
	}

	assert.True(t, found, "Expected to find a function with the name 'newName'")
}
