package config_access_test

import (
	ca "github.com/graniticio/config-access"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestQuietAccessArrays(t *testing.T) {

	var invoked bool

	errorFunc := func(path string, err error) {
		invoked = true
	}

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		base := ca.NewDefaultSelector(node, true, true)
		cs := ca.NewDeferredErrorQuietSelector(base, errorFunc)

		sa := cs.Array("simpleOne.StringArray")
		assert.False(t, invoked)
		assert.EqualValues(t, sa[1].(string), "b")
		invoked = false

		sa = cs.Array("missing.StringArray")
		assert.True(t, invoked)
		assert.Nil(t, sa)
		invoked = false

		ssa := cs.StringArray("simpleOne.StringArray")
		assert.False(t, invoked)
		assert.EqualValues(t, ssa[1], "b")
		invoked = false

		ssa = cs.StringArray("missing.StringArray")
		assert.True(t, invoked)
		assert.Nil(t, ssa)
		invoked = false

		ia := cs.IntArray("simpleOne.IntArray")
		assert.False(t, invoked)
		assert.EqualValues(t, ia[1], 2)
		invoked = false

		ia = cs.IntArray("missing.IntArray")
		assert.True(t, invoked)
		assert.Nil(t, ia)
		invoked = false
	}
}
