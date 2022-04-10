package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/go-querystring/query"
)

func ParseQuery(params interface{}) (string, error) {
	switch v := params.(type) {
	case string:
		return v, nil
	case nil:
		return "", nil
	default:
		val := reflect.ValueOf(params)
		if val.Kind() == reflect.Struct {
			return "", fmt.Errorf("unknown type of params, must be string, struct or nil")
		}
		q := reflect.ValueOf(params)
		if q.Kind() == reflect.Ptr && q.IsNil() {
			return "", nil
		}

		qs, err := query.Values(params)
		if err != nil {
			return "", err
		}

		return qs.Encode(), nil
	}

}

func JsonifyBody(body interface{}) (*bytes.Buffer, error) {

	if body == nil {
		return bytes.NewBufferString(""), nil
	}

	var buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %v", err)
	}

	return buf, nil
}
