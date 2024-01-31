package config_navigator_test

import (
	cn "config-navigator"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path"
	"testing"
)

func TestSimpleConfig(t *testing.T) {

	node := loadTestFile(t, "simple.json")

	s, err := cn.StringVal("simpleOne.String", node)

	assert.EqualValues(t, "abc", s)
	assert.NoError(t, err)

	b, err := cn.BoolVal("simpleOne.Bool", node)
	assert.True(t, b)
	assert.NoError(t, err)

	i, err := cn.IntVal("simpleOne.Int", node)
	assert.NoError(t, err)
	assert.EqualValues(t, 32, i)

	f, err := cn.Float64Val("simpleOne.Float", node)
	assert.NoError(t, err)
	assert.EqualValues(t, 32.22, f)

	sa, err := cn.Array("simpleOne.StringArray", node, false)
	assert.NoError(t, err)
	assert.EqualValues(t, sa[1].(string), "b", node)

	sa, err = cn.Array("simpleOne.StringArrayX", node, false)
	assert.NoError(t, err)

	sa, err = cn.Array("simpleOne.Bool", node, false)
	assert.Error(t, err)

	o, err := cn.ObjectVal("simpleOne.Bool", node, false)
	assert.Nil(t, o)
	assert.Error(t, err)

	o, err = cn.ObjectVal("simpleOne", node, false)
	assert.NotNil(t, o)
	assert.NoError(t, err)

}

func TestMissingPath(t *testing.T) {

	node := loadTestFile(t, "simple.json")

	_, err := cn.StringVal("missing.path", node)
	assert.Error(t, err)

	_, err = cn.BoolVal("missing.path", node)
	assert.Error(t, err)

	_, err = cn.IntVal("missing.path", node)
	assert.Error(t, err)

	_, err = cn.Float64Val("missing.path", node)
	assert.Error(t, err)

	_, err = cn.Array("missing.path", node, true)
	assert.Error(t, err)

	_, err = cn.Array("missing.pathx", node, false)
	assert.NoError(t, err)

	_, err = cn.ObjectVal("missing.path", node, true)
	assert.Error(t, err)

	o, err := cn.ObjectVal("simpleOnex", node, false)
	assert.Nil(t, o)
	assert.NoError(t, err)

}

func TestPathExistence(t *testing.T) {

	node := loadTestFile(t, "simple.json")

	assert.True(t, cn.PathExists("simpleOne.Bool", node))

	assert.False(t, cn.PathExists("simpleX.Bool", node))
	assert.False(t, cn.PathExists("", node))
	assert.False(t, cn.PathExists(".....", node))

}

func TestWrongType(t *testing.T) {
	node := loadTestFile(t, "simple.json")

	i, err := cn.IntVal("simpleOne.String", node)
	assert.EqualValues(t, 0, i)
	assert.Error(t, err)

	b, err := cn.BoolVal("simpleOne.String", node)
	assert.False(t, b)
	assert.Error(t, err)

	f, err := cn.Float64Val("simpleOne.String", node)
	assert.EqualValues(t, 0, f)
	assert.Error(t, err)

	s, err := cn.StringVal("simpleOne.Bool", node)
	assert.EqualValues(t, "", s)
	assert.Error(t, err)
}

func loadTestFile(t *testing.T, file string) cn.ConfigNode {
	fp := path.Join("testdata", file)

	f, err := os.Open(fp)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Fatalf("Unable to close test file %s: %s", file, err.Error())
		}
	}(f)

	if err != nil {
		t.Fatalf(err.Error())
	}
	var result cn.ConfigNode
	bytes, _ := io.ReadAll(f)

	if err = json.Unmarshal(bytes, &result); err != nil {
		t.Fatalf("Problem unmarshalling JSON from %s: %s", file, err.Error())
	}

	return result
}
