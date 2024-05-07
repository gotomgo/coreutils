package data

import (
	"reflect"

	"github.com/gotomgo/coreutils/errors"
)

// IValidation is implemented by types that provide validation
type IValidation interface {
	Validate() error
}

// IValidationDeep is implemented by types to provide a deeper, more
// introspective (and costly) validation
type IValidationDeep interface {
	IValidation

	ValidateDeep() error
}

//	--------------------------------------------------------------------------
//	The following methods do not differ when the RELEASE tag is/isn't used
//	so they are defined here to avoid redundant definitions
//	--------------------------------------------------------------------------

func isNilFixed(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		//use of IsNil method
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func FindMethod(i interface{}, methodName string) (result *reflect.Value) {
	var ptr reflect.Value
	var value reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		result = &method
	} else {
		// check for method on pointer
		method = ptr.MethodByName(methodName)
		if method.IsValid() {
			result = &method
		}
	}

	return
}

func CallMethodWithItem(i interface{}, method reflect.Value) (result interface{}) {
	if method.IsValid() {
		result = method.Call([]reflect.Value{})[0].Interface()
	} else {
		// log.Errorf("CallMethodWithItem method is invalid")
	}

	return
}

func CallMethod(i interface{}, methodName string) (result interface{}) {
	method := FindMethod(i, methodName)
	if method != nil {
		result = CallMethodWithItem(i, *method)
	}

	return
}

// MustValidate panics if item fails Validate
func MustValidate(item interface{}) {
	if err := Validate(item); errors.IsError(err) {
		panic(err)
	}
}

// MustValidateNotNil panics if item is nil or fails Validate
func MustValidateNotNil(item interface{}) {
	if err := ValidateNotNil(item); errors.IsError(err) {
		panic(err)
	}
}
