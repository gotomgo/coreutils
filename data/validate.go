//go:build !RELEASE
// +build !RELEASE

package data

import (
	"github.com/gotomgo/coreutils/errors"
)

// Validate validates an item if it is non-nil and implements IValidationDeep
// or IValidation
func Validate(item interface{}) (err error) {
	if !isNilFixed(item) {
		if val, ok := item.(IValidationDeep); ok {
			err = val.ValidateDeep()
		} else if val, ok := item.(IValidation); ok {
			err = val.Validate()
		} else {
			method := FindMethod(item, "ValidateDeep")
			if method == nil {
				method = FindMethod(item, "Validate")
			}

			if method != nil {
				_err := CallMethodWithItem(item, *method)
				if _err != nil {
					err, _ = _err.(error)
				}
			} else {
				// log.Warningf("not validating type %T", item)
			}
		}
	}

	return
}

// ValidateLight validates an item if it is non-nil and implements IValidation
func ValidateLight(item interface{}) (err error) {
	if !isNilFixed(item) {
		if val, ok := item.(IValidation); ok {
			err = val.Validate()
		}
	}

	return
}

// ValidateNotNil performs a Validate on an item that cannot be nil
func ValidateNotNil(item interface{}) (err error) {
	if isNilFixed(item) {
		err = errors.ErrArgumentNull.Instance("item")
	} else {
		err = Validate(item)
	}

	return
}

// ValidateNotNilLight performs a Validate on an item that cannot be nil
func ValidateNotNilLight(item interface{}) (err error) {
	if isNilFixed(item) {
		err = errors.ErrArgumentNull.Instance("item")
	} else {
		if val, ok := item.(IValidation); ok {
			err = val.Validate()
		}
	}

	return
}
