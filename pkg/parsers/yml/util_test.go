package yml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestSmallUnit_Node2NodeMap(t *testing.T) {

	body := ([]byte)(`
version: 1.2
hoo:
  alphabet: abc
  number: 123
  list: [1, 2, 3]
bar: null
`)

	nodeMap := NodeMap{}
	if assert.NoError(t, yaml.Unmarshal(body, &nodeMap)) {
		_, exist := nodeMap.Get("version")
		assert.Truef(t, exist, "version should be existed")
		_, exist = nodeMap.Get("hoo")
		assert.Truef(t, exist, "hoo should be existed")
		_, exist = nodeMap.Get("bar")
		assert.Truef(t, exist, "bar should be existed")
		_, exist = nodeMap.Get("piyo")
		assert.Falsef(t, exist, "piyo should not be existed")
	}

}
