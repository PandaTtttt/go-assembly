package simplejson_test

import (
	"github.com/PandaTtttt/go-assembly/simplejson"
	"github.com/PandaTtttt/go-assembly/util/jsonutil"
	"github.com/PandaTtttt/go-assembly/util/m"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMerge(t *testing.T) {
	type _subFoo struct {
		String string `json:"string"`
		Int    int    `json:"int"`
		Ints   []int  `json:"ints"`
		Bool   bool   `json:"bool"`
	}

	type _foo struct {
		String          string  `json:"string"`
		UnchangedString string  `json:"unchanged_string"`
		privateString   string  // should not change during merge operation since private field is invisible to json marshal.
		Int             int     `json:"int"`
		Ints            []int   `json:"ints"`
		SubFoo          _subFoo `json:"sub_foo"`
	}

	foo := _foo{
		String:          "string",
		UnchangedString: "unchanged_string",
		privateString:   "private_string",
		Int:             10,
		Ints:            []int{1, 2, 3, 4, 5},
		SubFoo: _subFoo{
			String: "sub_string",
			Int:    20,
			Ints:   []int{1, 2},
		},
	}

	j := simplejson.NewFrom(m.M{
		"string":         "new_string",
		"ints":           []int{8, 9},
		"private_string": "new_private_string",
		"sub_foo": m.M{
			"string": "new_sub_string",
			"int":    30,
			"bool":   true,
		},
	})

	err := jsonutil.Merge(j, &foo)
	assert.Equal(t, nil, err)

	assert.Equal(t, "new_string", foo.String)
	assert.Equal(t, "unchanged_string", foo.UnchangedString)
	assert.Equal(t, "private_string", foo.privateString)
	assert.Equal(t, 10, foo.Int)
	assert.Equal(t, []int{8, 9}, foo.Ints)
	assert.Equal(t, "new_sub_string", foo.SubFoo.String)
	assert.Equal(t, 30, foo.SubFoo.Int)
	assert.Equal(t, []int{1, 2}, foo.SubFoo.Ints)
	assert.Equal(t, true, foo.SubFoo.Bool)
}
