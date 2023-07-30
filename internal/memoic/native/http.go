package native

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"memoic/internal/memoic"
)

var client = resty.New()

var _ = memoic.AddFunction("http.request", memoic.Function{
	func(stack *memoic.Stack) (any, error) {
		stack.Interpolate()

		request := client.R()

		if method, ok := stack.Parameters["method"]; ok {
			method, ok := method.(string)
			if !ok {
				return nil, errors.New("method can only be a string")
			}
			request.Method = method
		} else {
			return nil, errors.New("cannot find the method to send http request")
		}

		if link, ok := stack.Parameters["link"]; ok {
			link, ok := link.(string)
			if !ok {
				return nil, errors.New("link can only be a string")
			}
			request.URL = link
		} else {
			return nil, errors.New("cannot find the link to send http request")
		}

		if headers, ok := stack.Parameters["headers"]; ok {
			headers, ok := headers.(map[string]string)
			if !ok {
				return nil, errors.New("headers can only be string:string")
			}
			request.SetHeaders(headers)
		}

		if body, ok := stack.Parameters["body"]; ok {
			request.SetBody(body)
		}

		response, err := request.Send()
		if err != nil {
			return nil, err
		}
		return string(response.Body()), nil
	},
})
