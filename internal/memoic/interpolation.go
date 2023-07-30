package memoic

import "encoding/json"

type marshaler func(v any) ([]byte, error)

var marshalers = map[string]marshaler{
	"json": json.Marshal,
}
