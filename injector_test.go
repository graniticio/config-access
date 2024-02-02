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

func TestPopulateObject(t *testing.T) {

	config := loadTestFile(t, "simple.json")

	var sc SimpleConfig

	err := ca.Populate("simpleOne", &sc, config)

	assert.NoError(t, err)

	assert.EqualValues(t, "abc", sc.String)

	assert.True(t, sc.Bool)

	assert.EqualValues(t, 32, sc.Int)

	assert.EqualValues(t, 32.22, sc.Float)

	m := sc.StringMap

	assert.NotNil(t, m)

	assert.EqualValues(t, 3, len(sc.FloatArray))

}

func TestPopulateOutOfBoundsNumbers(t *testing.T) {

	config := loadTestFile(t, "flipped-numbers.json")

	var sc SimpleConfig

	err := ca.Populate("simpleOne", &sc, config)

	assert.Nil(t, err)
}

func TestSetField(t *testing.T) {

	config := loadTestFile(t, "simple.json")

	var sc SimpleConfig

	if err := ca.SetField("String", "simpleOne.String", &sc, config); err != nil {
		t.FailNow()
	}

	if err := ca.SetField("Bool", "simpleOne.Bool", &sc, config); err != nil {
		t.FailNow()
	}

	if err := ca.SetField("Int", "simpleOne.Int", &sc, config); err != nil {
		t.FailNow()
	}

	if err := ca.SetField("Float", "simpleOne.Float", &sc, config); err != nil {
		t.FailNow()
	}

	if err := ca.SetField("IntArray", "simpleOne.IntArray", &sc, config); err != nil {
		t.FailNow()
	}

	if err := ca.SetField("StringMap", "simpleOne.StringMap", &sc, config); err != nil {
		t.FailNow()
	}

	if err := ca.SetField("Unsupported", "simpleOne.IntArray", &sc, config); err == nil {
		t.FailNow()
	}

	if err := ca.SetField("StringMap", "missing.path", &sc, config); err == nil {
		t.FailNow()
	}

	if err := ca.SetField("StringMap", "simpleOne.Bool", &sc, config); err == nil {
		t.FailNow()
	}

	if err := ca.SetField("StringMap", "simpleOne.BoolA", &sc, config); err == nil {
		t.FailNow()
	}

	if _, err := ca.ObjectVal("simpleOne.Bool", config, false); err == nil {
		t.FailNow()
	}

	if err := ca.SetField("StringArrayMap", "simpleOne.StringArrayMap", &sc, config); err != nil {
		t.FailNow()
	}

	if err := ca.SetField("StringArrayMap", "simpleOne.EmptyStringArrayMap", &sc, config); err == nil {
		t.FailNow()
	}

	if err := ca.SetField("StringArrayMap", "simpleOne.BoolArrayMap", &sc, config); err == nil {
		t.FailNow()
	}
}

func TestPopulateObjectMissingPath(t *testing.T) {
	config := loadTestFile(t, "simple.json")
	var sc SimpleConfig

	err := ca.Populate("undefined", &sc, config)

	if mpe, okay := err.(ca.MissingPathError); okay {
		assert.NotEmpty(t, mpe.Error())
	}

	assert.NotNil(t, err)

}

func TestPopulateInvalid(t *testing.T) {

	config := loadTestFile(t, "simple.json")
	var sc SimpleConfig

	err := ca.Populate("invalidConfig", &sc, config)
	assert.NoError(t, err)
}

func TestPopulateWithUnmarshalableContent(t *testing.T) {
	config := loadTestFile(t, "simple.json")
	var sc SimpleConfig

	so := config["simpleOne"].(ca.ConfigNode)
	so["Int"] = make(chan int)

	err := ca.Populate("simpleOne", &sc, config)
	assert.Error(t, err)
}
