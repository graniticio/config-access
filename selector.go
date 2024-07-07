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
	Value(path string, o ...Opts) interface{}
	ObjectVal(path string, o ...Opts) (ConfigNode, error)
	StringVal(path string, o ...Opts) (string, error)
	IntVal(path string, o ...Opts) (int, error)
	Float64Val(path string, o ...Opts) (float64, error)
	Array(path string, o ...Opts) ([]interface{}, error)
	StringArray(path string, o ...Opts) ([]string, error)
	IntArray(path string, o ...Opts) ([]int, error)
	Float64Array(path string, o ...Opts) ([]float64, error)
	BoolVal(path string, o ...Opts) (bool, error)
	Flush()
	Config() ConfigNode
}

// Opts defines optional behaviour for accessing and interpreting config values
type Opts struct {
	// If this value is set, it will be returned instead of an error if there is no value at the requested path
	OnMissing any
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

func (dfe *DefaultSelector) Value(path string, o ...Opts) interface{} {
	if v := Value(path, dfe.config); v != nil {
		return v
	} else {
		opts := options(o)
		return opts.OnMissing
	}
}

func (dfe *DefaultSelector) ObjectVal(path string, o ...Opts) (ConfigNode, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.(ConfigNode), nil
	}

	return ObjectVal(path, dfe.config, dfe.errorOnMissingObjectPath)
}

func (dfe *DefaultSelector) StringVal(path string, o ...Opts) (string, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.(string), nil
	}

	return StringVal(path, dfe.config)
}

func (dfe *DefaultSelector) IntVal(path string, o ...Opts) (int, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.(int), nil
	}

	return IntVal(path, dfe.config)
}

func (dfe *DefaultSelector) Float64Val(path string, o ...Opts) (float64, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.(float64), nil
	}

	return Float64Val(path, dfe.config)
}

func (dfe *DefaultSelector) Array(path string, o ...Opts) ([]interface{}, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.([]interface{}), nil
	}

	return Array(path, dfe.config, dfe.errorOnMissingArrayPath)
}

func (dfe *DefaultSelector) StringArray(path string, o ...Opts) ([]string, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.([]string), nil
	}

	return StringArray(path, dfe.config)
}

func (dfe *DefaultSelector) IntArray(path string, o ...Opts) ([]int, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.([]int), nil
	}

	return IntArray(path, dfe.config)
}

func (dfe *DefaultSelector) Float64Array(path string, o ...Opts) ([]float64, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.([]float64), nil
	}

	return Float64Array(path, dfe.config)
}

func (dfe *DefaultSelector) BoolVal(path string, o ...Opts) (bool, error) {

	opts := options(o)

	if opts.OnMissing != nil && !PathExists(path, dfe.config) {
		return opts.OnMissing.(bool), nil
	}

	return BoolVal(path, dfe.config)
}

func (dfe *DefaultSelector) Config() ConfigNode {
	return dfe.config
}

func options(o []Opts) Opts {
	if len(o) == 0 {
		return Opts{}
	} else {
		return o[0]
	}
}
