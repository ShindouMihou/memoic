package memoic

import (
	"encoding/json"
	"memoic/pkg/memoize"
	"reflect"
	"strings"
)

type marshaler func(v any) ([]byte, error)

var marshalers = map[string]marshaler{
	"json": json.Marshal,
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
