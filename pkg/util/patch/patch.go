package patch

import "encoding/json"

// Operation ...
type Operation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// WithAnnotations build patch by annotations
func WithAnnotations(old, new map[string]string) (patch []Operation) {
	var (
		patchAdd Operation = Operation{
			Op:    "add",
			Path:  "/metadata/annotations",
			Value: make(map[string]string, 0),
		}
	)

	for key, value := range new {
		if old == nil {
			patchAdd.Value.(map[string]string)[key] = value
		} else if old[key] == "" {
			patch = append(patch, Operation{
				Op:    "add",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		} else if old[key] != value {
			patch = append(patch, Operation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}

	if len(patchAdd.Value.(map[string]string)) != 0 {
		patch = append(patch, patchAdd)
	}

	return
}

// Encode patch by json
func Encode(patch []Operation) ([]byte, error) {
	return json.Marshal(patch)
}
