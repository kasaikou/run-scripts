package tests

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/kasaikou/markflow/pkg/models"
)

var currentDir = func() string {
	_, callerFile, _, _ := runtime.Caller(1)
	return filepath.Dir(callerFile)
}()

// TestCaseMarkflowMD contains markdown file data and parsed result.
type TestCaseMarkflowMD struct {
	Name       string
	File       []byte
	WorkingDir string
	Project    models.ProjectJSONContent
}

// NewTestCaseMarkflowMD creates a new TestCaseMarkflowMD instance from markdown file and expected parsed result.
func NewTestCaseMarkflowMD(name string, project models.ProjectJSONContent) TestCaseMarkflowMD {

	bytes, err := os.ReadFile(filepath.Join(currentDir, name))
	if err != nil {
		panic(err)
	}

	return TestCaseMarkflowMD{
		Name:       name,
		File:       bytes,
		WorkingDir: currentDir,
		Project:    project,
	}
}
