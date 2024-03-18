package config_access_test

import (
	ca "github.com/graniticio/config-access"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SimpleConfig struct {
	String         string
	Bool           bool
	Int            int
	Float          float64
	StringArray    []string
	FloatArray     []float64
	IntArray       []int
	StringMap      map[string]string
	Unsupported    *SimpleConfig
	StringArrayMap map[string][]string
}

func TestPopulationWithEntireDocument(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "root.json")
	yamlConf := loadYamlTestFile(t, "root.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {
		var sc SimpleConfig

		err := ca.PopulateFromRoot(&sc, node)

		assert.NoError(t, err)

		assert.EqualValues(t, "abc", sc.String)

		assert.True(t, sc.Bool)

		assert.EqualValues(t, 32, sc.Int)

		assert.EqualValues(t, 32.22, sc.Float)

		m := sc.StringMap

		assert.NotNil(t, m)

		assert.EqualValues(t, 3, len(sc.FloatArray))
	}
}

func TestPopulateObject(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		var sc SimpleConfig

		err := ca.Populate("simpleOne", &sc, node)

		assert.NoError(t, err)

		assert.EqualValues(t, "abc", sc.String)

		assert.True(t, sc.Bool)

		assert.EqualValues(t, 32, sc.Int)

		assert.EqualValues(t, 32.22, sc.Float)

		m := sc.StringMap

		assert.NotNil(t, m)

		assert.EqualValues(t, 3, len(sc.FloatArray))
	}

}

func TestPopulateOutOfBoundsNumbers(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "flipped-numbers.json")
	yamlConf := loadYamlTestFile(t, "flipped-numbers.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		var sc SimpleConfig

		err := ca.Populate("simpleOne", &sc, node)

		assert.Nil(t, err)
	}
}

func TestSetField(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {

		var sc SimpleConfig

		if err := ca.SetField("String", "simpleOne.String", &sc, node); err != nil {
			t.FailNow()
		}

		if err := ca.SetField("Bool", "simpleOne.Bool", &sc, node); err != nil {
			t.FailNow()
		}

		if err := ca.SetField("Int", "simpleOne.Int", &sc, node); err != nil {
			t.FailNow()
		}

		if err := ca.SetField("Float", "simpleOne.Float", &sc, node); err != nil {
			t.FailNow()
		}

		if err := ca.SetField("IntArray", "simpleOne.IntArray", &sc, node); err != nil {
			t.FailNow()
		}

		if err := ca.SetField("StringMap", "simpleOne.StringMap", &sc, node); err != nil {
			t.FailNow()
		}

		if err := ca.SetField("Unsupported", "simpleOne.IntArray", &sc, node); err == nil {
			t.FailNow()
		}

		if err := ca.SetField("StringMap", "missing.path", &sc, node); err == nil {
			t.FailNow()
		}

		if err := ca.SetField("StringMap", "simpleOne.Bool", &sc, node); err == nil {
			t.FailNow()
		}

		if err := ca.SetField("StringMap", "simpleOne.BoolA", &sc, node); err == nil {
			t.FailNow()
		}

		if _, err := ca.ObjectVal("simpleOne.Bool", node, false); err == nil {
			t.FailNow()
		}

		if err := ca.SetField("StringArrayMap", "simpleOne.StringArrayMap", &sc, node); err != nil {
			t.FailNow()
		}

		if err := ca.SetField("StringArrayMap", "simpleOne.EmptyStringArrayMap", &sc, node); err == nil {
			t.FailNow()
		}

		if err := ca.SetField("StringArrayMap", "simpleOne.BoolArrayMap", &sc, node); err == nil {
			t.FailNow()
		}
	}
}

func TestPopulateObjectMissingPath(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {
		var sc SimpleConfig

		err := ca.Populate("undefined", &sc, node)

		if mpe, okay := err.(ca.MissingPathError); okay {
			assert.NotEmpty(t, mpe.Error())
		}

		assert.NotNil(t, err)
	}
}

func TestPopulateInvalid(t *testing.T) {

	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {
		var sc SimpleConfig

		err := ca.Populate("invalidConfig", &sc, node)
		assert.NoError(t, err)
	}
}

func TestPopulateWithUnmarshalableContent(t *testing.T) {
	jsonConf := loadJsonTestFile(t, "simple.json")
	yamlConf := loadYamlTestFile(t, "simple.yaml")

	for _, node := range []ca.ConfigNode{jsonConf, yamlConf} {
		var sc SimpleConfig

		so := node["simpleOne"].(ca.ConfigNode)
		so["Int"] = make(chan int)

		err := ca.Populate("simpleOne", &sc, node)
		assert.Error(t, err)
	}
}
