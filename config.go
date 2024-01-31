// Copyright 2024 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package config_navigator

import (
	"errors"
	"fmt"
	"strings"
)

const PathSeparator = "."

type ConfigNode = map[string]interface{}

func PathExists(path string, node ConfigNode) bool {
	value := Value(path, node)

	return value != nil
}

// Value returns the value at the supplied path or nil if the path does not exist of points to a null value.
func Value(path string, node ConfigNode) interface{} {

	splitPath := strings.Split(path, PathSeparator)

	return configVal(splitPath, node)

}

// ObjectVal returns a map representing an object or nil if the path does not exist or points to a null value. An error
// is returned if the value cannot be interpreted as an object (a key in the configuration that has child keys rather than
// a value.
// If errIfMissing is set to true, an error will be return if the supplied path does not exist otherwise a nil
// array without and error will be returned.
func ObjectVal(path string, node ConfigNode, errIfMissing bool) (ConfigNode, error) {

	if errIfMissing && !PathExists(path, node) {
		return nil, errors.New("No such path " + path)
	}

	value := Value(path, node)

	if value == nil {
		return nil, nil
	} else if v, found := value.(ConfigNode); found {
		return v, nil
	}

	return nil, fmt.Errorf("unable to convert the value at %s to a ConfigNode", path)
}

// StringVal returns the string value of the string at the supplied path. Does not convert other types to
// a string, so will return an error if the value is not already a string.
func StringVal(path string, node ConfigNode) (string, error) {

	v := Value(path, node)

	if v == nil {
		return "", errors.New("No string value found at " + path)
	}

	s, found := v.(string)

	if found {
		return s, nil
	}

	return "", fmt.Errorf("Value at %s is %q and cannot be converted to a string", path, v)

}

// IntVal returns the int value of the  number at the supplied path. JSON numbers
// are internally represented by Go as a float64, so no error will be returned, but data might be lost
// if the number does not actually represent an int. An error will be returned if the value is not a number
// or cannot be converted to an int.
func IntVal(path string, node ConfigNode) (int, error) {

	v := Value(path, node)

	if v == nil {
		return 0, errors.New("No such path " + path)
	} else if f, found := v.(float64); found {
		return int(f), nil
	}

	return 0, fmt.Errorf("alue at %s is %q and cannot be converted to an int", path, v)

}

// Float64Val returns the float64 value of the  number at the supplied path. An error will be returned if the value is not a number.
func Float64Val(path string, node ConfigNode) (float64, error) {

	v := Value(path, node)

	if v == nil {
		return 0, errors.New("No such path " + path)
	} else if f, found := v.(float64); found {
		return f, nil
	}

	return 0, fmt.Errorf("value at %s is %q and cannot be converted to a float64", path, v)
}

// Array returns the value of an array of objects at the supplied path. Caution should be used when calling this method
// as behaviour is undefined for arrays of types other than []ConfigNode.
//
// If errIfMissing is set to true, an error will be return if the supplied path does not exist otherwise a nil
// array without and error will be returned.
func Array(path string, node ConfigNode, errIfMissing bool) ([]interface{}, error) {

	if errIfMissing && !PathExists(path, node) {
		return nil, errors.New("No such path " + path)
	}

	value := Value(path, node)

	if value == nil {
		return nil, nil
	} else if v, found := value.([]interface{}); found {
		return v, nil
	}

	return nil, fmt.Errorf("unable to convert the value at %s to an array", path)

}

// BoolVal returns the bool value of the bool at the supplied path. An error will be returned if the value is not a JSON bool.
// Note this method only supports the JSON definition of bool (true, false) not the Go definition (true, false, 1, 0 etc) or
// extended YAML definitions.
func BoolVal(path string, node ConfigNode) (bool, error) {

	v := Value(path, node)

	if v == nil {
		return false, errors.New("No such path " + path)
	}

	if b, found := v.(bool); found {
		return b, nil
	}

	return false, fmt.Errorf("Value at %s is %q and cannot be converted to a bool", path, v)

}

func configVal(path []string, jsonMap ConfigNode) interface{} {

	var result interface{}
	result = jsonMap[path[0]]

	if result == nil {
		return nil
	}

	if len(path) == 1 {
		return result
	}

	remainPath := path[1:]
	return configVal(remainPath, result.(ConfigNode))
}
