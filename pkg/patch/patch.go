package patch

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
	// ...
}

// WriteToFile writes a Patch's collection of FileDiffs to a file at the specified output path.
func (p *Patch) WriteToFile(outputPath string) error {
	// ...
	return nil
}

// GenerateDiff computes the diff between the original and modified text,
// returning the diff as a string or an error if the diff cannot be generated.
func GenerateDiff(original, modified string) (string, error) {
	// ...
	return "", nil
}
