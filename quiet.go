package config_access

// QuietSelector does not return errors for missing paths or incompatible types
type QuietSelector interface {
	PathExists(path string) bool
	Value(path string) interface{}
	ObjectVal(path string) ConfigNode
	StringVal(path string) string
	StringArray(path string) []string
	IntArray(path string) []int
	Float64Array(path string) []float64
	IntVal(path string) int
	Float64Val(path string) float64
	Array(path string) []interface{}
	BoolVal(path string) bool
}

func NewDeferredErrorQuietSelector(conf Selector, errorFunc func(path string, err error)) QuietSelector {
	dqs := new(DeferredErrorQuietSelector)
	dqs.conf = conf
	dqs.handleError = errorFunc

	return dqs
}

// DeferredErrorQuietSelector does not return errors for missing paths or incompatible types, instead
// executes a supplied error handling function and returns the appropriate zero value for the requested config item.
type DeferredErrorQuietSelector struct {
	conf        Selector
	handleError func(path string, err error)
}

func (dqs *DeferredErrorQuietSelector) PathExists(path string) bool {

	return dqs.conf.PathExists(path)

}

func (dqs *DeferredErrorQuietSelector) Value(path string) interface{} {

	return dqs.conf.Value(path)

}

func (dqs *DeferredErrorQuietSelector) ObjectVal(path string) ConfigNode {

	if v, err := dqs.conf.ObjectVal(path); err != nil {
		dqs.handleError(path, err)
		return nil
	} else {
		return v
	}

}

func (dqs *DeferredErrorQuietSelector) StringVal(path string) string {

	if v, err := dqs.conf.StringVal(path); err != nil {
		dqs.handleError(path, err)
		return ""
	} else {
		return v
	}

}

func (dqs *DeferredErrorQuietSelector) IntVal(path string) int {

	if v, err := dqs.conf.IntVal(path); err != nil {
		dqs.handleError(path, err)
		return 0
	} else {
		return v
	}

}

func (dqs *DeferredErrorQuietSelector) Float64Val(path string) float64 {

	if v, err := dqs.conf.Float64Val(path); err != nil {
		dqs.handleError(path, err)
		return 0
	} else {
		return v
	}

}

func (dqs *DeferredErrorQuietSelector) Array(path string) []interface{} {

	if v, err := dqs.conf.Array(path); err != nil {
		dqs.handleError(path, err)
		return nil
	} else {
		return v
	}

}

func (dqs *DeferredErrorQuietSelector) StringArray(path string) []string {
	if v, err := dqs.conf.StringArray(path); err != nil {
		dqs.handleError(path, err)
		return nil
	} else {
		return v
	}
}

func (dqs *DeferredErrorQuietSelector) IntArray(path string) []int {
	if v, err := dqs.conf.IntArray(path); err != nil {
		dqs.handleError(path, err)
		return nil
	} else {
		return v
	}
}

func (dqs *DeferredErrorQuietSelector) Float64Array(path string) []float64 {
	if v, err := dqs.conf.Float64Array(path); err != nil {
		dqs.handleError(path, err)
		return nil
	} else {
		return v
	}
}

func (dqs *DeferredErrorQuietSelector) BoolVal(path string) bool {

	if v, err := dqs.conf.BoolVal(path); err != nil {
		dqs.handleError(path, err)
		return false
	} else {
		return v
	}

}

// QuietSelectorFromPathValues creates a new QuietSelector populated with a map of complete paths (e.g. "my.config.path": "value")
func QuietSelectorFromPathValues(pv map[string]interface{}, errorFunc func(path string, err error)) QuietSelector {
	return NewDeferredErrorQuietSelector(SelectorFromPathValues(pv), errorFunc)
}
