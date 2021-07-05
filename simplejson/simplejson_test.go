package simplejson

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewFromReader(t *testing.T) {
	//Use New Constructor
	buf := bytes.NewBuffer([]byte(`{
		"test": {
			"array": [1, "2", 3],
			"array_with_subs": [
				{"sub_key_one": 1},
				{"sub_key_two": 2, "sub_key_three": 3}
			],
			"big_num": 9223372036854775807,
			"uint64": 18446744073709551615
		}
	}`))
	js, err := NewFromReader(buf)

	//Standard Test Case
	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)

	arr := js.Get("test", "array").Array()
	assert.NotEqual(t, nil, arr)
	for i, v := range arr {
		var iv int
		switch v.(type) {
		case json.Number:
			i64, err := v.(json.Number).Int64()
			assert.Equal(t, nil, err)
			iv = int(i64)
		case string:
			iv, _ = strconv.Atoi(v.(string))
		}
		assert.Equal(t, i+1, iv)
	}

	ma := js.Get("test", "array").Array()
	assert.Equal(t, ma, []interface{}{json.Number("1"), "2", json.Number("3")})

	mm := js.Get("test", "array_with_subs", 0).Map()
	assert.Equal(t, mm, map[string]interface{}{"sub_key_one": json.Number("1")})

	assert.Equal(t, js.Get("test", "big_num").Int64(), int64(9223372036854775807))
	assert.Equal(t, js.Get("test", "uint64").Uint64(), uint64(18446744073709551615))
}

func TestSimplejson1(t *testing.T) {
	var ok bool
	var err error

	js, err := NewJSON([]byte(`{
		"test": {
			"string_array": ["foo", "bar", "baz"],
			"string_array_null": ["abc", null, "efg"],
			"array": [1, "2", 3],
			"array_with_subs": [{"sub_key_one": 1},
			{"sub_key_two": 2, "sub_key_three": 3}],
			"int": 10,
			"float": 5.150,
			"string": "simplejson",
			"bool": true,
			"sub_obj": {"a": 1}
		}
	}`))

	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)

	_, ok = js.CheckGet("test")
	assert.Equal(t, true, ok)

	_, ok = js.CheckGet("missing_key")
	assert.Equal(t, false, ok)

	aws := js.Get("test", "array_with_subs")
	assert.NotEqual(t, nil, aws)
	var awsVal int
	awsVal = aws.Get(0, "sub_key_one").Int()
	assert.Equal(t, 1, awsVal)
	awsVal = aws.Get(1, "sub_key_two").Int()
	assert.Equal(t, 2, awsVal)
	awsVal = aws.Get(1, "sub_key_three").Int()
	assert.Equal(t, 3, awsVal)

	i := js.Get("test", "int").Int()
	assert.Equal(t, 10, i)

	f := js.Get("test", "float").Float64()
	assert.Equal(t, 5.150, f)

	s := js.Get("test", "string").String()
	assert.Equal(t, "simplejson", s)

	b := js.Get("test", "bool").Bool()
	assert.Equal(t, true, b)

	mi2 := js.Get("test").Get("missing_int").Int(5150)
	assert.Equal(t, 5150, mi2)

	ms2 := js.Get("test").Get("missing_string").String("bing")
	assert.Equal(t, "bing", ms2)

	ma2 := js.Get("test", "missing_array").Array([]interface{}{"1", 2, "3"})
	assert.Equal(t, ma2, []interface{}{"1", 2, "3"})

	msa := js.Get("test", "string_array").Array()
	assert.Equal(t, msa[0], "foo")
	assert.Equal(t, msa[1], "bar")
	assert.Equal(t, msa[2], "baz")

	mm2 := js.Get("test").Get("missing_map").Map(map[string]interface{}{"found": false})
	assert.Equal(t, mm2, map[string]interface{}{"found": false})

	strs2 := js.Get("test").Get("string_array_null").Array()
	assert.Equal(t, strs2[0], "abc")
	assert.Equal(t, strs2[1], nil)
	assert.Equal(t, strs2[2], "efg")

	assert.Equal(t, js.Get("test").Get("bool").Bool(), true)

	js.Set("float2", 300.0)
	assert.Equal(t, js.Get("float2").Float64(), 300.0)

	js.Set("test2", "setTest")
	assert.Equal(t, "setTest", js.Get("test2").String())

	js.Del("test2")
	assert.NotEqual(t, "setTest", js.Get("test2").String())

	js.Get("test").Get("sub_obj").Set("a", 2)
	assert.Equal(t, 2, js.Get("test").Get("sub_obj").Get("a").Int())

	js.Get("test", "sub_obj").Set("a", 3)
	assert.Equal(t, 3, js.Get("test", "sub_obj", "a").Int())
}

