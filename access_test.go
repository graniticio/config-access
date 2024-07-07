package config_access_test

import (
	ca "github.com/graniticio/config-access"
	"github.com/stretchr/testify/assert"
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

func TestStringArray(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {
		s, err := ca.StringArray("simpleOne.StringArray", node)

		assert.NoError(t, err)
		assert.EqualValues(t, []string{"a", "b", "c"}, s)

		s, err = ca.StringArray("missing.StringArray", node)
		assert.Error(t, err)

		s, err = ca.StringArray("simpleOne.IntArray", node)
		assert.Error(t, err)
	}
}

func TestIntArray(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {
		ia, err := ca.IntArray("simpleOne.IntArray", node)

		assert.NoError(t, err)
		assert.EqualValues(t, []int{1, 2, 3}, ia)

		ia, err = ca.IntArray("simpleOne.StringArray", node)
		assert.Error(t, err)

		ia, err = ca.IntArray("missing.IntArray", node)
		assert.Error(t, err)
	}
}

func TestFloat64Array(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {
		ia, err := ca.Float64Array("simpleOne.FloatArray", node)

		assert.NoError(t, err)
		assert.EqualValues(t, []float64{1.0, 2.0, 3.0}, ia)

		ia, err = ca.Float64Array("simpleOne.StringArray", node)
		assert.Error(t, err)

		ia, err = ca.Float64Array("simpleOne.IntArray", node)
		assert.NoError(t, err)
	}
}
