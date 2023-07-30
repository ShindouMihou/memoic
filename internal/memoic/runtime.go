package memoic

import (
	"errors"
	"memoic/pkg/memoize"
	"reflect"
	"strings"
)

// Runtime contains the current Runtime information such as the parameters,
// the heap, not to be confused by Stack's sector and the end result of the
// Runtime. Unlike Stack, the heap is protected and can only be written by
// the runtime.
type Runtime struct {
	Parameters Sector
	heap       Sector
	stacks     []Stack
	Result     *any
}

// Get retrieves an item from the heap.
func (runtime *Runtime) Get(key string) any {
	if value, ok := runtime.heap[key]; ok {
		return value
	}
	return nil
}

// Load loads a function with the given parameters.
func (runtime *Runtime) Load(fn *Function) error {
	for _, pipe := range *fn {
		stack := runtime.newStack()
		stack.value = pipe.Value
		stack.interpolate()
		result, err := pipe.Invoke(stack)
		if err != nil {
			return err
		}
		if result != nil && pipe.As != nil {
			runtime.heap[*pipe.As] = result
		}
	}
	return nil
}

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

func (runtime *Runtime) newStack() *Stack {
	stack := Stack{Runtime: runtime}
	runtime.stacks = append(runtime.stacks, stack)
	return &stack
}

func (stack *Stack) interpolate() {
	if params, ok := stack.params(); ok {
		stack.recursiveInterpolate(params)
	} else {
		if text, ok := (*stack.RawValue()).(string); ok {
			stack.replaceValue(stack.textInterpolate(text))
		} else {
			stack.typedInterpolate(stack.RawValue())
		}
	}
}

// typedInterpolate interpolates the value with its expected value.
// this does not support string as value, refer to textInterpolate instead.
func (stack *Stack) typedInterpolate(value *any) {
	if sector, ok := (*value).(Sector); ok {
		stack.recursiveInterpolate(sector)
	}
	if array, ok := (*value).([]string); ok {
		for index, val := range array {
			val := val
			array[index] = stack.textInterpolate(val)
		}
	}
	if array, ok := (*value).([]any); ok {
		for index, val := range array {
			val := val
			if sector, ok := val.(map[string]any); ok {
				stack.recursiveInterpolate(sector)
				array[index] = sector
			}
		}
	}
}

func (stack *Stack) replaceValue(value any) {
	stack.value = &value
}

func (stack *Stack) recursiveInterpolate(sector Sector) {
	for key, value := range sector {
		value := value
		// control for string interpolation has to be inlined.
		// since `typedInterpolate` needs something directly modifiable.
		if text, ok := value.(string); ok {
			value = stack.textInterpolate(text)
		} else {
			stack.typedInterpolate(&value)
		}
		sector[key] = value
	}
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

func (stack *Stack) textInterpolate(text string) string {
	directives := memoize.InterpolatingDirectors(text)
	for _, directive := range directives {
		item, ok := stack.Pull(directive.Directive)
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
		text = strings.ReplaceAll(text, directive.Source, result)
	}
	return text
}

func NewRuntime(parameters Sector) *Runtime {
	return &Runtime{Parameters: parameters, heap: Sector{}}
}
