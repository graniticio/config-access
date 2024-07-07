// Copyright 2024 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package config_access

import "strings"

const PathSeparator = "."

type ConfigNode = map[string]interface{}

// MissingPathError indicates that the a problem was caused by there being no value at the supplied
// config path
type MissingPathError struct {
	message string
}

func (mp MissingPathError) Error() string {
	return mp.message
}

type Selector interface {
	PathExists(path string) bool
	Value(path string) interface{}
	ObjectVal(path string) (ConfigNode, error)
	StringVal(path string) (string, error)
	IntVal(path string) (int, error)
	Float64Val(path string) (float64, error)
	Array(path string) ([]interface{}, error)
	StringArray(path string) ([]string, error)
	IntArray(path string) ([]int, error)
	Float64Array(path string) ([]float64, error)
	BoolVal(path string) (bool, error)
	Flush()
	Config() ConfigNode
}

// Opts defines optional behaviour for accessing and interpreting config values
type Opts struct {
}

// SelectorFromPathValues creates a Selector from a map of config paths (e.g. my.config.path) and their
// associated values. Empty path values are ignored.
func SelectorFromPathValues(pathValues map[string]interface{}) Selector {

	store := make(map[string]interface{})

	for k, v := range pathValues {

		if strings.TrimSpace(k) == "" {
			continue
		}

		addValue(strings.Split(k, "."), v, store)

	}

	return NewDefaultSelector(store, true, true)

}

func addValue(path []string, value any, store map[string]interface{}) {

	first := path[0]

	if len(path) == 1 {
		store[first] = value
	} else {

		storeForFirst := store[first]

		if storeForFirst == nil {
			newFirst := make(map[string]interface{})
			store[first] = newFirst
			addValue(path[1:], value, newFirst)
		} else {
			addValue(path[1:], value, storeForFirst.(map[string]interface{}))
		}

	}

}

func NewDefaultSelector(config ConfigNode, errorOnMissingObjectPath, errorOnMissingArrayPath bool) Selector {
	ds := new(DefaultSelector)
	ds.config = config
	ds.errorOnMissingArrayPath = errorOnMissingArrayPath
	ds.errorOnMissingObjectPath = errorOnMissingObjectPath

	return ds
}

func NewGraniticSelector(config ConfigNode) Selector {
	ds := new(DefaultSelector)
	ds.config = config

	return ds
}

type DefaultSelector struct {
	errorOnMissingObjectPath bool
	errorOnMissingArrayPath  bool
	config                   ConfigNode
}

func (dfe *DefaultSelector) Flush() {
	dfe.config = nil
}

func (dfe *DefaultSelector) PathExists(path string) bool {
	return PathExists(path, dfe.config)
}

func (dfe *DefaultSelector) Value(path string) interface{} {
	return Value(path, dfe.config)
}

func (dfe *DefaultSelector) ObjectVal(path string) (ConfigNode, error) {
	return ObjectVal(path, dfe.config, dfe.errorOnMissingObjectPath)
}

func (dfe *DefaultSelector) StringVal(path string) (string, error) {
	return StringVal(path, dfe.config)
}

func (dfe *DefaultSelector) IntVal(path string) (int, error) {
	return IntVal(path, dfe.config)
}

func (dfe *DefaultSelector) Float64Val(path string) (float64, error) {
	return Float64Val(path, dfe.config)
}

func (dfe *DefaultSelector) Array(path string) ([]interface{}, error) {
	return Array(path, dfe.config, dfe.errorOnMissingArrayPath)
}

func (dfe *DefaultSelector) StringArray(path string) ([]string, error) {
	return StringArray(path, dfe.config)
}

func (dfe *DefaultSelector) IntArray(path string) ([]int, error) {
	return IntArray(path, dfe.config)
}

func (dfe *DefaultSelector) Float64Array(path string) ([]float64, error) {
	return Float64Array(path, dfe.config)
}

func (dfe *DefaultSelector) BoolVal(path string) (bool, error) {
	return BoolVal(path, dfe.config)
}

func (dfe *DefaultSelector) Config() ConfigNode {
	return dfe.config
}
