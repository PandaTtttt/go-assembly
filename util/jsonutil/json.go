package jsonutil

import (
	"encoding/json"
	"github.com/PandaTtttt/go-assembly/simplejson"
)

// Merge performs a partial update and apply it to "to".
// We require the "from"" to be a simplejson so what key is updated is expressed more precisely.
func Merge(from *simplejson.JSON, to interface{}) error {
	toV := simplejson.New()
	if err := Remarshal(to, toV); err != nil {
		return err
	}
	toV.MergeFrom(from)
	return Remarshal(toV, to)
}

// Remarshal converts one struct to another, assuming they share the same json structure.
func Remarshal(from interface{}, to interface{}) error {
	encoded, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, to)
}
