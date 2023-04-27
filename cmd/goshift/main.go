package main

import (
	"flag"
	"fmt"
	"github.com/albertywu/goshift/pkg/asttransform"
	"github.com/albertywu/goshift/pkg/fileutil"
	"github.com/albertywu/goshift/pkg/patch"
	"go/ast"
	"go/token"
	"log"
	"os"
)

var (
	root       string
	pattern    string
	output     string
	transforms []asttransform.TransformFunc
)

func init() {
	flag.StringVar(&root, "root", ".", "Root directory to apply transformations")
	flag.StringVar(&pattern, "pattern", "*.go", "File matching pattern for Go source files")
	flag.StringVar(&output, "output", "output.patch", "Output file for the patch")
}

func main() {
	flag.Parse()

	// Add AST transformations here
	transforms = append(
		transforms,
		renameFunc("oldName", "newName"),
		renameFunc("foo", "bar"),
		renameFunc("one", "two"),
	)

	// Create a patch to store file diffs
	p := patch.Patch{}

	// Traverse the directory and process Go source files
	err := fileutil.WalkMatch(root, pattern, func(path string, info os.FileInfo) error {
		if fileutil.IsGoFile(path, info) {
			filePatch, err := processFile(path, transforms)
			if err != nil {
				return err
			}
			p.FileDiffs = append(p.FileDiffs, filePatch.FileDiffs...)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error processing files: %v", err)
	}

	// Write the patch to the output file
	err = p.WriteToFile(output)
	if err != nil {
		log.Fatalf("Error writing patch to file: %v", err)
	}

	fmt.Printf("Patch written to %s\n", output)
}

func processFile(filePath string, transforms []asttransform.TransformFunc) (patch.Patch, error) {
	// Load the file's AST once to create a backup
	astOrig, err := asttransform.ParseGoFile(filePath)
	if err != nil {
		return patch.Patch{}, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// Load the file's AST for transformation
	astFile, err := asttransform.ParseGoFile(filePath)
	if err != nil {
		return patch.Patch{}, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// Apply the AST transformations
	modifiedAST, err := asttransform.ApplyTransforms(astFile, transforms)
	if err != nil {
		return patch.Patch{}, fmt.Errorf("failed to apply transformations to %s: %w", filePath, err)
	}

	// Generate a diff between the original and modified ASTs
	originalSrc, err := asttransform.Print(astOrig)
	if err != nil {
		return patch.Patch{}, fmt.Errorf("failed to print original AST for %s: %w", filePath, err)
	}

	fmt.Println("original src:")
	fmt.Println(originalSrc)

	modifiedSrc, err := asttransform.Print(modifiedAST.(*ast.File))
	if err != nil {
		return patch.Patch{}, fmt.Errorf("failed to print modified AST for %s: %w", filePath, err)
	}

	fmt.Println("modified src:")
	fmt.Println(modifiedSrc)

	diff, err := patch.GenerateDiff(filePath, originalSrc, filePath, modifiedSrc)
	if err != nil {
		return patch.Patch{}, fmt.Errorf("failed to generate diff for %s: %w", filePath, err)
	}

	p := patch.Patch{}
	p.AddFileDiff(patch.FileDiff{
		OriginalPath: filePath,
		ModifiedPath: filePath,
		Diff:         diff,
	})

	return p, nil
}

// renameFunc is a sample AST transformation function that renames a function with the old name to the new name
func renameFunc(oldName, newName string) asttransform.TransformFunc {
	return func(node ast.Node) (ast.Node, error) {
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == oldName {
					x.Name.Name = newName
				}
			case *ast.Ident:
				if x.Name == oldName {
					x.Name = newName
				}
			}
			return true
		})
		return node, nil
	}
}

func prettyPrintAST(node ast.Node) {
	fset := token.NewFileSet()
	ast.Print(fset, node)
}
