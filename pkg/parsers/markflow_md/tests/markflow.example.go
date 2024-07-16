package tests

import "github.com/kasaikou/markflow/pkg/models"

// TestCaseExample is expected parse result of ./markflow.example.co file.
var TestCaseExample = NewTestCaseMarkflowMD("markflow.example.md", models.ProjectJSONContent{
	Executions: []models.ExecutionJSONContent{
		{
			Name:               "シェルスクリプトの実行",
			Descriptions:       []string{"シェルで echo \"hello world\" を実行します。"},
			PrevExecutionNames: []string{},
			Lang:               "sh",
			Script:             "echo \"hello world\"\n",
			WorkingDir:         currentDir,
			Environments: []string{
				"TEST=environment",
			},
		},
		{
			Name:               "Python の実行",
			Descriptions:       []string{"Python で print(\"hello world\") を実行します。"},
			PrevExecutionNames: []string{},
			Lang:               "py",
			Script:             "print(\"hello world\")\n",
			WorkingDir:         currentDir,
		},
	},
	Pipelines: []models.PipelineJSONContent{},
})
