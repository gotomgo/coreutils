package convert

import (
	"encoding/base64"
	jsoniter "github.com/json-iterator/go"
)

// EncodeStruct converts a struct to JSON and Base64 encodes it
func EncodeStruct(obj interface{}) (string, error) {
	bytes, err := jsoniter.Marshal(obj)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

// DecodeStruct decodes a base64 encoded JSON struct representation to an
// instance of T
//
//	Notes
//		template should be of type T
//		template should be a pointer to a struct that is T (but not strictly
//		*T, which is why it is typed as interface{})
func DecodeStruct[T any](base64str string, template interface{}) (T, error) {
	bytes, err := base64.URLEncoding.DecodeString(base64str)
	if err != nil {
		var def T
		return def, err
	}

	err = jsoniter.Unmarshal(bytes, template)
	if err != nil {
		var def T
		return def, err
	}

	return template.(T), nil
}
