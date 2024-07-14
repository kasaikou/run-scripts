package markflow_md

import (
	"fmt"

	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/parsers/md"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

type section struct {
	Heading     md.Heading
	Description []md.TextBlock
	Config      sectionConfig
	Script      *md.FencedCodeBlock
	BeginAt     int
	EndAt       int
}

func wrapModelValidateErrorWithSection(s section, err error) error {
	return models.WrapModelValidateError(fmt.Sprintf("sections['%s']", s.Heading.Plaintext()), err)
}

type parse2SectionsResponse struct {
	sections []section
}

func parseToSections(src []byte) (res parse2SectionsResponse, err error) {

	var crtSection *section
	node := goldmark.DefaultParser().Parse(text.NewReader(src))

	for node := node.FirstChild(); node != nil; node = node.NextSibling() {
		switch node.Kind() {
		case ast.KindHeading:
			beginAt := node.Lines().At(0).Start

			if crtSection != nil {
				crtSection.EndAt = beginAt
			}

			res.sections = append(res.sections, section{
				Heading: md.NewHeading(node.(*ast.Heading), src),
				BeginAt: beginAt,
			})

			crtSection = &res.sections[len(res.sections)-1]

		case ast.KindParagraph:
			if crtSection != nil {
				crtSection.Description = append(crtSection.Description, md.NewParagraph(node.(*ast.Paragraph), src))
			}

		case ast.KindFencedCodeBlock:
			codeBlock := md.NewFencedCodeBlock(node.(*ast.FencedCodeBlock), src)
			lang := codeBlock.Lang()

			switch lang {
			case "yaml", "yml":
				yamlConfig := codeBlock.CodeBytes()
				if err := yaml.Unmarshal(yamlConfig, &crtSection.Config); err != nil {
					return res, wrapModelValidateErrorWithSection(*crtSection, models.WrapModelValidateError("config", err))
				}
			default:
				if crtSection.Script != nil {
					return res, wrapModelValidateErrorWithSection(*crtSection, models.WrapModelValidateError("script",
						fmt.Errorf("multiple script in a section is not allowed")))
				}
				crtSection.Script = &codeBlock
			}
		}
	}

	crtSection.EndAt = len(src)
	if len(res.sections) == 0 {
		return res, nil
	}

	sections := make([]section, 0, len(res.sections))
	for i := 0; i < len(res.sections); i++ {
		if res.sections[i].Config.Kind == "" && res.sections[i].Script == nil {
			if len(sections) > 0 {
				sections[len(sections)-1].Description = append(sections[len(sections)-1].Description, res.sections[i].Heading)
				sections[len(sections)-1].Description = append(sections[len(sections)-1].Description, res.sections[i].Description...)
				sections[len(sections)-1].EndAt = res.sections[i].EndAt
			}
		} else {
			sections = append(sections, res.sections[i])
		}
	}

	res.sections = sections
	return res, nil
}
