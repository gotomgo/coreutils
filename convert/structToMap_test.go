package convert

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type A struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type B struct {
	Pairs []A      `json:"pairs"`
	Types []string `json:"types"`
}

type C struct {
	KeyPairs map[string]interface{}
}

func TestIsArrayOfStruct(t *testing.T) {
	a := []A{}
	assert.True(t, isArrayOfStruct(a))
	assert.True(t, isArrayOfStruct(&a))

	b := []*A{}
	assert.True(t, isArrayOfStruct(b))
	assert.True(t, isArrayOfStruct(&b))

	c := []interface{}{}
	assert.True(t, isArrayOfStruct(c))
	assert.True(t, isArrayOfStruct(&c))

	d := []string{}
	assert.False(t, isArrayOfStruct(d))
	assert.False(t, isArrayOfStruct(&d))

	e := []int{1, 2, 3}
	assert.False(t, isArrayOfStruct(e))
	assert.False(t, isArrayOfStruct(&e))
}

func TestArrayOfStructToMap(t *testing.T) {
	a := []A{{Name: "test", ID: 1234}}
	assert.True(t, isArrayOfStruct(a))
	assert.True(t, isArrayOfStruct(&a))

	// use an anonymous struct to convert, and return the converted array
	convertMap := ToMap(struct {
		A []A `json:"a"`
	}{A: a})
	spew.Dump(convertMap["a"])
}

type CC struct {
	Name string
}

type DD struct {
	ID string `json:"id"`
	CC CC     `json:"cc"`
}

func (cc CC) ToMap(tagType string) map[string]interface{} {
	return map[string]interface{}{"alias_name": cc.Name}
}

func TestMapConverter(t *testing.T) {
	cc := CC{Name: "test converter"}
	spew.Dump(ToMap(&cc))

	dd := DD{ID: "1234", CC: cc}
	spew.Dump(ToMap(&dd))
}

func TestArrayOfStruct(t *testing.T) {
	a := A{Name: "test", ID: 1234}
	b := []A{a, A{Name: "test2", ID: 2345}}
	spew.Dump(ToArray(b))

	aa := []string{"a", "b", "c"}
	spew.Dump(ToArray(aa))
}

type FF struct {
	Value int64 `json:"value"`
}

type EE struct {
	A
	FF
}

func TestAnonymous(t *testing.T) {
	ee := EE{A: A{Name: "test"}, FF: FF{Value: 42}}
	spew.Dump(ToMap(&ee))
}

type GG struct {
	A

	ID  int64  `json:"id"`
	Str string `json:"str"`
}

func TestAnonymous2(t *testing.T) {
	ee := GG{A: A{Name: "test"}, ID: 42, Str: "testStr"}
	spew.Dump(ToMap(&ee))

}
