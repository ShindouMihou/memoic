package native

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"memoic/internal/memoic"
)

var client = resty.New()

var _ = memoic.AddFunction("http.request", memoic.Function{
	{
		Invoke: func(stack *memoic.Stack) (any, error) {
			params, err := stack.MappedParameters()
			if err != nil {
				return nil, err
			}
			request := client.R()
			method, err := memoic.GetFrom[string](params, "method", true)
			if err != nil {
				return nil, err
			}
			request.Method = *method

			link, err := memoic.GetFrom[string](params, "link", false)
			if err != nil {
				return nil, err
			}
			if link != nil {
				request.URL = *link
			}

			headers, err := memoic.GetFrom[map[string]any](params, "headers", false)
			if err != nil {
				return nil, err
			}
			if headers != nil {
				mkay := make(map[string]string)
				for key, value := range *headers {
					if text, ok := value.(string); ok {
						mkay[key] = text
						continue
					}
					bytes, err := json.Marshal(value)
					if err != nil {
						return nil, errors.Join(errors.New("failed to jsonify header "+key), err)
					}
					mkay[key] = string(bytes)
				}
				request.SetHeaders(mkay)
			}

			if body, ok := params["body"]; ok {
				request.SetBody(body)
			}

			response, err := request.Send()
			if err != nil {
				return nil, err
			}
			return string(response.Body()), nil
		},
	},
})
