package config_access_test

import (
	"encoding/json"
	ca "github.com/graniticio/config-access"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"testing"
)

func TestSimpleConfig(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		s, err := ca.StringVal("simpleOne.String", node)

		assert.EqualValues(t, "abc", s)
		assert.NoError(t, err)

		b, err := ca.BoolVal("simpleOne.Bool", node)
		assert.True(t, b)
		assert.NoError(t, err)

		i, err := ca.IntVal("simpleOne.Int", node)
		assert.NoError(t, err)
		assert.EqualValues(t, 32, i)

		f, err := ca.Float64Val("simpleOne.Float", node)
		assert.NoError(t, err)
		assert.EqualValues(t, 32.22, f)

		sa, err := ca.Array("simpleOne.StringArray", node, false)
		assert.NoError(t, err)
		assert.EqualValues(t, sa[1].(string), "b", node)

		sa, err = ca.Array("simpleOne.StringArrayX", node, false)
		assert.NoError(t, err)

		sa, err = ca.Array("simpleOne.Bool", node, false)
		assert.Error(t, err)

		o, err := ca.ObjectVal("simpleOne.Bool", node, false)
		assert.Nil(t, o)
		assert.Error(t, err)

		o, err = ca.ObjectVal("simpleOne", node, false)
		assert.NotNil(t, o)
		assert.NoError(t, err)
	}

}

func TestSimpleConfigViaSelector(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		cs := ca.NewDefaultSelector(node, false, false)

		s, err := cs.StringVal("simpleOne.String")

		assert.EqualValues(t, "abc", s)
		assert.NoError(t, err)

		b, err := cs.BoolVal("simpleOne.Bool")
		assert.True(t, b)
		assert.NoError(t, err)

		i, err := cs.IntVal("simpleOne.Int")
		assert.NoError(t, err)
		assert.EqualValues(t, 32, i)

		f, err := cs.Float64Val("simpleOne.Float")
		assert.NoError(t, err)
		assert.EqualValues(t, 32.22, f)

		sa, err := cs.Array("simpleOne.StringArray")
		assert.NoError(t, err)
		assert.EqualValues(t, sa[1].(string), "b")

		sa, err = cs.Array("simpleOne.StringArrayX")
		assert.NoError(t, err)

		sa, err = cs.Array("simpleOne.Bool")
		assert.Error(t, err)

		o, err := cs.ObjectVal("simpleOne.Bool")
		assert.Nil(t, o)
		assert.Error(t, err)

		o, err = cs.ObjectVal("simpleOne")
		assert.NotNil(t, o)
		assert.NoError(t, err)
	}

}

func TestPostFlushBehaviour(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		cs := ca.NewDefaultSelector(node, false, false)

		s, err := cs.StringVal("simpleOne.String")

		assert.EqualValues(t, "abc", s)
		assert.NoError(t, err)

		cs.Flush()

		b := cs.PathExists("simpleOne")
		assert.False(t, b)

		v := cs.Value("simpleOne")
		assert.Nil(t, v)

		s, err = cs.StringVal("simpleOne.String")
		assert.Error(t, err)

		_, err = cs.BoolVal("simpleOne.Bool")
		assert.Error(t, err)

		_, err = cs.IntVal("simpleOne.Int")
		assert.Error(t, err)

		_, err = cs.Float64Val("simpleOne.Float")
		assert.Error(t, err)

		_, err = cs.Array("simpleOne.StringArray")
		assert.Error(t, err)

		_, err = cs.Array("simpleOne.StringArrayX")
		assert.Error(t, err)

		_, err = cs.ObjectVal("simpleOne.Bool")
		assert.Error(t, err)
	}

}

func TestMapCallByReferenceInSelector(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		cs := ca.NewDefaultSelector(node, false, false)

		s, err := cs.StringVal("simpleOne.String")

		assert.EqualValues(t, "abc", s)
		assert.NoError(t, err)

		so := node["simpleOne"].(ca.ConfigNode)
		so["String"] = "123"

		s, err = cs.StringVal("simpleOne.String")

		assert.EqualValues(t, "123", s)
		assert.NoError(t, err)

		c := cs.Config()["simpleOne"].(ca.ConfigNode)
		c["String"] = "XYZ"

		s, err = cs.StringVal("simpleOne.String")

		assert.EqualValues(t, "XYZ", s)
		assert.NoError(t, err)
	}

}

func TestMissingPath(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		_, err := ca.StringVal("missing.path", node)
		assert.Error(t, err)

		_, err = ca.BoolVal("missing.path", node)
		assert.Error(t, err)

		_, err = ca.IntVal("missing.path", node)
		assert.Error(t, err)

		_, err = ca.Float64Val("missing.path", node)
		assert.Error(t, err)

		_, err = ca.Array("missing.path", node, true)
		assert.Error(t, err)

		_, err = ca.Array("missing.pathx", node, false)
		assert.NoError(t, err)

		_, err = ca.ObjectVal("missing.path", node, true)
		assert.Error(t, err)

		o, err := ca.ObjectVal("simpleOnex", node, false)
		assert.Nil(t, o)
		assert.NoError(t, err)
	}

}

