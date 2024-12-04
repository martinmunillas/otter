package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenJson(t *testing.T) {
	testcases := []struct {
		in  map[string]interface{}
		out map[string]string
		err bool
	}{
		{
			in: map[string]interface{}{
				"name": "John",
				"address": map[string]interface{}{
					"city": "New York",
					"zip": map[string]interface{}{
						"code":      "10001",
						"extension": "1234",
					},
				},
				"age": "30",
			},
			out: map[string]string{
				"name":                  "John",
				"address.city":          "New York",
				"address.zip.code":      "10001",
				"address.zip.extension": "1234",
				"age":                   "30",
			},
			err: false,
		},
		{
			in: map[string]interface{}{
				"name": "John",
				"address": map[string]interface{}{
					"city": "New York",
					"zip": map[string]interface{}{
						"code":      "10001",
						"extension": "1234",
					},
				},
				"age": 30, // only strings allowed
			},
			out: map[string]string{},
			err: true,
		},
	}

	for _, testcase := range testcases {
		m, err := flattenJson(testcase.in)
		if testcase.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, testcase.out, m)
		}
	}
}
