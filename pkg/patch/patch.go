package patch

import (
	"os"

	"github.com/pmezard/go-difflib/difflib"
)

// Patch is a collection of FileDiffs that represent changes made to multiple files.
type Patch struct {
	FileDiffs []FileDiff
}

// FileDiff represents the diff between the original and modified versions of a file,
// including the file paths and the diff text.
type FileDiff struct {
	OriginalPath string
	ModifiedPath string
	Diff         string
}

// AddFileDiff appends a FileDiff to a Patch's collection of FileDiffs.
func (p *Patch) AddFileDiff(fileDiff FileDiff) {
	p.FileDiffs = append(p.FileDiffs, fileDiff)
}

// WriteToFile writes a Patch's collection of FileDiffs to a file at the specified output path.
func (p *Patch) WriteToFile(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, fileDiff := range p.FileDiffs {
		_, err = file.WriteString(fileDiff.Diff)
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateDiff computes the diff between the original and modified text,
// returning the diff as a string or an error if the diff cannot be generated.
func GenerateDiff(originalPath, original, modifiedPath, modified string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(original),
		B:        difflib.SplitLines(modified),
		FromFile: originalPath,
		ToFile:   modifiedPath,
		Context:  3,
	}
	result, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", err
	}
	return result, nil
}
