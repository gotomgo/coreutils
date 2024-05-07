package convert

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gotomgo/coreutils/errors"
)

// ToInt32 converts a value to an int32
func ToInt32(input interface{}) (result int32, err error) {
	switch input := input.(type) {
	case string:
		var i64 int64
		if i64, err = strconv.ParseInt(input, 0, 32); err == nil {
			result = int32(i64)
		}
	case uint8:
		result = int32(input)
	case uint16:
		result = int32(input)
	case uint32:
		result = int32(input)
	case uint64:
		result = int32(input)
	case uint:
		result = int32(input)
	case int8:
		result = int32(input)
	case int16:
		result = int32(input)
	case int32:
		result = input
	case int:
		result = int32(input)
	case int64:
		result = int32(input)
	case float32:
		if math.IsInf(float64(input), 0) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert +Inf,-Inf to int64")
		} else if math.IsNaN(float64(input)) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert NaN to int64")
		} else {
			result = int32(input)
		}
	case float64:
		if math.IsInf(input, 0) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert +Inf,-Inf to int64")
		} else if math.IsNaN(input) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert NaN to int64")
		} else {
			result = int32(input)
		}
	default:
		err = fmt.Errorf("wrong type")
	}

	return
}

// ToInt converts a value to an int
func ToInt(input interface{}) (result int, err error) {
	switch input := input.(type) {
	case string:
		var i64 int64
		if i64, err = strconv.ParseInt(input, 0, 64); err == nil {
			result = int(i64)
		}
	case uint8:
		result = int(input)
	case uint16:
		result = int(input)
	case uint32:
		result = int(input)
	case uint64:
		result = int(input)
	case int8:
		result = int(input)
	case int16:
		result = int(input)
	case int32:
		result = int(input)
	case uint:
		result = int(input)
	case int:
		result = input
	case *int:
		result = *input
	case int64:
		result = int(input)
	case float32:
		if math.IsInf(float64(input), 0) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert +Inf,-Inf to int64")
		} else if math.IsNaN(float64(input)) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert NaN to int64")
		} else {
			result = int(input)
		}
	case float64:
		if math.IsInf(input, 0) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert +Inf,-Inf to int64")
		} else if math.IsNaN(input) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert NaN to int64")
		} else {
			result = int(input)
		}
	default:
		err = fmt.Errorf("cannot convert %T to int", input)
	}

	if err != nil {
		err = errors.ErrConversionFailed.Instance("to int").WithInner(err)
	}

	return
}

// ToInt64 converts a value to an int64
func ToInt64(input interface{}) (result int64, err error) {
	switch input := input.(type) {
	case string:
		result, err = strconv.ParseInt(input, 0, 64)
	case uint8:
		result = int64(input)
	case uint16:
		result = int64(input)
	case uint32:
		result = int64(input)
	case uint64:
		result = int64(input)
	case uint:
		result = int64(input)
	case int8:
		result = int64(input)
	case int16:
		result = int64(input)
	case int32:
		result = int64(input)
	case int:
		result = int64(input)
	case int64:
		result = input
	case *int64:
		result = *input
	case float32:
		if math.IsInf(float64(input), 0) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert +Inf,-Inf to int64")
		} else if math.IsNaN(float64(input)) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert NaN to int64")
		} else {
			result = int64(input)
		}
	case float64:
		if math.IsInf(input, 0) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert +Inf,-Inf to int64")
		} else if math.IsNaN(input) {
			err = errors.ErrArgumentOutOfRange.Instancef("cannot convert NaN to int64")
		} else {
			result = int64(input)
		}
	default:
		err = fmt.Errorf("cannot convert %T to int64", input)
	}

	if err != nil {
		err = errors.ErrConversionFailed.Instance("to int64").WithInner(err)
	}

	return
}

// ToFloat64 converts a value to a float64
func ToFloat64(input interface{}) (result float64, err error) {
	switch input := input.(type) {
	case string:
		result, err = strconv.ParseFloat(input, 64)
	case uint8:
		result = float64(input)
	case uint16:
		result = float64(input)
	case uint32:
		result = float64(input)
	case uint64:
		result = float64(input)
	case uint:
		result = float64(input)
	case int8:
		result = float64(input)
	case int16:
		result = float64(input)
	case int32:
		result = float64(input)
	case int:
		result = float64(input)
	case int64:
		result = float64(input)
	case float32:
		result = float64(input)
	case float64:
		result = input
	case *float64:
		result = *input
	default:
		err = fmt.Errorf("cannot convert %T to float64", input)
	}

	if err != nil {
		err = errors.ErrConversionFailed.Instance("to float64").WithInner(err)
	}

	return
}

func ToBool(input interface{}) (result bool, err error) {
	switch input := input.(type) {
	case string:
		result, err = strconv.ParseBool(input)
	case bool:
		result = input
	case *bool:
		result = *input
	case int:
		result = input != 0
	default:
		err = fmt.Errorf("the type %T could not be converted to bool", input)
	}

	if err != nil {
		err = errors.ErrConversionFailed.Instance("to bool").WithInner(err)
	}

	return
}

func ToString(input interface{}) (result string, err error) {
	switch input := input.(type) {
	case string:
		result = input
	case *string:
		if input != nil {
			result = *input
		}
	default:
		if stringer, ok := input.(fmt.Stringer); ok {
			result = stringer.String()
		} else {
			err = fmt.Errorf("the type %T could not be converted to string", input)
		}
	}

	if err != nil {
		err = errors.ErrConversionFailed.Instance("to string").WithInner(err)
	}

	return
}
