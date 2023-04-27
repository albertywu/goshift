package asttransform_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/albertywu/goshift/pkg/asttransform"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestParseGoFile(t *testing.T) {
	// Create a temporary Go file for testing
	tempDir, err := os.MkdirTemp("", "asttransform_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testGoFile := filepath.Join(tempDir, "test.go")
	err = ioutil.WriteFile(testGoFile, []byte("package main\n\nfunc main() {}\n"), 0644)
	require.NoError(t, err)

	// Attempt to parse the temporary Go file
	parsedFile, err := asttransform.ParseGoFile(testGoFile)
	require.NoError(t, err)
	require.NotNil(t, parsedFile)

	// Check if the parsed file has the expected package name and function name
	require.Equal(t, "main", parsedFile.Name.Name)
	require.Len(t, parsedFile.Decls, 1)

	funcDecl, ok := parsedFile.Decls[0].(*ast.FuncDecl)
	require.True(t, ok)
	require.Equal(t, "main", funcDecl.Name.Name)
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name    string
		srcCode string
	}{
		{
			name: "simple function",
			srcCode: `package main

func main() {
	println("Hello, world!")
}
`,
		},
		{
			name: "multiple functions",
			srcCode: `package main

func foo() {
	println("Foo")
}

func bar() {
	println("Bar")
}

func main() {
	foo()
	bar()
}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the source code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tc.srcCode, parser.ParseComments)
			require.NoError(t, err)

			// Print the AST
			printedSrc, err := asttransform.Print(file)
			require.NoError(t, err)

			// Re-parse the printed source code
			reParsedFile, err := parser.ParseFile(fset, "", printedSrc, parser.ParseComments)
			require.NoError(t, err)

			// Check if the original and re-parsed ASTs are equal
			assert.Empty(t, cmp.Diff(file, reParsedFile,
				cmpopts.IgnoreFields(ast.Ident{}, "NamePos"),
				cmpopts.IgnoreFields(ast.BasicLit{}, "ValuePos"),
				cmpopts.IgnoreFields(ast.FuncType{}, "Func"),
				cmpopts.IgnoreFields(ast.FieldList{}, "Opening", "Closing"),
				cmpopts.IgnoreFields(ast.BlockStmt{}, "Lbrace", "Rbrace"),
				cmpopts.IgnoreFields(ast.CallExpr{}, "Lparen", "Rparen"),
				cmpopts.IgnoreFields(ast.File{}, "Package"),
			), "Original and reparsed ASTs should be equal")
		})
	}
}
