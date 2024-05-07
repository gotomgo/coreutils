// +build !RELEASE

package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type vtest struct {
	err error
}

func (v *vtest) Validate() error {
	if v.err != nil {
		return v.err
	}

	return nil
}

type vtestDeep struct {
	err error
}

func (v *vtestDeep) Validate() error {
	return nil
}

func (v *vtestDeep) ValidateDeep() error {
	if v.err != nil {
		return v.err
	}

	return nil
}

func Test_Validate(t *testing.T) {
	t.Run("data.Validate() pass", func(t *testing.T) {
		t.Parallel()

		v := &vtest{}

		if !assert.NoError(t, Validate(v), "vtest.Validate() should pass") {
			return
		}
	})

	t.Run("data.Validate() fail", func(t *testing.T) {
		t.Parallel()

		v := &vtest{err: fmt.Errorf("test")}

		if !assert.Error(t, Validate(v), "vtest.Validate() should fail") {
			return
		}
	})

	t.Run("data.Validate() pass deep", func(t *testing.T) {
		t.Parallel()

		v := &vtestDeep{}

		if !assert.NoError(t, Validate(v), "vtestDeep.Validate() should pass") {
			return
		}
	})

	t.Run("data.Validate() fail deep", func(t *testing.T) {
		t.Parallel()

		v := &vtestDeep{err: fmt.Errorf("test")}

		if !assert.Error(t, Validate(v), "vtestDeep.Validate() should fail") {
			return
		}
	})

	t.Run("data.Validate() pass not * receiver", func(t *testing.T) {
		t.Parallel()

		v := vtestDeep{err: fmt.Errorf("test")}

		if !assert.Error(t, Validate(v), "vtestDeep.Validate() should fail") {
			return
		}
	})

}

func Test_ValidateLight(t *testing.T) {
	t.Run("data.ValidateLight() pass", func(t *testing.T) {
		t.Parallel()

		v := &vtest{}

		if !assert.NoError(t, ValidateLight(v), "vtest.Validate() should pass") {
			return
		}
	})

	t.Run("data.ValidateLight() fail", func(t *testing.T) {
		t.Parallel()

		v := &vtest{err: fmt.Errorf("test")}

		if !assert.Error(t, ValidateLight(v), "vtest.Validate() should fail") {
			return
		}
	})

	t.Run("data.ValidateLight() pass deep", func(t *testing.T) {
		t.Parallel()

		v := &vtestDeep{}

		if !assert.NoError(t, ValidateLight(v), "vtestDeep.Validate() should pass") {
			return
		}
	})

	// this will pass (by design) because it only fails if deep validation
	// occurs which ValidateLight does not do
	t.Run("data.ValidateLight() pass deep #2", func(t *testing.T) {
		t.Parallel()

		v := &vtestDeep{err: fmt.Errorf("test")}

		if !assert.NoError(t, ValidateLight(v), "vtestDeep.Validate() should pass") {
			return
		}
	})
}

func Test_ValidateNotNil(t *testing.T) {
	t.Run("data.ValidateNotNil() pass not-nil", func(t *testing.T) {
		t.Parallel()

		if !assert.NoError(t, ValidateNotNil(&vtest{}), "vtest.ValidateNotNil() should pass") {
			return
		}
	})

	t.Run("data.ValidateNotNil() fail nil", func(t *testing.T) {
		t.Parallel()

		if !assert.Error(t, ValidateNotNil(nil), "vtest.ValidateNotNil() should fail") {
			return
		}
	})
}

func Test_ValidateNotNilLight(t *testing.T) {
	t.Run("data.ValidateNotNil() pass not-nil", func(t *testing.T) {
		t.Parallel()

		if !assert.NoError(t, ValidateNotNilLight(&vtestDeep{err: fmt.Errorf("test")}), "vtestDeep.ValidateNotNil() should pass") {
			return
		}
	})

	t.Run("data.ValidateNotNil() fail nil", func(t *testing.T) {
		t.Parallel()

		if !assert.Error(t, ValidateNotNilLight(nil), "vtest.ValidateNotNil() should fail") {
			return
		}
	})
}

func Test_MustValidate(t *testing.T) {
	t.Run("data.MustValidate() pass", func(t *testing.T) {
		t.Parallel()

		v := &vtest{}

		if !assert.NotPanics(t, func() { MustValidate(v) }, "MustValidate() should pass") {
			return
		}
	})

	t.Run("data.Validate() fail", func(t *testing.T) {
		t.Parallel()

		v := &vtest{err: fmt.Errorf("test")}

		if !assert.Panics(t, func() { MustValidate(v) }, "MustValidate() should pass") {
			return
		}
	})
}

func Test_MustValidateNotNil(t *testing.T) {
	t.Run("data.MustValidateNotNil() pass", func(t *testing.T) {
		t.Parallel()

		v := &vtest{}

		if !assert.NotPanics(t, func() { MustValidateNotNil(v) }, "MustValidateNotNil() should pass") {
			return
		}
	})

	t.Run("data.MustValidateNotNil() fail", func(t *testing.T) {
		t.Parallel()

		if !assert.Panics(t, func() { MustValidateNotNil(nil) }, "MustValidateNotNil() should pass") {
			return
		}
	})
}
