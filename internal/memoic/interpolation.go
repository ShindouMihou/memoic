package memoic

import (
	"fmt"
	"github.com/bytedance/sonic"
	"memoic/pkg/memoize"
	"reflect"
	"strconv"
	"strings"
)

type marshaler func(v any) ([]byte, error)

var marshalers = map[string]marshaler{
	"json": sonic.Marshal,
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
			value := reflect.ValueOf(item)
			switch value.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result = strconv.FormatInt(value.Int(), 10)
			case reflect.String:
				result = value.String()
			case reflect.Bool:
				result = strconv.FormatBool(value.Bool())
			case reflect.Float32, reflect.Float64:
				result = fmt.Sprintf("%f", value.Float())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				result = strconv.FormatUint(value.Uint(), 10)
			default:
				bytes, err := sonic.Marshal(item)
				if err != nil {
					continue
				}
				result = string(bytes)
			}
		}
		text = strings.ReplaceAll(text, directive.Source, result)
	}
	return text
}
