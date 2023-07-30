package memoize

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Pipe struct {
	Directive
	As    *string
	Pipes []*Pipe
}

type Directive struct {
	Director string
	Keys     []string
	Value    *any
}

type SourcedDirective struct {
	Directive
	Source string
	As     *string
}

type Root struct {
	Metadata  Metadata
	Functions []FunctionDeclaration
}

type Metadata struct {
	Flags   []string `json:"flags"`
	Package string   `json:"package"`
}

type FunctionDeclaration struct {
	Name     string
	Pipeline []*Pipe
}

type declaration struct {
	FunctionDeclaration
	Pipes []map[string]any `json:"pipeline"`
}

func InterpolatingDirectors(src string) []SourcedDirective {
	var directives []SourcedDirective
	start, end, valid := -1, -1, false

	for index, char := range src {
		index, char := index, char

		if char == '{' {
			if start != -1 {
				continue
			}
			start = index
			valid = src[index+1] == '$'
		} else if char == '}' && start != -1 && valid {
			end = index
			source := src[start : end+1]

			// ignore those that repeat, so we don't duplicate it when we happen to
			// do the replacements.
			appended := false
			for _, value := range directives {
				if strings.EqualFold(value.Source, source) {
					appended = true
					break
				}
			}
			if appended {
				continue
			}

			selection := src[start+1 : end]

			if !strings.HasPrefix(selection, DirectorPrefix) || !strings.Contains(selection, Splitter) {
				continue
			}
			ident := strings.ToLower(selection)

			idents := strings.Split(ident[1:], " ")
			var as *string
			if len(idents) > 2 {
				if strings.EqualFold(idents[1], AsToken) {
					as = &idents[2]
				}
			}
			ident = idents[0]

			tokens := strings.Split(ident, Splitter)
			if len(tokens) < 2 {
				return nil
			}
			director := tokens[0]
			keys := tokens[1:]

			directives = append(directives, SourcedDirective{Source: source, As: as, Directive: Directive{
				Director: director,
				Keys:     keys,
			}})
			start, end, valid = -1, -1, false
		}
	}

	return directives
}

func Directors(src string) []Directive {
	var directives []Directive
	parts := strings.Split(src, " ")
	for _, part := range parts {
		if !strings.HasPrefix(part, DirectorPrefix) || !strings.Contains(part, Splitter) {
			continue
		}
		ident := strings.ToLower(part)[1:]
		tokens := strings.Split(ident, Splitter)
		if len(tokens) < 2 {
			return nil
		}
		director := tokens[0]
		keys := tokens[1:]
		directives = append(directives, Directive{
			Director: director,
			Keys:     keys,
			Value:    nil,
		})
	}
	return directives
}

func Parse(src []byte) (*Root, error) {
	table := make(map[string]json.RawMessage)
	if err := json.Unmarshal(src, &table); err != nil {
		return nil, err
	}
	root := Root{}
	for key, value := range table {
		if key == "$metadata" {
			if err := json.Unmarshal(value, &root.Metadata); err != nil {
				return nil, err
			}
			continue
		}
		var function declaration
		if err := json.Unmarshal(value, &function); err != nil {
			return nil, err
		}
		function.Name = key
		for _, pipe := range function.Pipes {
			for key, value := range pipe {
				pipe := parsePipe(key, value)
				if pipe != nil {
					function.Pipeline = append(function.Pipeline, pipe)
				}
			}
		}
		root.Functions = append(root.Functions, function.FunctionDeclaration)
	}
	return &root, nil
}

func parsePipe(key string, value any) *Pipe {
	root := markPipe(key)
	if root == nil {
		return nil
	}
	if children, ok := value.(map[string]any); ok {
		digPipes(root, children)
		if len(root.Pipes) == 0 {
			root.Value = &value
		}
	} else {
		root.Value = &value
	}
	return root
}

func digPipes(root *Pipe, children map[string]any) {
	for key, value := range children {
		child := markPipe(key)
		if child == nil {
			continue
		}
		if children, ok := value.(map[string]any); ok {
			digPipes(child, children)
			if len(child.Pipes) == 0 {
				child.Value = &value
			}
		} else {
			child.Value = &value
		}
		root.Pipes = append(root.Pipes, child)
	}
	if len(root.Pipes) == 0 {
		value := reflect.ValueOf(children).Interface()
		root.Value = &value
	}
}

func markPipe(key string) *Pipe {
	if !strings.HasPrefix(key, DirectorPrefix) || !strings.Contains(key, Splitter) {
		return nil
	}

	key = strings.ToLower(key)

	idents := strings.Split(key[1:], " ")
	var as *string
	if len(idents) > 2 {
		if strings.EqualFold(idents[1], AsToken) {
			as = &idents[2]
		}
	}
	tokens := strings.Split(idents[0], Splitter)
	if len(tokens) < 2 {
		return nil
	}
	director := tokens[0]
	keys := tokens[1:]
	return &Pipe{
		Directive: Directive{
			Director: director,
			Keys:     keys,
			Value:    nil,
		},
		As:    as,
		Pipes: []*Pipe{},
	}
}
