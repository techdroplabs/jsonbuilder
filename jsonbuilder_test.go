package jsonbuilder_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torbinsky/jsonbuilder"
)

type Foo struct {
	Name    string
	Age     int
	Address string
}

func TestCustomMarshallingUsed(t *testing.T) {
	// Test a marshaller that is different but actually gonna make some valid JSON
	m := func(i interface{}) ([]byte, error) {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		// Default marshaller uses default encoding settings, so we try to set them to non-default values
		enc.SetIndent("", "    ")
		enc.SetEscapeHTML(false)

		err := enc.Encode(i)

		return buf.Bytes(), err
	}
	foo := Foo{"someone", 25, "earth"}

	// Marshall and compare result with what comes out of the original function. They should be the same.
	a, _ := m(foo)
	jh := jsonbuilder.FromMarshaller(foo, m)
	b := jh.Marshal()
	// Use JSONEq because things like field ordering may change, but the effective JSON should be equivalent
	assert.JSONEq(t, string(a), b, "The encodings should produce the same effective JSON.")

	// Try something a little more obvious now
	m = func(i interface{}) ([]byte, error) {
		return []byte("test"), nil
	}
	jh = jsonbuilder.FromMarshaller(foo, m)
	b = jh.Marshal()
	assert.Equal(t, "test", string(b), `"We should get "test"`)
}
