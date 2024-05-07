package convert

import (
	"github.com/gotomgo/coreutils/errors"
	jsoniter "github.com/json-iterator/go"
)

// SimpleStructToMap Converts a struct to a map via JSON Marshalling
func SimpleStructToMap(something interface{}) (map[string]interface{}, error) {
	var aMap map[string]interface{}

	marshal, err := jsoniter.Marshal(something)

	if err != nil {
		return nil, err
	}

	if err = jsoniter.Unmarshal(
		marshal,
		&aMap); errors.IsError(err) {
		return nil, err
	}

	return aMap, nil
}
