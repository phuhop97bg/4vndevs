package json_debug

import (
	"encoding/json"
)

var (
	// empty
	emptyBytes = []byte("{}")
)

// JSONDebugData func
func JSONDebugData(message interface{}) []byte {
	if data, err := json.Marshal(message); err == nil {
		return data
	}
	return emptyBytes
}

// JSONDebugDataString func
func JSONDebugDataString(message interface{}) string {
	return string(JSONDebugData(message))
}
