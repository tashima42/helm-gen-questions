// SPDX-FileCopyrightText: SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package genquestions

import (
	"slices"
	"strconv"

	"github.com/tidwall/gjson"
	"go.yaml.in/yaml/v4"
)

func isNull(r gjson.Result) bool { return !r.Exists() || r.Type == gjson.Null }

// firstSet returns the first result that is set and not null or false, like jq's // operator.
func firstSet(results ...gjson.Result) gjson.Result {
	for _, r := range results {
		if !isNull(r) && r.Type != gjson.False {
			return r
		}
	}
	return gjson.Result{}
}

func scalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

// jsonNode converts a JSON value into a yaml node, keeping the key order.
func jsonNode(r gjson.Result) *yaml.Node {
	if !r.Exists() {
		return nil
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(r.Raw), &doc); err != nil || len(doc.Content) == 0 {
		return nil
	}
	return copyNode(doc.Content[0])
}

// resolve follows aliases.
func resolve(n *yaml.Node) *yaml.Node {
	for n != nil && n.Kind == yaml.AliasNode {
		n = n.Alias
	}
	return n
}

// copyNode deep copies a node, without the style, anchors and comments of the source.
func copyNode(n *yaml.Node) *yaml.Node {
	n = resolve(n)
	if n == nil {
		return nil
	}
	c := &yaml.Node{Kind: n.Kind, Tag: n.Tag, Value: n.Value}
	for _, child := range n.Content {
		c.Content = append(c.Content, copyNode(child))
	}
	return c
}

func mapGet(m *yaml.Node, key string) *yaml.Node {
	m = resolve(m)
	if m == nil || m.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return resolve(m.Content[i+1])
		}
	}
	return nil
}

// lookup gets the node at a key path, nil if it does not exist.
func lookup(n *yaml.Node, path []string) *yaml.Node {
	for _, segment := range path {
		n = resolve(n)
		switch {
		case n == nil:
			return nil
		case n.Kind == yaml.SequenceNode:
			i, err := strconv.Atoi(segment)
			if err != nil || i < 0 || i >= len(n.Content) {
				return nil
			}
			n = n.Content[i]
		default:
			n = mapGet(n, segment)
		}
	}
	return resolve(n)
}

// without copies a mapping node, leaving out the keys.
func without(m *yaml.Node, keys ...string) *yaml.Node {
	c := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if !slices.Contains(keys, m.Content[i].Value) {
			c.Content = append(c.Content, m.Content[i], m.Content[i+1])
		}
	}
	return c
}
