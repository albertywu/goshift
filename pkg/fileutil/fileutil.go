package fileutil

import (
	"os"
	"path/filepath"
)

// FileVisitorFunc is a function that takes a file path and file info as input,
// and performs some action on the file. It returns an error if the action fails.
type FileVisitorFunc func(path string, info os.FileInfo) error

// WalkMatch recursively traverses a directory, applying the visitor function
// to all files and directories that match the specified glob pattern.
func WalkMatch(root string, pattern string, visitor FileVisitorFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}

		if matched {
			return visitor(path, info)
		}

		return nil
	})
}

// IsGoFile checks if a given file path corresponds to a Go source file,
// based on the file's extension and whether it is a regular file.
func IsGoFile(path string, info os.FileInfo) bool {
	return info.Mode().IsRegular() && filepath.Ext(path) == ".go"
}
