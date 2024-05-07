package convert

// I got this from some random project, and extended it make it "better"

import (
	"fmt"
	"strings"

	"reflect"
)

// MapConverter is an interface that allows a struct type to specify a custom
// conversion of the struct instance to a map representation
type MapConverter interface {
	ToMap(tagType string) map[string]interface{}
}

func ToMap(s interface{}) (result map[string]interface{}) {
	return ToMapWithTag(s, "json")
}

func ToMapWithTag(s interface{}, tagType string) map[string]interface{} {
	if s == nil {
		return nil
	}

	return newStructWithTag(s, tagType).ToMap()
}

func ToArray(s interface{}) []interface{} {
	return ToArrayWithTag(s, "json")
}

func ToArrayWithTag(s interface{}, tagType string) (result []interface{}) {
	if !isArray(s) {
		panic("value is not []")
	}

	str := structVal{TagName: tagType}

	if isArrayOfStruct(s) {
		result = str.nested(reflect.ValueOf(s)).([]interface{})
	} else {
		result = str.makeSlice(reflect.ValueOf(s)).([]interface{})
	}

	return
}

// Struct encapsulates a struct type to provide several high level functions
// around the struct.
type structVal struct {
	raw     interface{}
	value   reflect.Value
	TagName string
}

// tagOptions contains a slice of tag options
type tagOptions []string

// Has returns true if the given option is available in tagOptions
func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

// parseTag splits a struct field's tag into its name and a list of options
// which comes after a name. A tag is in the form of: "name,option1,option2".
// The name can be neglected.
func parseTag(tag string) (string, tagOptions) {
	// tag is one of followings:
	// ""
	// "name"
	// "name,opt"
	// "name,opt,opt2"
	// ",opt"

	res := strings.Split(tag, ",")
	return res[0], res[1:]
}

// New returns a new *structVal with the struct s. It panics if the s's kind is
// not struct.
func newStruct(s interface{}) *structVal {
	return newStructWithTag(s, "json")
}

// New returns a new *structVal with the struct s. It panics if the s's kind is
// not struct.
func newStructWithTag(s interface{}, tagName string) *structVal {
	return &structVal{
		raw:     s,
		value:   getStructVal(s),
		TagName: tagName,
	}
}

// ToMap converts the given struct to a map[string]interface{}, where the keys
// of the map are the field names and the values of the map the associated
// values of the fields. The default key string is the struct field name but
// can be changed in the struct field's tag value. The "structs" key in the
// struct's field tag value is the key name. Example:
//
//	// Field appears in map as key "myName".
//	Name string `structs:"myName"`
//
// A tag value with the content of "-" ignores that particular field. Example:
//
//	// Field is ignored by this package.
//	Field bool `structs:"-"`
//
// A tag value with the content of "string" uses the stringer to get the value. Example:
//
//	// The value will be output of Animal's String() func.
//	// Map will panic if Animal does not implement String().
//	Field *Animal `structs:"field,string"`
//
// A tag value with the option of "flatten" used in a struct field is to flatten its fields
// in the output map. Example:
//
//	// The FieldStruct's fields will be flattened into the output map.
//	FieldStruct time.Time `structs:",flatten"`
//
// A tag value with the option of "omitnested" stops iterating further if the type
// is a struct. Example:
//
//	// Field is not processed further by this package.
//	Field time.Time     `structs:"myName,omitnested"`
//	Field *http.Request `structs:",omitnested"`
//
// A tag value with the option of "omitempty" ignores that particular field if
// the field value is empty. Example:
//
//	// Field appears in map as key "myName", but the field is
//	// skipped if empty.
//	Field string `structs:"myName,omitempty"`
//
//	// Field appears in map as key "Field" (the default), but
//	// the field is skipped if empty.
//	Field string `structs:",omitempty"`
//
// Note that only exported fields of a struct can be accessed, non exported
// fields will be neglected.
func (s *structVal) ToMap() map[string]interface{} {
	if converter, ok := s.raw.(MapConverter); ok {
		return converter.ToMap(s.TagName)
	}

	out := make(map[string]interface{})
	s.IntoMap(out)
	return out
}

// IntoMap is the same as Map. Instead of returning the output, it fills the
// given map.
func (s *structVal) IntoMap(out map[string]interface{}) {
	if out == nil {
		return
	}

	fields := s.structFields()

	for _, field := range fields {
		name := field.Name
		val := s.value.FieldByName(name)
		isSubStruct := false
		forceFlatten := false
		var finalVal interface{}

		tagName, tagOpts := parseTag(field.Tag.Get(s.TagName))
		if tagName != "" {
			name = tagName
		}

		// if the value is a zero value and the field is marked as omitempty do
		// not include
		if tagOpts.Has("omitempty") {
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()

			if reflect.DeepEqual(current, zero) {
				continue
			}
		}

		if !tagOpts.Has("omitnested") {
			if converter, ok := val.Interface().(MapConverter); ok {
				finalVal = converter.ToMap(s.TagName)
			} else {

				finalVal = s.nested(val)

				v := reflect.ValueOf(val.Interface())
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}

				switch v.Kind() {
				case reflect.Map, reflect.Struct:
					isSubStruct = true
					forceFlatten = field.Anonymous && (len(tagName) == 0)
				}
			}
		} else {
			finalVal = val.Interface()
		}

		if tagOpts.Has("string") {
			s, ok := val.Interface().(fmt.Stringer)
			if ok {
				out[name] = s.String()
			}
			continue
		}

		if isSubStruct && (forceFlatten || tagOpts.Has("flatten")) {
			for k := range finalVal.(map[string]interface{}) {
				out[k] = finalVal.(map[string]interface{})[k]
			}
		} else {
			out[name] = finalVal
		}
	}
}

