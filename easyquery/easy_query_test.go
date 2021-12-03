package easyquery

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	vals := url.Values{}
	vals.Add("N", "10")
	vals.Add("Ss", "1")
	vals.Add("Ss", "2")
	vals.Add("Bs", "true")
	vals.Add("int_8_s", "1")
	vals.Add("int_8_s", "2")
	vals.Add("Int64s", "3")
	vals.Add("Int64s", "4")

	s := EncodeValues(vals)
	t.Log(s)
	assert.Equal(t, "Bs=true&Int64s=3&Int64s=4&N=10&Ss=1&Ss=2&int_8_s=1&int_8_s=2", s)

	vals2, err := ParseQuery(s)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(vals, vals2))
}

type AI interface {
}

type A struct {
	N          int
	N2         int
	ND         int `query:"nd,99"`
	S          string
	Ss         []string
	Bs         []bool
	Int8s      []int8
	Int64s     []int64
	AI         AI
	M          map[string]int
	MyJSONOk   float64
	MyJSANFall float32
}

type B struct {
	A *A `param:"A"`
}

func Test2(t *testing.T) {
	vals := url.Values{}
	vals.Add("N", "10")
	vals.Add("Ss", "1")
	vals.Add("Ss", "2")
	vals.Add("Bs", "true")
	vals.Add("int_8_s", "1")
	vals.Add("int_8_s", "2")
	vals.Add("Int64s", "3")
	vals.Add("Int64s", "4")

	b := &B{}
	err := Values2Orm(vals, b)
	assert.Nil(t, err)

	vals2 := url.Values{}
	err = Orm2Values(b, vals2)
	assert.Nil(t, err)

	b2 := &B{}
	err = Values2Orm(vals2, b2)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(b, b2))
}
