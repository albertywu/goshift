package fileutil

import (
	"os"
)

// FileVisitorFunc is a function that takes a file path and file info as input,
// and performs some action on the file. It returns an error if the action fails.
type FileVisitorFunc func(path string, info os.FileInfo) error

// WalkMatch recursively traverses a directory, applying the visitor function
// to all files and directories that match the specified glob pattern.
func WalkMatch(root string, pattern string, visitor FileVisitorFunc) error {
	// ...
	return nil
}

// IsGoFile checks if a given file path corresponds to a Go source file,
// based on the file's extension and whether it is a regular file.
func IsGoFile(path string, info os.FileInfo) bool {
	// ...
	return false
}
