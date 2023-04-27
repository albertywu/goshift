# goshift
A toolkit for go codemods

## Overview

goshift is a command-line tool that performs user-defined AST transformations on a directory of Go source code files. The tool will traverse the directory structure recursively, applying the specified transformation functions to all matching files. The output will be a single patch file containing all the changes made by the AST transformations.

## Features

- Apply one or more AST transformations to Go source code files in a directory.
- Traverse directories recursively, processing all nested Go files.
- Support file matching using glob patterns to include or exclude specific files and folders.
- Generate a single patch file containing all changes made by the transformations.

## Usage

```
goshift [options] <source_directory> <output_patch_file>
```

### Options

- `-t, --transform`: Path to a Go source code file containing one or more AST transformation functions. This option can be specified multiple times to include multiple transformation functions.
- `-g, --glob`: A glob pattern to filter the files and folders to be processed. By default, all `.go` files in the specified directory and its subdirectories will be processed.

### Transformation Functions

The AST transformation functions should be written in Go and have the following signature:

```go
func Transform(node ast.Node) (ast.Node, error)
```

The `Transform` function takes an `ast.Node` as input, representing the current node in the AST being processed. The function should return a modified `ast.Node` or an error if the transformation cannot be applied.

### Examples

Apply a single transformation to all Go files in the `./src` directory and output the changes to `changes.patch`:

```
goshift -t transform_function.go ./src changes.patch
```

Apply multiple transformations to Go files in the `./src` directory, but only process files inside the `./src/models` and `./src/controllers` folders:

```
goshift -t transform_function1.go -t transform_function2.go -g "./src/{models,controllers}/**/*.go" ./src changes.patch
```

## Implementation

The Go AST Transformer tool can be implemented using the following steps:

1. Parse the command-line arguments, including the source directory, output patch file, transformation function files, and file matching glob pattern.
2. Load and compile the transformation function(s) using the `plugin` package, which allows loading Go functions from source code files at runtime.
3. Traverse the source directory recursively, processing all files that match the specified glob pattern.
    - For each matching file, use the `go/parser` package to parse the file into an AST.
    - Apply the transformation function(s) to the AST using the `ast.Inspect` function.
    - Use the `go/format` package to convert the modified AST back to Go source code.
    - Generate a diff between the original and transformed source code using a library like `github.com/sergi/go-diff/diffmatchpatch`.
    - Append the diff to the output patch file.
4. Save the output patch file to the specified location.

## Dependencies

- The Go standard library packages: `go/parser`, `go/token`, `go/ast`, `go/format`
- [github.com/sergi/go-diff](https://github.com/sergi/go-diff): A library to generate diffs between text files.
- The `plugin` package from the Go standard library for loading transformation functions from source code files. Note that this package is only supported on certain platforms (e.g., Linux and macOS).