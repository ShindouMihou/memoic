package memoic

import (
	"memoic/pkg/memoize"
	"reflect"
	"strings"
)

// Runtime contains the current Runtime information such as the parameters,
// the heap, not to be confused by Stack's sector and the end result of the
// Runtime.
type Runtime struct {
	Parameters Sector
	Heap       Sector
	Stacks     []Stack
	Result     *any
}

// Stack is the entity that contains the memory maps of a Pipe.
// A stack exists to prevent functions from accidentally overriding
// one another, or reusing the variables of one another.
type Stack struct {
	Sector     Sector
	Parameters Sector
	Runtime    *Runtime
}

func (runtime *Runtime) newStack() *Stack {
	stack := Stack{Runtime: runtime}
	runtime.Stacks = append(runtime.Stacks, stack)
	return &stack
}

func (stack *Stack) Interpolate() {
	stack.recursiveInterpolate(stack.Parameters)
}

func (stack *Stack) recursiveInterpolate(sector Sector) {
	for key, value := range sector {
		if text, ok := value.(string); ok {
			directives := memoize.InterpolatingDirectors(text)
			for _, directive := range directives {
				mem := stack.Runtime.Heap
				if strings.EqualFold(directive.Director, memoize.GlobalDirector) {
					mem = GlobalSector
				}
				if strings.EqualFold(directive.Director, memoize.ParamsDirector) {
					mem = stack.Runtime.Parameters
				}
				item, ok := mem[directive.Keys[0]]
				if !ok {
					continue
				}
				if len(directive.Keys) > 1 {
					for _, key := range directive.Keys[1:] {
						item = reflect.ValueOf(item).FieldByName(key).Interface()
					}
				}
				var result string
				if directive.As != nil {
					as := strings.ToLower(*directive.As)
					if marshal, ok := marshalers[as]; ok {
						bytes, err := marshal(item)
						if err != nil {
							continue
						}
						result = string(bytes)
					}
				}
				if result == "" {
					result = reflect.ValueOf(item).String()
				}
				sector[key] = strings.ReplaceAll(text, directive.Source, result)
			}
			if sector, ok := value.(Sector); ok {
				stack.recursiveInterpolate(sector)
			}
		}
	}
}

func newRuntime(parameters Sector) *Runtime {
	return &Runtime{Parameters: parameters}
}
