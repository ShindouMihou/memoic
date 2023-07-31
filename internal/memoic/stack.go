package memoic

import (
	"errors"
	"memoic/pkg/memoize"
	"strings"
)

// Stack is the entity that contains the memory maps of a Pipe.
// A stack exists to prevent functions from accidentally overriding
// one another, or reusing the variables of one another.
type Stack struct {
	Sector  Sector
	value   *any
	Runtime *Runtime
}

// RawValue retrieves the raw value, this can be null or anything.
func (stack *Stack) RawValue() *any {
	return stack.value
}

func (stack *Stack) params() (Sector, bool) {
	params, ok := (*stack.value).(map[string]any)
	return params, ok
}

// MappedParameters gets the value of this function as a Sector (otherwise, a map of anything).
func (stack *Stack) MappedParameters() (Sector, error) {
	if params, ok := stack.params(); ok {
		return params, nil
	}
	return nil, errors.New("the function's value is not a map")
}

func (stack *Stack) replaceValue(value any) {
	stack.value = &value
}

func (stack *Stack) Pull(directive memoize.Directive) (any, bool) {
	mem := stack.Runtime.heap
	if strings.EqualFold(directive.Director, memoize.GlobalDirector) {
		mem = GlobalSector
	}
	if strings.EqualFold(directive.Director, memoize.ParamsDirector) {
		mem = stack.Runtime.Parameters
	}
	item, ok := mem[directive.Keys[0]]
	return item, ok
}