// structFields returns the exported struct fields for a given s struct. This
// is a convenient helper method to avoid duplicate code in some of the
// functions.
func (s *structVal) structFields() []reflect.StructField {
	t := s.value.Type()

	var f []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}

		// don't check if it's omitted
		if tag := field.Tag.Get(s.TagName); tag == "-" {
			continue
		}

		f = append(f, field)
	}

	return f
}

func getStructVal(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("the value specified is not a struct")
	}

	return v
}

func isArray(s interface{}) (result bool) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// uninitialized zero value of a struct
	if v.Kind() == reflect.Invalid {
		return
	}

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		result = true
	}

	return
}

func isArrayOfStruct(s interface{}) (result bool) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// uninitialized zero value of a struct
	if v.Kind() == reflect.Invalid {
		return
	}

	if v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Array {

		if v.Type().Elem().Kind() == reflect.Interface {
			result = true
		} else if (v.Type().Elem().Kind() == reflect.Struct) ||
			((v.Type().Elem().Kind() == reflect.Ptr) &&
				(v.Type().Elem().Elem().Kind() == reflect.Struct)) {
			result = true
		}
	}

	return
}

// isStruct returns true if the given variable is a struct or a pointer to
// struct or a **struct.
func isStruct(s interface{}) bool {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// uninitialized zero value of a struct
	if v.Kind() == reflect.Invalid {
		return false
	}

	return v.Kind() == reflect.Struct
}

func (s *structVal) makeSlice(val reflect.Value) interface{} {
	slice := make([]interface{}, val.Len())
	for x := 0; x < val.Len(); x++ {
		slice[x] = s.nested(val.Index(x))
	}
	return slice
}

// nested retrieves recursively all types for the given value and returns the
// nested value.
func (s *structVal) nested(val reflect.Value) interface{} {
	var finalVal interface{}

	if converter, ok := val.Interface().(MapConverter); ok {
		return converter.ToMap(s.TagName)
	}

	v := reflect.ValueOf(val.Interface())
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		n := newStructWithTag(val.Interface(), s.TagName)
		m := n.ToMap()

		// do not add the converted value if there are no exported fields, ie:
		// time.Time
		// BUG: This causes problems when a type has all fields are omitempty,
		// and the map ends up being empty
		if len(m) == 0 {
			finalVal = val.Interface()
		} else {
			finalVal = m
		}
	case reflect.Map:
		// get the element type of the map
		mapElem := val.Type()
		switch val.Type().Kind() {
		case reflect.Ptr, reflect.Array, reflect.Map,
			reflect.Slice, reflect.Chan:
			mapElem = val.Type().Elem()
			if mapElem.Kind() == reflect.Ptr {
				mapElem = mapElem.Elem()
			}
		}

		// only iterate over struct types, ie: map[string]StructType,
		// map[string][]StructType,
		if (mapElem.Kind() == reflect.Struct || mapElem.Kind() == reflect.Interface) ||
			(mapElem.Kind() == reflect.Slice &&
				(mapElem.Elem().Kind() == reflect.Struct || mapElem.Kind() == reflect.Interface)) {
			m := make(map[string]interface{}, val.Len())
			for _, k := range val.MapKeys() {
				m[k.String()] = s.nested(reflect.ValueOf(val.MapIndex(k).Interface()))
			}
			finalVal = m
			break
		}

		finalVal = val.Interface()
	case reflect.Slice, reflect.Array:
		if val.Type().Kind() == reflect.Interface {
			finalVal = val.Interface()
			break
		}

		if val.Type().Elem().Kind() == reflect.Interface {
			var hasStruct bool
			for x := 0; x < val.Len(); x++ {
				if isStruct(val.Index(x).Interface()) {
					hasStruct = true
					break
				}
			}
			if hasStruct {
				finalVal = s.makeSlice(val)
				break
			}

			// use array as is
			finalVal = val.Interface()
			break
		}

		// do not iterate of non struct types, just pass the value. Ie: []int,
		// []string, co... We only iterate further if it's a struct.
		// i.e []foo or []*foo
		if val.Type().Elem().Kind() != reflect.Struct &&
			!(val.Type().Elem().Kind() == reflect.Ptr &&
				val.Type().Elem().Elem().Kind() == reflect.Struct) {
			finalVal = val.Interface()
			break
		}

		finalVal = s.makeSlice(val)
	default:
		finalVal = val.Interface()
	}

	return finalVal
}
