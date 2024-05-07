package convert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToInt(t *testing.T) {
	t.Run("ToInt success", func(t *testing.T) {
		t.Parallel()

		val, err := ToInt(uint8(64))
		if !assert.NoError(t, err, "int64 from uint8 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from uint8 should work") {
			return
		}

		val, err = ToInt(uint16(64))
		if !assert.NoError(t, err, "int64 from uint16 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from uint16 should work") {
			return
		}

		val, err = ToInt(uint32(64))
		if !assert.NoError(t, err, "int64 from uint32 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from uint32 should work") {
			return
		}

		val, err = ToInt(uint(64))
		if !assert.NoError(t, err, "int64 from uint64 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from uint64 should work") {
			return
		}

		val, err = ToInt(int8(64))
		if !assert.NoError(t, err, "int64 from int8 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from int8 should work") {
			return
		}

		val, err = ToInt(int16(64))
		if !assert.NoError(t, err, "int64 from int16 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from int16 should work") {
			return
		}

		val, err = ToInt(int32(64))
		if !assert.NoError(t, err, "int64 from int32 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from int32 should work") {
			return
		}

		val, err = ToInt(int(64))
		if !assert.NoError(t, err, "int64 from int should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from int should work") {
			return
		}

		val, err = ToInt(int(64))
		if !assert.NoError(t, err, "int64 from int64 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from int64 should work") {
			return
		}

		temp := int(64)
		val, err = ToInt(&temp)
		if !assert.NoError(t, err, "int64 from *int64 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from *int64 should work") {
			return
		}

		val, err = ToInt(float32(64.0))
		if !assert.NoError(t, err, "int64 from float32 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from float32 should work") {
			return
		}

		val, err = ToInt(64.0)
		if !assert.NoError(t, err, "int64 from float64 should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from float64 should work") {
			return
		}

		val, err = ToInt("64")
		if !assert.NoError(t, err, "int64 from string should work") {
			return
		}
		if !assert.Equal(t, int(64), val, "int64 from string should work") {
			return
		}
	})

	t.Run("ToInt uint boundary success", func(t *testing.T) {
		t.Parallel()

		val, err := ToInt(uint8(math.MaxUint8))
		if !assert.NoError(t, err, "int64 from uint8 should work") {
			return
		}
		if !assert.Equal(t, int(math.MaxUint8), val, "int64 from uint8 should work") {
			return
		}

		val, err = ToInt(uint16(math.MaxUint16))
		if !assert.NoError(t, err, "int64 from uint16 should work") {
			return
		}
		if !assert.Equal(t, int(math.MaxUint16), val, "int64 from uint16 should work") {
			return
		}

		val, err = ToInt(uint32(math.MaxUint32))
		if !assert.NoError(t, err, "int64 from uint32 should work") {
			return
		}
		if !assert.Equal(t, int(math.MaxUint32), val, "int64 from uint32 should work") {
			return
		}

		// sign bit results in -1 so we have to check that otherwise compiler
		// detects overflow with const (but not variable)
		val, err = ToInt(uint64(math.MaxUint64))
		if !assert.NoError(t, err, "int64 from uint64 should work") {
			return
		}
		if !assert.Equal(t, int(-1), int(val), "int64 from uint64 should work") {
			return
		}

		// Test w/ variable and prove no overflow error
		temp := uint64(math.MaxUint64)
		val, err = ToInt(temp)
		if !assert.NoError(t, err, "int64 from uint64 should work") {
			return
		}
		if !assert.Equal(t, int(temp), int(val), "int64 from uint64 should work") {
			return
		}

		// avoid compiler const checks by using a variable
		tempf := math.MaxFloat64
		val, err = ToInt(math.MaxFloat64)
		if !assert.NoError(t, err, "int64 from float64 should work") {
			return
		}
		if !assert.Equal(t, int(tempf), int(val), "int64 from float64 should work") {
			return
		}
	})

	t.Run("ToInt float invalid fails", func(t *testing.T) {
		t.Parallel()

		_, err := ToInt(math.Inf(+1))
		if !assert.Error(t, err, "int64 from +Inf should fail") {
			return
		}

		_, err = ToInt(math.Inf(-1))
		if !assert.Error(t, err, "int64 from -Inf should fail") {
			return
		}

		_, err = ToInt(math.NaN())
		if !assert.Error(t, err, "int64 from NaN should fail") {
			return
		}
	})

	t.Run("ToInt non-supported type fails", func(t *testing.T) {
		t.Parallel()

		temp := int64(64)
		_, err := ToInt(&temp)
		if !assert.Error(t, err, "int from *int64 should fail") {
			return
		}
	})

	t.Run("ToInt more strings pass", func(t *testing.T) {
		t.Parallel()

		val, err := ToInt("0xDEADBEEF")
		if !assert.NoError(t, err, "int64 from hex string should pass") {
			return
		}
		if !assert.Equal(t, int(0xDEADBEEF), val, "int64 from hex string should work") {
			return
		}

		val, err = ToInt("0777")
		if !assert.NoError(t, err, "int64 from octal string should pass") {
			return
		}
		if !assert.Equal(t, int(0777), val, "int64 octal hex string should work") {
			return
		}
	})

	t.Run("ToInt bad string fails", func(t *testing.T) {
		t.Parallel()

		_, err := ToInt("")
		if !assert.Error(t, err, "int64 from empty string should fail") {
			return
		}

		_, err = ToInt("14Gf32")
		if !assert.Error(t, err, "int64 from non-numeric string should fail") {
			return
		}
	})

}
