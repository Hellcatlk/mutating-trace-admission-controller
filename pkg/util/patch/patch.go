package patch

import (
	"encoding/json"
)

// Operation ...
type Operation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// WithAnnotations build patch by annotations
func WithAnnotations(old, new map[string]string) (patchs []Operation) {
	if old == nil {
		old = make(map[string]string)
	}
	patch := Operation{
		Op:    "replace",
		Path:  "/metadata/annotations",
		Value: old,
	}

	for key, value := range new {
		patch.Value.(map[string]string)[key] = value
	}
	patchs = append(patchs, patch)

	return
}

// Encode patch by json
func Encode(patch []Operation) ([]byte, error) {
	return json.Marshal(patch)
}
