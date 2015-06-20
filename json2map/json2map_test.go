package json2map

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {

	reader := bytes.NewBufferString(`{"a": 1, "b": "asdf"}`)
	m, err := NewWithRootKeyName("root").Convert(reader)
	require.NoError(t, err)
	assert.Equal(t, "1", m["root.a"].(json.Number).String())
	assert.Equal(t, "asdf", m["root.b"].(string))
}

func TestLargeInt(t *testing.T) {

	reader := bytes.NewBufferString(`{"a": 1356998399}`)
	m, err := NewWithRootKeyName("root").Convert(reader)
	require.NoError(t, err)
	assert.Equal(t, "1356998399", m["root.a"].(json.Number).String())
}

func TestFloat(t *testing.T) {

	reader := bytes.NewBufferString(`{"a": 1356998399.32}`)
	m, err := NewWithRootKeyName("root").Convert(reader)
	require.NoError(t, err)
	assert.Equal(t, "1356998399.32", m["root.a"].(json.Number).String())
}

func TestNested(t *testing.T) {

	reader := bytes.NewBufferString(`{"a": {"b": "asdf"}}`)
	m, err := NewWithRootKeyName("root").Convert(reader)
	require.NoError(t, err)
	assert.Equal(t, "asdf", m["root.a.b"].(string))
}

func TestArray(t *testing.T) {

	reader := bytes.NewBufferString(`{"a": [{"a0": "asdf"}, {"a1": [1, false]}]}`)
	m, err := NewWithRootKeyName("root").Convert(reader)
	require.NoError(t, err)
	assert.Equal(t, "asdf", m["root.a[0].a0"].(string))
	assert.Equal(t, "1", m["root.a[1].a1[0]"].(json.Number).String())
	assert.Equal(t, false, m["root.a[1].a1[1]"].(bool))
}
