package fileutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/albertywu/goshift/pkg/fileutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalkMatch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "walkmatch")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	err = os.WriteFile(filepath.Join(tempDir, "file1.go"), []byte("package main"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("Hello, world!"), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(subDir, "file3.go"), []byte("package main"), 0644)
	require.NoError(t, err)

	visited := make(map[string]bool)

	err = fileutil.WalkMatch(tempDir, "*.go", func(path string, info os.FileInfo) error {
		visited[path] = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, visited[filepath.Join(tempDir, "file1.go")])
	assert.False(t, visited[filepath.Join(tempDir, "file2.txt")])
	assert.True(t, visited[filepath.Join(subDir, "file3.go")])
}

func TestIsGoFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "isgofile")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	goFilePath := filepath.Join(tempDir, "file.go")
	err = os.WriteFile(goFilePath, []byte("package main"), 0644)
	require.NoError(t, err)

	txtFilePath := filepath.Join(tempDir, "file.txt")
	err = os.WriteFile(txtFilePath, []byte("Hello, world!"), 0644)
	require.NoError(t, err)

	goFileInfo, err := os.Stat(goFilePath)
	require.NoError(t, err)

	txtFileInfo, err := os.Stat(txtFilePath)
	require.NoError(t, err)

	assert.True(t, fileutil.IsGoFile(goFilePath, goFileInfo), "Expected IsGoFile to return true for .go files")
	assert.False(t, fileutil.IsGoFile(txtFilePath, txtFileInfo), "Expected IsGoFile to return false for non-.go files")
}
