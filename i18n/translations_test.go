package i18n

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessJson(t *testing.T) {
	testcases := []struct {
		in  io.Reader
		out map[string]string
		err bool
	}{
		{
			in: strings.NewReader(`{
		"name": "John",
		"address": {
			"city": "New York",
			"zip": {
				"code": "10001",
				"extension": "1234"
			}
		},
		"age": "30"
	}`),
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
			in: strings.NewReader(`{
		"name": "John",
		"address": {
			"city": "New York",
			"zip": {
				"code": "10001",
				"extension": "1234"
			}
		},
		"age": 30
	}`),
			out: map[string]string{},
			err: true,
		},
	}

	for _, testcase := range testcases {
		m, err := processLang(testcase.in)
		if testcase.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, testcase.out, m)
		}
	}
}
