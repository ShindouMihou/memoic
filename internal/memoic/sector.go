package memoic

import (
	"fmt"
)

func SectorGet[T any](sector Sector, key string, must bool) (*T, error) {
	if cast, ok := sector[key]; ok {
		cast, ok := cast.(T)
		if !ok {
			return nil, fmt.Errorf("%s can only be of type %T, instead we got %T", key, *new(T), sector[key])
		}
		return &cast, nil
	} else {
		if must {
			return nil, fmt.Errorf("cannot find %s in parameters", key)
		}
		return nil, nil
	}
}
