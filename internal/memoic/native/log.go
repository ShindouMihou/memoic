package native

import (
	"fmt"
	"memoic/internal/memoic"
)

var _ = memoic.AddFunction("log", memoic.Function{
	{
		Invoke: func(stack *memoic.Stack) (any, error) {
			fmt.Println(*stack.RawValue())
			return nil, nil
		},
	},
})
