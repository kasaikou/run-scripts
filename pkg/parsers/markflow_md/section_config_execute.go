package markflow_md

type sectionConfigExecute struct {
	Aliases        []string `yaml:"alias"`
	PrevExecutions []string `yaml:"prev"`
	Environments   []string `yaml:"environments"`
}