func TestStdlibInterfaces(t *testing.T) {
	val := new(struct {
		Name   string `json:"name"`
		Params *JSON  `json:"params"`
	})
	val2 := new(struct {
		Name   string `json:"name"`
		Params *JSON  `json:"params"`
	})

	raw := `{"name":"my_object","params":{"string":"simplejson"}}`

	assert.Equal(t, nil, json.Unmarshal([]byte(raw), val))

	assert.Equal(t, "my_object", val.Name)
	assert.NotEqual(t, nil, val.Params.data)
	s := val.Params.Get("string").String()
	assert.Equal(t, "simplejson", s)

	p, err := json.Marshal(val)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, json.Unmarshal(p, val2))
	assert.Equal(t, val, val2) // stable
}

func TestSet(t *testing.T) {
	js, err := NewJSON([]byte(`{}`))
	assert.Equal(t, nil, err)

	js.Set("baz", "bing")

	s := js.Get("baz").String()
	assert.Equal(t, "bing", s)
}

func TestReplace(t *testing.T) {
	js, err := NewJSON([]byte(`{}`))
	assert.Equal(t, nil, err)

	err = js.UnmarshalJSON([]byte(`{"baz":"bing"}`))
	assert.Equal(t, nil, err)

	s := js.Get("baz").String()
	assert.Equal(t, "bing", s)
}

func TestSetPath(t *testing.T) {
	js, err := NewJSON([]byte(`{}`))
	assert.Equal(t, nil, err)

	js.SetPath([]string{"foo", "bar"}, "baz")

	s := js.Get("foo", "bar").String()
	assert.Equal(t, "baz", s)
}

func TestSetPathNoPath(t *testing.T) {
	js, err := NewJSON([]byte(`{"some":"data","some_number":1.0,"some_bool":false}`))
	assert.Equal(t, nil, err)

	f := js.Get("some_number").Float64(99.0)
	assert.Equal(t, f, 1.0)

	js.SetPath([]string{}, map[string]interface{}{"foo": "bar"})

	s := js.Get("foo").String()
	assert.Equal(t, "bar", s)

	f = js.Get("some_number").Float64(99.0)
	assert.Equal(t, f, 99.0)
}

func TestPathWillAugmentExisting(t *testing.T) {
	js, err := NewJSON([]byte(`{"this":{"a":"aa","b":"bb","c":"cc"}}`))
	assert.Equal(t, nil, err)

	js.SetPath([]string{"this", "d"}, "dd")

	cases := []struct {
		path    []interface{}
		outcome string
	}{
		{
			path:    []interface{}{"this", "a"},
			outcome: "aa",
		},
		{
			path:    []interface{}{"this", "b"},
			outcome: "bb",
		},
		{
			path:    []interface{}{"this", "c"},
			outcome: "cc",
		},
		{
			path:    []interface{}{"this", "d"},
			outcome: "dd",
		},
	}

	for _, tc := range cases {
		s, ok := js.Get(tc.path...).CheckString()
		assert.Equal(t, true, ok)
		assert.Equal(t, tc.outcome, s)
	}
}

func TestPathWillOverwriteExisting(t *testing.T) {
	// notice how "a" is 0.1 - but then we'll try to set at path a, foo
	js, err := NewJSON([]byte(`{"this":{"a":0.1,"b":"bb","c":"cc"}}`))
	assert.Equal(t, nil, err)

	js.SetPath([]string{"this", "a", "foo"}, "bar")

	s, ok := js.Get("this", "a", "foo").CheckString()
	assert.Equal(t, true, ok)
	assert.Equal(t, "bar", s)
}
