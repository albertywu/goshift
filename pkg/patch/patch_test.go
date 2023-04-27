package patch_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/albertywu/goshift/pkg/patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatch_WriteToFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "patchtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "patchfile.diff")

	p := patch.Patch{}
	fileDiff := patch.FileDiff{
		OriginalPath: "original.go",
		ModifiedPath: "modified.go",
		Diff: `--- original.go
+++ modified.go
@@ -1,5 +1,5 @@
 package main

-func oldName() {
+func newName() {
 	println("Hello, world!")
 }
`,
	}
	p.AddFileDiff(fileDiff)

	err = p.WriteToFile(outputPath)
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	assert.Equal(t, fileDiff.Diff, string(content), "Expected written diff to match the added FileDiff")
}

func TestGenerateDiff(t *testing.T) {
	original := `package main

func oldName() {
	println("Hello, world!")
}

func main() {
	oldName()
}
`

	modified := `package main

func newName() {
	println("Hello, world!")
}

func main() {
	newName()
}
`

	expectedDiff := `--- original.go
+++ modified.go
@@ -1,10 +1,10 @@
 package main
 
-func oldName() {
+func newName() {
 	println("Hello, world!")
 }
 
 func main() {
-	oldName()
+	newName()
 }
`

	originalPath := "original.go"
	modifiedPath := "modified.go"

	diff, err := patch.GenerateDiff(originalPath, original, modifiedPath, modified)
	require.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(expectedDiff), strings.TrimSpace(diff), "Expected generated diff to match the expected diff")
}
