package markflow_md

import (
	"fmt"
	"testing"

	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/parsers/markflow_md/tests"
	"github.com/stretchr/testify/assert"
)

func TestMediumIntegration_ParseMarkflowMD(t *testing.T) {

	testCases := []struct {
		TestCaseMarkflowMD tests.TestCaseMarkflowMD
	}{
		{
			TestCaseMarkflowMD: tests.TestCaseExample,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("name=%s", testCase.TestCaseMarkflowMD.Name), func(t *testing.T) {
			t.Parallel()

			project := models.NewProject()
			if assert.NoError(t, Parse(project, testCase.TestCaseMarkflowMD.WorkingDir, testCase.TestCaseMarkflowMD.File)) {
				assert.Equal(t, testCase.TestCaseMarkflowMD.Project, project.ToJSONContent())
			}
		})
	}
}
