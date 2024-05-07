package convert

import jsoniter "github.com/json-iterator/go"

// MapToStruct converts a map into a struct via Json Marshal/Unmarshal
func MapToStruct(convertMap map[string]interface{}, convertTo interface{}) error {
	marshal, err := jsoniter.Marshal(convertMap)

	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(
		marshal,
		convertTo)

	return err
}

// MapToStructEx converts a map into a struct via Json Marshal/Unmarshal
// and converts any child maps of type map[interface{}]interface{} to
// map[string]interface{}
func MapToStructEx(convertMap map[string]interface{}, convertTo interface{}) error {
	convertMap = ChildMapsToStringKey(convertMap)

	marshal, err := jsoniter.Marshal(convertMap)

	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(
		marshal,
		convertTo)

	return err
}

// ChildMapsToStringKey converts a map[interface{}]interface{} to
// map[string]interface{} recursively, assuming that the input map key values
// are explicitly of type string
func ChildMapsToStringKey(convertMap map[string]interface{}) map[string]interface{} {
	temp := map[string]interface{}{}

	for k, v := range convertMap {
		// recursively convert maps
		if childMap, ok := v.(map[interface{}]interface{}); ok {
			v = ToMapWithStringKey(childMap)
		}
		temp[k] = v
	}

	return temp
}

// ToMapWithStringKey converts a map[interface{}]interface{} to
// map[string]interface{} recursively, assuming that the input map key values
// are explicitly of type string
func ToMapWithStringKey(convertMap map[interface{}]interface{}) map[string]interface{} {
	temp := map[string]interface{}{}

	for k, v := range convertMap {
		// recursively convert maps
		if childMap, ok := v.(map[interface{}]interface{}); ok {
			v = ToMapWithStringKey(childMap)
		}
		temp[k.(string)] = v
	}

	return temp
}

// GenericMapToStruct converts a map into a struct via Json Marshal/Unmarshal
func GenericMapToStruct(convertMap map[interface{}]interface{}, thing interface{}) error {
	// convert to map[string]interface{} (recursively)
	return MapToStruct(ToMapWithStringKey(convertMap), thing)
}
