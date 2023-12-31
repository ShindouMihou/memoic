package memoic

import (
	"errors"
	"fmt"
	"memoic/pkg/memoize"
	"strings"
)

var reservedKeys = []string{
	"return",
	"assert",
}

var functions = make(map[string]Function)

type Pipe struct {
	Invoke func(stack *Stack) (any, error)
	Value  *any
	As     *string
}
type Function []Pipe

func AddFunction(key string, function Function) bool {
	for _, reserved := range reservedKeys {
		if fn := Get(reserved); fn != nil && strings.EqualFold(key, reserved) {
			panic(fmt.Errorf("%s is a reserved function key", reserved))
		}
	}
	functions[key] = function
	return true
}

func Get(key string) *Function {
	if fn, ok := functions[key]; ok {
		return &fn
	}
	return nil
}

func LoadFunctions(root *memoize.Root) error {
	pkg := root.Metadata.Package
	if pkg == "" {
		return errors.New("package cannot be empty")
	}
	for _, declaration := range root.Functions {
		declaration := declaration

		var fn []Function
		fn = linkPipes(fn, declaration.Pipeline)

		// flatten the function into one level
		var fnc Function
		for _, function := range fn {
			for _, pipe := range function {
				fnc = append(fnc, pipe)
			}
		}
		AddFunction(pkg+"."+declaration.Name, fnc)
	}
	return nil
}

func linkPipes(base []Function, pipes []*memoize.Pipe) []Function {
	for _, child := range pipes {
		if len(child.Pipes) > 0 {
			base = linkPipes(base, child.Pipes)
		}
		if !strings.EqualFold(child.Director, memoize.FunctionDirector) {
			fmt.Println(child.Keys, " is not a function pipeline.")
			continue
		}
		key := strings.Join(child.Keys, ".")
		fn, ok := functions[key]
		if !ok {
			fmt.Println(key, " has no linked function.")
			continue
		}

		cpy := make(Function, len(fn))
		copy(cpy, fn)
		fn = cpy

		if child.As != nil {
			fn[len(fn)-1].As = child.As
		}

		if child.Value != nil {
			fn[len(fn)-1].Value = child.Value
		}

		base = append(base, fn)
	}
	return base
}
