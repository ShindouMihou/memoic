package native

import (
	"errors"
	"fmt"
	"memoic/internal/memoic"
	"memoic/pkg/memoize"
	"strings"
)

var invalidReturnValueErr = errors.New("value of `return` should be a string directing to only one value (such as `$local.result`)")
var tooManyDirectivesErr = errors.New("value of `return` should only be one (such as `$local.result`)")

var _ = memoic.AddFunction("return", memoic.Function{
	{
		Invoke: func(stack *memoic.Stack) (any, error) {
			val, ok := (*stack.RawValue()).(string)
			if !ok {
				return nil, invalidReturnValueErr
			}
			directives := memoize.Directors(val)
			if len(directives) == 0 {
				return nil, invalidReturnValueErr
			}
			if len(directives) > 1 {
				return nil, tooManyDirectivesErr
			}
			directive := directives[0]
			value, ok := stack.Pull(directive)
			if !ok {
				return nil, fmt.Errorf("cannot find any value for `%s.%s`", directive.Director, strings.Join(directive.Keys, "."))
			}
			stack.Runtime.Result = &value
			return nil, nil
		},
	},
})
