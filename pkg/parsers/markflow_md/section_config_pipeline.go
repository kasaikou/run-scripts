package markflow_md

type sectionConfigPipeline struct {
	Aliases []string                    `yaml:"aliases"`
	Steps   []sectionConfigPipelineStep `yaml:"steps"`
}

type sectionConfigPipelineStep struct {
	Name       string   `yaml:"name"`
	executions []string `yaml:"executes"`
}
