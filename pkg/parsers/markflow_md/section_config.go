package markflow_md

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/parsers/yml"
	"gopkg.in/yaml.v3"
)

type sectionConfig struct {
	Kind     sectionConfigKind      `json:"kind"`
	Global   *sectionConfigGlobal   `json:"global,omitempty"`
	Execute  *sectionConfigExecute  `json:"execute,omitempty"`
	Pipeline *sectionConfigPipeline `json:"pipeline,omitempty"`
}

func (c *sectionConfig) UnmarshalYAML(node *yaml.Node) error {
	*c = sectionConfig{}

	mapping := yml.NodeMap{}
	if err := node.Decode(&mapping); err != nil {
		if _, ok := err.(models.ModelValidateError); ok {
			return err
		} else {
			return models.NewModelValidateError(err)
		}
	}

	kind, isKindExist := mapping.Get("kind")
	content, isContentExist := mapping.Get("content")

	if !isKindExist {
		return models.NewModelValidateError(fmt.Errorf("'kind' option is required"))
	} else if !isContentExist {
		return models.NewModelValidateError(fmt.Errorf("'content' option is required"))
	}

	if err := kind.Decode(&c.Kind); err != nil {
		return models.WrapModelValidateError("kind", err)
	}

	switch c.Kind {
	case "global":
		c.Global = &sectionConfigGlobal{}
		if err := content.Decode(c.Global); err != nil {
			return models.WrapModelValidateError("content", err)
		}
	case "execute":
		c.Execute = &sectionConfigExecute{}
		if err := content.Decode(c.Execute); err != nil {
			return models.WrapModelValidateError("content", err)
		}
	case "pipeline":
		c.Pipeline = &sectionConfigPipeline{}
		if err := content.Decode(c.Pipeline); err != nil {
			return models.WrapModelValidateError("content", err)
		}
	}

	return nil
}

type sectionConfigKind string

func (c *sectionConfigKind) UnmarshalYAML(content *yaml.Node) error {
	contentKind := content.Kind
	if contentKind != yaml.ScalarNode {
		return models.NewModelValidateError(fmt.Errorf("kind should be %s, but this is %s node", yml.NodeKind2String(yaml.ScalarNode), yml.NodeKind2String(contentKind)))
	}

	contentTag := content.ShortTag()
	if contentTag != "!!str" {
		return models.NewModelValidateError(fmt.Errorf("kind should be !!str tag, but this is %s tag", contentTag))
	}

	options := []string{"global", "execute", "pipeline"}
	if !slices.Contains(options, content.Value) {
		return models.NewModelValidateError(fmt.Errorf("kind should be one of type (%s), but found '%s'", strings.Join(options, ", "), content.Value))
	}

	*c = sectionConfigKind(content.Value)
	return nil
}
