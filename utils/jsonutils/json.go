package jsonutils

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

// IsJSONObject checks if the passed data is a JSON Object
func IsJSONObject(data []byte) bool {
	var d = struct{}{}
	return json.Unmarshal(data, &d) == nil
}

// IsJSONArray checks if the passed data is a JSON Array
func IsJSONArray(data []byte) bool {
	d := make([]interface{}, 0)
	return json.Unmarshal(data, &d) == nil
}

var emptyArray = []byte{
	'[',
	']',
}

func MergeJSONObjects(overwrite bool, jsons ...[]byte) ([]byte, error) {
	data := make([]map[string]json.RawMessage, len(jsons))
	for i, j := range jsons {
		m := make(map[string]json.RawMessage)
		if err := json.Unmarshal(j, &m); err != nil {
			return nil, err
		}
		data[i] = m
	}
	final := map[string]json.RawMessage{}
	for _, d := range data {
		for k, v := range d {
			finalV, finalFound := final[k]
			if !finalFound {
				final[k] = v
				continue
			}
			if IsJSONArray(finalV) {
				if !IsJSONArray(v) {
					v = Arrayify(v)
				}
				var err error
				final[k], err = MergeJSONArrays(finalV, v)
				if err != nil {
					return nil, err
				}
				continue
			}
			if IsJSONObject(finalV) {
				if !IsJSONObject(v) {
					return nil, errors.Errorf("cannot merge json object and non-json object for key '%s'", k)
				}
				merged, err := MergeJSONObjects(overwrite, finalV, v)
				if err != nil {
					return nil, err
				}
				final[k] = merged
			}
			if overwrite {
				final[k] = v
			}
		}
	}
	return json.Marshal(final)
}

// MergeJSONArrays merges two json arrays into one
func MergeJSONArrays(arrays ...[]byte) ([]byte, error) {
	datas := make([][]json.RawMessage, len(arrays))
	for i, a := range arrays {
		j := make([]json.RawMessage, 0)
		if err := json.Unmarshal(a, &j); err != nil {
			return nil, err
		}
		datas[i] = j
	}
	out := make([]json.RawMessage, 0)
	for _, a := range datas {
		for _, aEl := range a {
			found := false
			for _, oEl := range out {
				if bytes.Equal(aEl, oEl) {
					found = true
					break
				}
			}
			if !found {
				out = append(out, aEl)
			}
		}
	}

	return json.Marshal(out)
}

// Arrayify creates an JSON array with the passed object in
func Arrayify(o []byte) []byte {
	return bytes.Join(
		[][]byte{
			{'['},
			o,
			{']'},
		}, nil,
	)
}

// UnwrapString removes the " around a string if they exist
func UnwrapString(s []byte) []byte {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
