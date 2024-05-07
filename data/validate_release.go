//go:build RELEASE
// +build RELEASE

package data

import (
	"github.com/gotomgo/coreutils/errors"
)

// Validate validates an item if it is non-nil and implements IValidation
func Validate(item interface{}) (err error) {
	if !isNilFixed(item) {
		if val, ok := item.(IValidation); ok {
			err = val.Validate()
		}
	}

	return
}

// ValidateLight validates an item if it is non-nil and implements IValidation
func ValidateLight(item interface{}) error {
	return nil
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
func ValidateRequiredOpt(item interface{}) (err error) {
	if isNilFixed(item) {
		err = errors.ErrArgumentNull.Instance("item")
	}

	return
}
