package query

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type OmitTestSubObj struct {
	SS string
}

type OmitTestObj struct {
	S   string          `query:"s,omitempty"`
	N   int             `query:",omitempty"`
	Obj *OmitTestSubObj `query:",omitempty"`
}

func TestMarshalOmit1(t *testing.T) {
	obj := &OmitTestObj{
		S: "",
		N: 10,
	}

	values := url.Values{}
	err := Marshal(obj, values)
	assert.Nil(t, err)
	t.Log(values)
}