func TestPathExistence(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		assert.True(t, ca.PathExists("simpleOne.Bool", node))

		assert.False(t, ca.PathExists("simpleX.Bool", node))
		assert.False(t, ca.PathExists("", node))
		assert.False(t, ca.PathExists(".....", node))
	}
}

func TestWrongType(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		i, err := ca.IntVal("simpleOne.String", node)
		assert.EqualValues(t, 0, i)
		assert.Error(t, err)

		b, err := ca.BoolVal("simpleOne.String", node)
		assert.False(t, b)
		assert.Error(t, err)

		f, err := ca.Float64Val("simpleOne.String", node)
		assert.EqualValues(t, 0, f)
		assert.Error(t, err)

		s, err := ca.StringVal("simpleOne.Bool", node)
		assert.EqualValues(t, "", s)
		assert.Error(t, err)
	}
}

func loadJsonTestFile(t *testing.T, file string) ca.ConfigNode {
	return loadAndParse(t, file, parseJson)
}

func loadYamlTestFile(t *testing.T, file string) ca.ConfigNode {
	return loadAndParse(t, file, parseYaml)
}

func loadAndParse(t *testing.T, file string, parse func(source []byte, target interface{}) error) ca.ConfigNode {
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
	var result ca.ConfigNode
	bytes, _ := io.ReadAll(f)

	if err = parse(bytes, &result); err != nil {
		t.Fatalf("Problem unmarshalling JSON from %s: %s", file, err.Error())
	}

	return result
}

func parseJson(source []byte, target interface{}) error {
	return json.Unmarshal(source, &target)
}

func parseYaml(source []byte, target interface{}) error {
	return yaml.Unmarshal(source, target)
}

func TestSimpleConfigQuietAccess(t *testing.T) {

	var invoked bool

	errorFunc := func(path string, err error) {
		invoked = true
	}

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		base := ca.NewDefaultSelector(node, true, true)
		cs := ca.NewDeferredErrorQuietSelector(base, errorFunc)

		assert.True(t, cs.PathExists("simpleOne.String"))
		assert.False(t, invoked)
		invoked = false

		s := cs.StringVal("simpleOne.String")

		assert.EqualValues(t, "abc", s)
		assert.False(t, invoked)
		invoked = false

		s = cs.StringVal("missing.String")

		assert.EqualValues(t, "", s)
		assert.True(t, invoked)
		invoked = false

		b := cs.BoolVal("simpleOne.Bool")
		assert.True(t, b)
		assert.False(t, invoked)
		invoked = false

		b = cs.BoolVal("missing.Bool")
		assert.False(t, b)
		assert.True(t, invoked)
		invoked = false

		i := cs.IntVal("simpleOne.Int")
		assert.EqualValues(t, 32, i)
		assert.False(t, invoked)
		invoked = false

		i = cs.IntVal("missing.Int")
		assert.EqualValues(t, 0, i)
		assert.True(t, invoked)
		invoked = false

		f := cs.Float64Val("simpleOne.Float")
		assert.False(t, invoked)
		assert.EqualValues(t, 32.22, f)
		invoked = false

		f = cs.Float64Val("missing.Float")
		assert.True(t, invoked)
		assert.EqualValues(t, 0, f)
		invoked = false

		sa := cs.Array("simpleOne.StringArray")
		assert.False(t, invoked)
		assert.EqualValues(t, sa[1].(string), "b")
		invoked = false

		sa = cs.Array("missing.StringArray")
		assert.True(t, invoked)
		assert.Nil(t, sa)
		invoked = false

		ov := cs.ObjectVal("simpleOne.StringMap")
		assert.False(t, invoked)
		assert.NotNil(t, ov)
		invoked = false

		ov = cs.ObjectVal("missing.StringMap")
		assert.True(t, invoked)
		assert.Nil(t, ov)
		invoked = false

		is := cs.Value("simpleOne.String")

		assert.NotNil(t, is)
		assert.False(t, invoked)
		invoked = false

		is = cs.Value("missing.String")

		assert.Nil(t, is)
		assert.False(t, invoked)
		invoked = false

	}
}

func TestGraniticSelectorErrorBehaviour(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		cs := ca.NewGraniticSelector(node)

		o, err := cs.ObjectVal("missing.path")
		assert.Nil(t, o)
		assert.NoError(t, err)

		a, err := cs.Array("missing.path")
		assert.Nil(t, a)
		assert.NoError(t, err)

	}
}

func TestNodeBehaviour(t *testing.T) {
	i, err := ca.IntVal("some.path", nil)

	assert.Equal(t, 0, i)
	assert.Error(t, err)

}
