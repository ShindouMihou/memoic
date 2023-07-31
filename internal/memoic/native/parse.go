package native

import (
	"errors"
	"github.com/bytedance/sonic"
	"memoic/internal/memoic"
)

var _ = memoic.AddFunction("json.parse", memoic.Function{
	{
		Invoke: func(stack *memoic.Stack) (any, error) {
			v := stack.RawValue()
			if v == nil {
				return nil, errors.New("cannot find `value` in `json.parse`")
			}
			value, ok := (*v).(string)
			if !ok {
				return nil, errors.New("`value` in `json.parse` isn't a string")
			}
			var t any
			err := sonic.UnmarshalString(value, &t)
			if err != nil {
				return nil, err
			}
			return t, nil
		},
	},
})
