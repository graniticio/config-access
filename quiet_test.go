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

func TestQuietSelectorFromPathValues(t *testing.T) {

	var invoked bool

	errorFunc := func(path string, err error) {
		invoked = true
	}

	pv := map[string]interface{}{
		"a.b.c": 1,
	}

	qs := ca.QuietSelectorFromPathValues(pv, errorFunc)

	i := qs.IntVal("a.b.c")
	assert.EqualValues(t, 1, i)
	assert.False(t, invoked)

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

		fa := cs.Float64Array("simpleOne.FloatArray")
		assert.False(t, invoked)
		assert.EqualValues(t, fa[1], 2.0)
		invoked = false

		fa = cs.Float64Array("missing.FloatArray")
		assert.True(t, invoked)
		assert.Nil(t, ia)
		invoked = false
	}
}

func TestQuietDefaultValueOnMissing(t *testing.T) {
	pv := map[string]interface{}{}

	errorFunc := func(path string, err error) {

	}

	sel := ca.QuietSelectorFromPathValues(pv, errorFunc)

	s := sel.Value("x.y.x", ca.Opts{OnMissing: "default"})
	assert.EqualValues(t, "default", s)

	dov := make(ca.ConfigNode)
	ov := sel.ObjectVal("x.y.z", ca.Opts{OnMissing: dov})
	assert.EqualValues(t, dov, ov)

	sdv := "default"
	sv := sel.StringVal("x.y.z", ca.Opts{OnMissing: sdv})
	assert.EqualValues(t, sdv, sv)

	idv := -2
	iv := sel.IntVal("x.y.z", ca.Opts{OnMissing: -2})
	assert.EqualValues(t, idv, iv)

	fdv := -3.2
	fv := sel.Float64Val("x.y.z", ca.Opts{OnMissing: -3.2})
	assert.EqualValues(t, fdv, fv)

	bdv := true
	bv := sel.BoolVal("x.y.z", ca.Opts{OnMissing: true})
	assert.EqualValues(t, bdv, bv)

	ifadv := []interface{}{true, 1}
	ifav := sel.Array("x.y.z", ca.Opts{OnMissing: ifadv})
	assert.EqualValues(t, ifav, ifadv)

	iadv := []int{5, 6}
	iav := sel.IntArray("x.y.z", ca.Opts{OnMissing: iadv})
	assert.EqualValues(t, iav, iadv)

	fadv := []float64{5.2, 6.3}
	fav := sel.Float64Array("x.y.z", ca.Opts{OnMissing: fadv})
	assert.EqualValues(t, fav, fadv)

	sadv := []string{"a", "d"}
	sav := sel.StringArray("x.y.z", ca.Opts{OnMissing: sadv})
	assert.EqualValues(t, sav, sadv)

}
