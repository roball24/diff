package diff

import (
	"os"
)

// FileChecker is a convenience checker for files
type FileChecker struct {
	*SingleChecker
	f1, f2 *os.File
}

// NewFileChecker initializes a FileChecker from files
func NewFileChecker(fname1, fname2 string) (*FileChecker, error) {
	f1, err := os.Open(fname1)
	if err != nil {
		return nil, err
	}

	f2, err := os.Open(fname2)
	if err != nil {
		return nil, err
	}

	schk, err := NewSingleChecker(f1, f2)
	if err != nil {
		f1.Close()
		f2.Close()
		return nil, err
	}

	return &FileChecker{
		SingleChecker: schk,
		f1:            f1,
		f2:            f2,
	}, nil
}

// Close closes files used by checker
func (fchk *FileChecker) Close() {
	fchk.f1.Close()
	fchk.f2.Close()
}

// DefaultFileCheck does a FullEqual check on two files and writes output to os.Stdout
func DefaultFileCheck(fname1, fname2 string, verbose bool) (bool, error) {
	checker, err := NewFileChecker(fname1, fname2)
	if err != nil {
		return false, err
	}
	defer checker.Close()

	checker.AddCondition(FullEqual, ErrorFail)
	if !verbose {
		checker.SetWriter(nil)
	}
	return checker.Run(), nil
}
