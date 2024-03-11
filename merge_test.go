package config_access_test

import (
	ca "github.com/graniticio/config-access"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeDetection(t *testing.T) {

	if ca.ConfigType("") != ca.ConfigString {
		t.FailNow()
	}

	if ca.ConfigType(true) != ca.ConfigBool {
		t.FailNow()
	}

	if ca.ConfigType(make(map[string]interface{})) != ca.ConfigMap {
		t.FailNow()
	}

	if ca.ConfigType([]interface{}{}) != ca.ConfigArray {
		t.FailNow()
	}

	if ca.ConfigType(1) != ca.ConfigUnknown {
		t.FailNow()
	}
}

func TestMergingWithArrayReplace(t *testing.T) {

	base := loadTestFile(t, "merge-base.json")
	additions := loadTestFile(t, "merge-additions.json")

	result := ca.Merge(base, additions, false)

	assert.NotNil(t, result)

	assert.EqualValues(t, result["baseOnly"].(string), "def")
	assert.EqualValues(t, result["baseString"].(string), "xyz")
	assert.EqualValues(t, result["baseNumber"].(float64), 200)
	assert.False(t, result["baseBool"].(bool))

	a, okay := result["baseArray"].([]interface{})

	assert.True(t, okay)
	assert.Len(t, a, 1)

	o, okay := result["baseObject"].(ca.ConfigNode)

	assert.True(t, okay)
	assert.EqualValues(t, o["objectField1"].(string), "inBase")
	assert.EqualValues(t, o["objectField2"].(string), "inAdditions")

}

func TestMergingWithArrayMerge(t *testing.T) {

	base := loadTestFile(t, "merge-base.json")
	additions := loadTestFile(t, "merge-additions.json")

	result := ca.Merge(base, additions, true)

	assert.NotNil(t, result)

	a, okay := result["baseArray"].([]interface{})

	assert.True(t, okay)
	assert.Len(t, a, 4)

	assert.EqualValues(t, a[0].(float64), 1)
	assert.EqualValues(t, a[1].(float64), 2)
	assert.EqualValues(t, a[2].(float64), 3)
	assert.EqualValues(t, a[3].(float64), 4)

}

/*baseString": "def",
  "baseNumber": 200,
  "baseBool": false,*/
