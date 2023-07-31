package memoic

import "fmt"

// Sector is a space where pipelines and others can store data to, this is simply
// a map.
type Sector map[string]any

// GlobalSector is a global space that can be accessed by all functions regardless of the runtime.
var GlobalSector = make(Sector)

// GetFrom gets the value of a key from the given sector and tries to cast it
// into T (the expected type),  if it cannot cast then it will error out, otherwise
// if the key is not there, it will either: (if must then error else ignore).
func GetFrom[T any](sector Sector, key string, must bool) (*T, error) {
	if cast, ok := sector[key]; ok {
		cast, ok := cast.(T)
		if !ok {
			return nil, fmt.Errorf("%s can only be of type %T, instead we got %T", key, *new(T), sector[key])
		}
		return &cast, nil
	} else {
		if must {
			return nil, fmt.Errorf("cannot find %s in memory", key)
		}
		return nil, nil
	}
}
