package memoic

import (
	"errors"
	"memoic/pkg/memoize"
	"reflect"
	"strconv"
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
	if len(directive.Keys) > 1 && ok {
		for _, key := range directive.Keys[1:] {
			reflection := reflect.ValueOf(item)
			if !reflection.IsValid() || reflection.IsNil() {
				return nil, true
			}
			kind := reflection.Kind()
			if kind != reflect.Struct && !(kind == reflect.Array || kind == reflect.Slice) && kind != reflect.Map {
				return nil, false
			}
			if kind == reflect.Array || kind == reflect.Slice {
				index, err := strconv.Atoi(key)
				if err != nil {
					return nil, false
				}
				value := reflection.Index(index)
				if !value.IsValid() || value.IsNil() {
					return nil, false
				}
				item = value.Interface()
				continue
			}
			if kind == reflect.Map {
				value := reflection.MapIndex(reflect.ValueOf(key))
				if !value.IsValid() || value.IsNil() {
					return nil, false
				}
				item = value.Interface()
				continue
			}
			value := reflection.FieldByName(key)
			if !value.IsValid() || value.IsNil() {
				return nil, false
			}
			item = value.Interface()
		}
	}
	return item, ok
}
