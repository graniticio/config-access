package config_access

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

func SetField(fieldName string, path string, target interface{}, config ConfigNode) error {

	if !PathExists(path, config) {
		return MissingPathError{message: "No value found at " + path}
	}

	targetReflect := reflect.ValueOf(target).Elem()
	targetField := targetReflect.FieldByName(fieldName)

	k := targetField.Type().Kind()

	switch k {
	case reflect.String:
		s, _ := StringVal(path, config)
		targetField.SetString(s)
	case reflect.Bool:
		b, _ := BoolVal(path, config)
		targetField.SetBool(b)
	case reflect.Int:
		i, _ := IntVal(path, config)
		targetField.SetInt(int64(i))
	case reflect.Float64:
		f, _ := Float64Val(path, config)
		targetField.SetFloat(f)
	case reflect.Map:

		if v, err := ObjectVal(path, config, false); err == nil {
			if err = populateMapField(targetField, v); err != nil {
				return err
			}
		} else {
			return err
		}
	case reflect.Slice:
		populateSlice(targetField, path, config)

	default:
		m := fmt.Sprintf("Unable to use value at path %s as target field %s is not a suppported type (%s)", path, fieldName, k)
		return errors.New(m)
	}

	return nil
}

// PopulateFromRoot sets the fields on the supplied target object using the whole supplied config document.
// This is achieved using Go's json.Marshal to convert the data
// back into text JSON and then json.Unmarshal to unmarshal back into the target.
func PopulateFromRoot(target interface{}, config ConfigNode) error {

	wrapper := make(ConfigNode)

	wrapper["root"] = config

	return Populate("root", target, wrapper)
}

// Populate sets the fields on the supplied target object using the data
// at the supplied path. This is achieved using Go's json.Marshal to convert the data
// back into text JSON and then json.Unmarshal to unmarshal back into the target.
func Populate(path string, target interface{}, config ConfigNode) error {
	if !PathExists(path, config) {
		return MissingPathError{message: "No value found at " + path}
	}

	//Already check if path exists
	object, _ := ObjectVal(path, config, false)

	if data, err := json.Marshal(object); err != nil {
		m := fmt.Sprintf("%T cannot be marshalled to JSON", object)
		return errors.New(m)
	} else if json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("%T cannot be populated with %v to JSON", object, data)
	}

	return nil

}

func populateSlice(targetField reflect.Value, path string, config ConfigNode) {

	v := Value(path, config)

	data, _ := json.Marshal(v)

	vt := targetField.Type()
	nt := reflect.New(vt)

	jTarget := nt.Interface()
	json.Unmarshal(data, &jTarget)

	vr := reflect.ValueOf(jTarget)
	targetField.Set(vr.Elem())

}

func populateMapField(targetField reflect.Value, contents map[string]interface{}) error {
	var err error

	m := reflect.MakeMap(targetField.Type())
	targetField.Set(m)

	for k, v := range contents {

		keyVal := reflect.ValueOf(k)
		vVal := reflect.ValueOf(v)

		if vVal.Kind() == reflect.Slice {
			vVal, err = arrayVal(vVal)

			if err != nil {
				return err
			}
		}

		m.SetMapIndex(keyVal, vVal)

	}

	return nil
}

func arrayVal(a reflect.Value) (reflect.Value, error) {

	v := a.Interface().([]interface{})
	l := len(v)

	if l == 0 {

		return reflect.Zero(reflect.TypeOf(v)), errors.New("cannot use an empty array as a value in a Map")

	}

	var s reflect.Value

	switch t := v[0].(type) {
	case string:
		s = reflect.MakeSlice(reflect.TypeOf([]string{}), 0, 0)
	default:
		m := fmt.Sprintf("Cannot use an array of %T as a value in a Map.", t)
		return reflect.Zero(reflect.TypeOf(v[0])), errors.New(m)
	}

	for _, elem := range v {

		s = reflect.Append(s, reflect.ValueOf(elem))

	}

	return s, nil
}
