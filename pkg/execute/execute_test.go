package execute

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestMidiumUnit_Execute(t *testing.T) {

	wd := util.Must(os.Getwd())

	testCases := []struct {
		Execution        models.Execution
		ExpectedExitCode int
		ExpectedStdout   string
		ExpectedStderr   string
		ExpectedEnviron  []string
	}{
		{
			Execution: models.Execution{
				ID:         uuid.Must(uuid.NewV7()),
				Name:       util.Must(models.ValidateExecutionName(`testExecuteContext['echo "Hello world" to stdout']`)),
				Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
				Script:     `echo "Hello world"`,
				WorkingDir: wd,
			},
			ExpectedExitCode: 0,
			ExpectedStdout:   "Hello world\n",
			ExpectedStderr:   "",
			ExpectedEnviron:  []string{},
		},
		{
			Execution: models.Execution{
				ID:         uuid.Must(uuid.NewV7()),
				Name:       util.Must(models.ValidateExecutionName(`testExecuteContext['export TEST_ENV="test"]`)),
				Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
				Script:     `echo 'TEST_ENV=test' >> ${MARKFLOW_EXPORT}`,
				WorkingDir: wd,
			},
			ExpectedExitCode: 0,
			ExpectedStdout:   "",
			ExpectedStderr:   "",
			ExpectedEnviron:  []string{"TEST_ENV=test"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("executionName=%s", testCase.Execution.Name.String()), func(t *testing.T) {
			t.Parallel()
			execution, err := BuildExecution(testCase.Execution)
			if !assert.NoError(t, err) {
				return
			}

			stdoutBuffer := bytes.NewBuffer([]byte{})
			stderrBuffer := bytes.NewBuffer([]byte{})
			wg := sync.WaitGroup{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				util.Must(io.Copy(stdoutBuffer, execution.Stdout()))
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				util.Must(io.Copy(stderrBuffer, execution.Stderr()))
			}()

			exitCode, err := execution.Execute(context.Background())

			if assert.NoError(t, err) {
				assert.Equal(t, testCase.ExpectedExitCode, exitCode)
				wg.Wait()
				assert.Equal(t, testCase.ExpectedStdout, stdoutBuffer.String())
				assert.Equal(t, testCase.ExpectedStderr, stderrBuffer.String())
				assert.Equal(t, testCase.ExpectedEnviron, execution.Environ())
			}
		})
	}

}
