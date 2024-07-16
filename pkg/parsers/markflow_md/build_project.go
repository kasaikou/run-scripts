package markflow_md

import (
	"errors"

	"github.com/kasaikou/markflow/pkg/models"
)

type buildProjectRequest struct {
	sections   []section
	workingDir string
	dest       *models.Project
}

func buildProject(req buildProjectRequest) error {

	var globalSection *section
	pipelineSections := make([]*section, 0, len(req.sections))
	executionSections := make([]*section, 0, len(req.sections))

	for _, section := range req.sections {
		section := section

		if section.Config.Kind == "global" {
			if globalSection != nil {
				return wrapModelValidateErrorWithSection(section, errors.New("global config must be one section in one markdown, not duplicated"))
			}
			globalSection = &section

		} else if section.Config.Kind == "pipeline" {
			pipelineSections = append(pipelineSections, &section)
		} else {
			executionSections = append(executionSections, &section)
		}
	}

	// TODO: default execution and pipeline from global config.

	for _, section := range executionSections {
		execution := models.NewExecution()

		if name, err := models.ValidateExecutionName(section.Heading.Plaintext()); err != nil {
			return wrapModelValidateErrorWithSection(*section, err)
		} else {
			execution.Name = name
		}

		execution.Descriptions = make([]string, 0, len(section.Description))
		for _, desc := range section.Description {
			execution.Descriptions = append(execution.Descriptions, desc.Plaintext())
		}

		if lang, err := models.ValidateExecutionLanguage(section.Script.Lang()); err != nil {
			return wrapModelValidateErrorWithSection(*section, err)
		} else {
			execution.Lang = lang
		}

		execution.Script = section.Script.Code()
		execution.WorkingDir = req.workingDir

		if section.Config.Execute != nil {
			execution.Aliases = section.Config.Execute.Aliases
			execution.Environments = section.Config.Execute.Environments

			for _, prev := range section.Config.Execute.PrevExecutions {
				execution.PrevExecutions = append(execution.PrevExecutions, models.NewReferenceExecution(req.dest, prev))
			}
		}

		if err := req.dest.AddExecution(execution); err != nil {
			return wrapModelValidateErrorWithSection(*section, err)
		}
	}

	for _, section := range pipelineSections {
		pipeline := models.NewPipeline()

		if name, err := models.ValidatePipelineName(section.Heading.Plaintext()); err != nil {
			return wrapModelValidateErrorWithSection(*section, err)
		} else {
			pipeline.Name = name
		}

		pipeline.Aliases = section.Config.Pipeline.Aliases

		for _, step := range section.Config.Pipeline.Steps {
			pipelineStep := models.PipelineStep{}
			pipelineStep.Name = step.Name
			for _, execution := range step.executions {
				pipelineStep.Executions = append(pipelineStep.Executions, models.NewReferenceExecution(req.dest, execution))
			}
			pipeline.Steps = append(pipeline.Steps, pipelineStep)
		}

		if err := req.dest.AddPipeline(pipeline); err != nil {
			return wrapModelValidateErrorWithSection(*section, err)
		}
	}

	return nil
}
