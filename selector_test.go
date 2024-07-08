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

		ia, err := cs.IntArray("simpleOne.IntArray")
		assert.NoError(t, err)
		assert.EqualValues(t, ia[1], 2)

		fa, err := cs.Float64Array("simpleOne.FloatArray")
		assert.NoError(t, err)
		assert.EqualValues(t, fa[1], 2.0)

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

func TestDefaultValueOnMissing(t *testing.T) {
	pv := map[string]interface{}{}
	sel := ca.SelectorFromPathValues(pv)

	s := sel.Value("x.y.x", ca.Opts{OnMissing: "default"})
	assert.EqualValues(t, "default", s)

	dov := make(ca.ConfigNode)
	ov, err := sel.ObjectVal("x.y.z", ca.Opts{OnMissing: dov})
	assert.NoError(t, err)
	assert.EqualValues(t, dov, ov)

	sdv := "default"
	sv, err := sel.StringVal("x.y.z", ca.Opts{OnMissing: sdv})
	assert.NoError(t, err)
	assert.EqualValues(t, sdv, sv)

	idv := -2
	iv, err := sel.IntVal("x.y.z", ca.Opts{OnMissing: -2})
	assert.NoError(t, err)
	assert.EqualValues(t, idv, iv)

	fdv := -3.2
	fv, err := sel.Float64Val("x.y.z", ca.Opts{OnMissing: -3.2})
	assert.NoError(t, err)
	assert.EqualValues(t, fdv, fv)

	bdv := true
	bv, err := sel.BoolVal("x.y.z", ca.Opts{OnMissing: true})
	assert.NoError(t, err)
	assert.EqualValues(t, bdv, bv)

	ifadv := []interface{}{true, 1}
	ifav, err := sel.Array("x.y.z", ca.Opts{OnMissing: ifadv})
	assert.NoError(t, err)
	assert.EqualValues(t, ifav, ifadv)

	iadv := []int{5, 6}
	iav, err := sel.IntArray("x.y.z", ca.Opts{OnMissing: iadv})
	assert.NoError(t, err)
	assert.EqualValues(t, iav, iadv)

	fadv := []float64{5.2, 6.3}
	fav, err := sel.Float64Array("x.y.z", ca.Opts{OnMissing: fadv})
	assert.NoError(t, err)
	assert.EqualValues(t, fav, fadv)

	sadv := []string{"a", "d"}
	sav, err := sel.StringArray("x.y.z", ca.Opts{OnMissing: sadv})
	assert.NoError(t, err)
	assert.EqualValues(t, sav, sadv)

}

func TestSelectorFromPathValues(t *testing.T) {

	pv := map[string]interface{}{
		"a.b.c.d": "e",
		"x.y.z":   10,
		"   ":     "ignore",
		"a.b.c.f": "S",
	}

	s := ca.SelectorFromPathValues(pv)
	assert.NotNil(t, s)

	s1, err := s.StringVal("a.b.c.d")
	assert.NoError(t, err)
	assert.EqualValues(t, "e", s1)

	i, err := s.IntVal("x.y.z")
	assert.NoError(t, err)
	assert.EqualValues(t, 10, i)

	s2, err := s.StringVal("a.b.c.f")
	assert.NoError(t, err)
	assert.EqualValues(t, "S", s2)

	_, err = s.StringVal("   ")
	assert.Error(t, err)
}

func TestStringOrEnv(t *testing.T) {

	ef := func(s string) string {
		if s == "ENV_NAME" {
			return "ENV_VALUE"
		} else {
			return ""
		}
	}

	pv := map[string]interface{}{
		"value":        "VALUE",
		"env":          "$ENV_NAME",
		"missing":      "$MISSING_ENV_NAME",
		"notOnOs":      "$MCASCASCASASCQWEQWE",
		"customPrefix": "#ENV_NAME",
		"notString":    1,
	}

	s := ca.SelectorFromPathValues(pv)
	assert.NotNil(t, s)

	ev, err := s.StringOrEnv("notString")
	assert.Error(t, err)
	assert.Zero(t, ev)

	ev, err = s.StringOrEnv("notOnOs")
	assert.Error(t, err)
	assert.Zero(t, ev)

	ev, err = s.StringOrEnv("value", ca.Opts{EnvAccessFunc: ef})
	assert.NoError(t, err)
	assert.Equal(t, "VALUE", ev)

	ev, err = s.StringOrEnv("env", ca.Opts{EnvAccessFunc: ef})
	assert.NoError(t, err)
	assert.Equal(t, "ENV_VALUE", ev)

	ev, err = s.StringOrEnv("missing", ca.Opts{EnvAccessFunc: ef})
	assert.Error(t, err)
	assert.Zero(t, ev)

	ev, err = s.StringOrEnv("customPrefix", ca.Opts{EnvAccessFunc: ef, EnvVarPrefix: "#"})
	assert.NoError(t, err)
	assert.Equal(t, "ENV_VALUE", ev)

}
