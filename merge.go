package config_access

const (
	ConfigUnset   = -2
	ConfigUnknown = -1
	ConfigString  = 1
	ConfigArray   = 2
	ConfigMap     = 3
	ConfigBool    = 4
)

// JSONType determines the apparent JSONType of the supplied Go interface.
func ConfigType(value interface{}) int {

	switch value.(type) {
	case string:
		return ConfigString
	case map[string]interface{}:
		return ConfigMap
	case bool:
		return ConfigBool
	case []interface{}:
		return ConfigArray
	default:
		return ConfigUnknown
	}
}

type ConfigMerger interface {
	Merge(base, additional ConfigNode) ConfigNode
}

func Merge(base, additional map[string]interface{}, mergeArrays bool) map[string]interface{} {

	for key, value := range additional {

		if existingEntry, ok := base[key]; ok {

			existingEntryType := ConfigType(existingEntry)
			newEntryType := ConfigType(value)

			if existingEntryType == ConfigMap && newEntryType == ConfigMap {
				Merge(existingEntry.(ConfigNode), value.(ConfigNode), mergeArrays)
			} else if mergeArrays && existingEntryType == ConfigArray && newEntryType == ConfigArray {
				base[key] = MergeArrays(existingEntry.([]interface{}), value.([]interface{}))
			} else {
				base[key] = value
			}
		} else {
			base[key] = value
		}

	}

	return base
}

func MergeArrays(a []interface{}, b []interface{}) []interface{} {
	return append(a, b...)
}
