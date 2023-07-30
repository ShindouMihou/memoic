package native

import (
	"fmt"
	"memoic/internal/memoic"
)

var _ = memoic.AddFunction("std.log", memoic.Function{
	func(stack *memoic.Stack) (any, error) {
		stack.Interpolate()
		fmt.Println(stack.Parameters["message"])
		return nil, nil
	},
})
