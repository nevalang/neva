package ast

import "github.com/nevalang/neva/pkg/core"

// Comments is a parsed comment payload attached to an entity.
type Comments struct {
	TextBlocks []string          `json:"textBlocks,omitempty"`
	Inports    map[string]string `json:"inports,omitempty"`
	Outports   map[string]string `json:"outports,omitempty"`
	Examples   []string          `json:"examples,omitempty"`
	Tags       []CommentTag      `json:"tags,omitempty"`
	Meta       core.Meta         `json:"meta"`
}

type CommentTag struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
