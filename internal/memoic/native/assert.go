package native

import (
	"encoding/json"
	"errors"
	"fmt"
	"memoic/internal/memoic"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type typeCheck func(value string) (bool, error)

func createTypeCheck(fn func(value string) error) typeCheck {
	return func(value string) (bool, error) {
		if err := fn(value); err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
}

var typeChecks = map[string]typeCheck{
	"json": createTypeCheck(func(value string) error {
		var t any
		return json.Unmarshal([]byte(value), &t)
	}),
	"duration": createTypeCheck(func(value string) error {
		_, err := time.ParseDuration(value)
		return err
	}),
	"int": createTypeCheck(func(value string) error {
		_, err := strconv.Atoi(value)
		return err
	}),
	"uint": createTypeCheck(func(value string) error {
		_, err := strconv.ParseUint(value, 10, 64)
		return err
	}),
	"bool": createTypeCheck(func(value string) error {
		if strings.EqualFold(value, "true") || strings.EqualFold(value, "false") {
			return nil
		}
		return fmt.Errorf("%s is not a boolean", value)
	}),
	"float": createTypeCheck(func(value string) error {
		_, err := strconv.ParseFloat(value, 64)
		return err
	}),
	"url": createTypeCheck(func(value string) error {
		_, err := url.Parse(value)
		return err
	}),
}

var _ = memoic.AddFunction("assert", memoic.Function{
	{
		Invoke: func(stack *memoic.Stack) (any, error) {
			asserts, ok := (*stack.RawValue()).([]any)
			if !ok {
				return nil, errors.New("cannot find array of asserts for `assert`")
			}
			for _, assert := range asserts {
				assert := assert.(map[string]any)
				v, ok := assert["value"]
				if !ok {
					return nil, errors.New("cannot find `value` in `assert.asserts.$`")
				}
				value, ok := v.(string)
				if !ok {
					return nil, errors.New("`value` in `assert.asserts.$` isn't a string")
				}
				regex := ""
				if pattern, ok := assert["regex"]; ok {
					if _, ok := pattern.(string); ok {
						regex = pattern.(string)
					} else {
						return nil, errors.New("`regex` in `assert.asserts.$` is not a string")
					}
				}
				_type, ok := assert["type"]
				if ok {
					typedef, ok := _type.(string)
					if !ok {
						return nil, errors.New("`type` in `assert.asserts.$` is not a string")
					}
					typedef = strings.ToLower(typedef)
					if check, ok := typeChecks[typedef]; ok {
						ok, err := check(value)
						if !ok {
							fail := fmt.Errorf("failed to assert %s as %s", value, typedef)
							if err != nil {
								return nil, errors.Join(fail, err)
							}
							return nil, fail
						}
						if ok {
							continue
						}
					} else {
						return nil, errors.New("`type` in `assert.asserts.$` is not a valid type")
					}
				}
				if regex == "" {
					return nil, errors.New("`regex` or `type` in `assert.asserts.$` is not found")
				}
				ok, err := regexp.MatchString(regex, value)
				if !ok {
					fail := fmt.Errorf("failed to match %s with %s", value, regex)
					if err != nil {
						return nil, errors.Join(fail, err)
					}
					return nil, fail
				}
			}
			return nil, nil
		},
	},
})
