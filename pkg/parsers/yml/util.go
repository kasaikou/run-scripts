package yml

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kasaikou/markflow/pkg/models"
	"gopkg.in/yaml.v3"
)

// NodeKind2String converts from [gopkg.in/yaml.v3.Kind] to string.
func NodeKind2String(kind yaml.Kind) string {
	switch kind {
	case yaml.AliasNode:
		return "AliasNode"
	case yaml.MappingNode:
		return "MappingNode"
	case yaml.DocumentNode:
		return "DocumentNode"
	case yaml.SequenceNode:
		return "SequenceNode"
	}

	panic("unknown kind")
}

// NodeMapContent contains value node associated with the key and key node,
type NodeMapContent struct {
	key       string
	keyNode   *yaml.Node
	valueNode *yaml.Node
}

// NodeMap contains pairs of value node associated keys.
type NodeMap struct {
	contents []NodeMapContent
}

// UnmarshalYAML parses from yaml's mapping node to NodeMap.
func (nm *NodeMap) UnmarshalYAML(node *yaml.Node) error {

	if node.Kind != yaml.MappingNode {
		return models.NewModelValidateError(fmt.Errorf("this is must be %s, but node is %s", NodeKind2String(yaml.MappingNode), NodeKind2String(node.Kind)))
	}

	nm.contents = make([]NodeMapContent, 0, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {

		nm.contents = append(nm.contents, NodeMapContent{
			key:       node.Content[i].Value,
			keyNode:   node.Content[i],
			valueNode: node.Content[i+1],
		})
	}

	slices.SortFunc(nm.contents, func(a, b NodeMapContent) int { return strings.Compare(a.key, b.key) })

	for i := 0; i < len(nm.contents)-1; i++ {
		if nm.contents[i].key == nm.contents[i+1].key {
			return models.NewModelValidateError(fmt.Errorf("'%s' is duplicated key", nm.contents[i].key))
		}
	}

	return nil
}

// Get gets yaml.Node associated with the key.
func (nm *NodeMap) Get(key string) (node *yaml.Node, exist bool) {
	idx, exist := slices.BinarySearchFunc(nm.contents, key, func(nmc NodeMapContent, s string) int { return strings.Compare(nmc.key, s) })
	if exist {
		return nm.contents[idx].valueNode, true
	} else {
		return nil, false
	}
}
