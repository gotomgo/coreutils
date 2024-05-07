package convert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToFloat64(t *testing.T) {
	t.Run("ToFloat64 success", func(t *testing.T) {
		t.Parallel()

		val, err := ToFloat64(uint8(64))
		if !assert.NoError(t, err, "float64 from uint8 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from uint8 should work") {
			return
		}

		val, err = ToFloat64(uint16(64))
		if !assert.NoError(t, err, "float64 from uint16 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from uint16 should work") {
			return
		}

		val, err = ToFloat64(uint32(64))
		if !assert.NoError(t, err, "float64 from uint32 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from uint32 should work") {
			return
		}

		val, err = ToFloat64(uint64(64))
		if !assert.NoError(t, err, "float64 from uint64 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from uint64 should work") {
			return
		}

		val, err = ToFloat64(int8(64))
		if !assert.NoError(t, err, "float64 from int8 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from int8 should work") {
			return
		}

		val, err = ToFloat64(int16(64))
		if !assert.NoError(t, err, "float64 from int16 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from int16 should work") {
			return
		}

		val, err = ToFloat64(int32(64))
		if !assert.NoError(t, err, "float64 from int32 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from int32 should work") {
			return
		}

		val, err = ToFloat64(int(64))
		if !assert.NoError(t, err, "float64 from int should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from int should work") {
			return
		}

		val, err = ToFloat64(float64(64))
		if !assert.NoError(t, err, "float64 from int64 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from int64 should work") {
			return
		}

		temp := float64(64)
		val, err = ToFloat64(&temp)
		if !assert.NoError(t, err, "float64 from *float64 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from *float64 should work") {
			return
		}

		val, err = ToFloat64(float32(64.0))
		if !assert.NoError(t, err, "float64 from float32 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from float32 should work") {
			return
		}

		val, err = ToFloat64(64.0)
		if !assert.NoError(t, err, "float64 from float64 should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from float64 should work") {
			return
		}

		val, err = ToFloat64("64")
		if !assert.NoError(t, err, "float64 from string should work") {
			return
		}
		if !assert.Equal(t, float64(64), val, "float64 from string should work") {
			return
		}
	})

	t.Run("ToFloat64 uint boundary success", func(t *testing.T) {
		t.Parallel()

		val, err := ToFloat64(uint8(math.MaxUint8))
		if !assert.NoError(t, err, "float64 from uint8 should work") {
			return
		}
		if !assert.Equal(t, float64(math.MaxUint8), val, "float64 from uint8 should work") {
			return
		}

		val, err = ToFloat64(uint16(math.MaxUint16))
		if !assert.NoError(t, err, "float64 from uint16 should work") {
			return
		}
		if !assert.Equal(t, float64(math.MaxUint16), val, "float64 from uint16 should work") {
			return
		}

		val, err = ToFloat64(uint32(math.MaxUint32))
		if !assert.NoError(t, err, "float64 from uint32 should work") {
			return
		}
		if !assert.Equal(t, float64(math.MaxUint32), val, "float64 from uint32 should work") {
			return
		}

		// sign bit results in -1 so we have to check that otherwise compiler
		// detects overflow with const (but not variable)
		val, err = ToFloat64(uint64(math.MaxUint64))
		if !assert.NoError(t, err, "float64 from uint64 should work") {
			return
		}
		/*if !assert.Equal(t, float64(-1), float64(val), "float64 from uint64 should work") {
			return
		}*/

		// Test w/ variable and prove no overflow error
		temp := uint64(math.MaxUint64)
		val, err = ToFloat64(temp)
		if !assert.NoError(t, err, "float64 from uint64 should work") {
			return
		}
		/*if !assert.Equal(t, float64(temp), float64(val), "float64 from uint64 should work") {
			return
		}*/

		// avoid compiler const checks by using a variable
		val, err = ToFloat64(math.MaxFloat64)
		if !assert.NoError(t, err, "float64 from float64 should work") {
			return
		}
		if !assert.Equal(t, math.MaxFloat64, float64(val), "float64 from float64 should work") {
			return
		}
	})

	t.Run("ToFloat64 non-supported type fails", func(t *testing.T) {
		t.Parallel()

		temp := 64
		_, err := ToFloat64(&temp)
		if !assert.Error(t, err, "float64 from *int should fail") {
			return
		}
	})

	t.Run("ToFloat64 more strings pass", func(t *testing.T) {
		t.Parallel()

		val, err := ToFloat64("1e6")
		if !assert.NoError(t, err, "float64 from hex string should pass") {
			return
		}
		if !assert.Equal(t, 1e6, val, "float64 from exp notation should work") {
			return
		}

	})

	t.Run("ToFloat64 bad string fails", func(t *testing.T) {
		t.Parallel()

		_, err := ToFloat64("")
		if !assert.Error(t, err, "float64 from empty string should fail") {
			return
		}

		_, err = ToFloat64("14Gf32")
		if !assert.Error(t, err, "float64 from non-numeric string should fail") {
			return
		}
	})

}
