package markflow_md

import "github.com/kasaikou/markflow/pkg/models"

// Parse parses markflow.md file and register into dest.
func Parse(dest *models.Project, workingDir string, body []byte) error {
	res, err := parseToSections(body)
	if err != nil {
		return err
	}

	err = buildProject(buildProjectRequest{
		dest:       dest,
		sections:   res.sections,
		workingDir: workingDir,
	})
	if err != nil {
		return err
	}

	return nil
}
