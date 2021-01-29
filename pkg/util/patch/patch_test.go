package patch

import (
	"reflect"
	"testing"
)

func TestWithAnnotations(t *testing.T) {
	cases := []struct {
		name     string
		new      map[string]string
		old      map[string]string
		expected []Operation
	}{
		{
			name: "old is nil",
			new: map[string]string{
				"k1": "1",
				"k2": "2",
			},
			expected: []Operation{
				{
					Op:   "replace",
					Path: "/metadata/annotations",
					Value: map[string]string{
						"k1": "1",
						"k2": "2",
					},
				},
			},
		},
		{
			name: "old is empty",
			old:  map[string]string{},
			new: map[string]string{
				"k1": "1",
				"k2": "2",
			},
			expected: []Operation{
				{
					Op:   "replace",
					Path: "/metadata/annotations",
					Value: map[string]string{
						"k1": "1",
						"k2": "2",
					},
				},
			},
		},
		{
			name: "old have same key with new",
			old: map[string]string{
				"k1": "0",
				"k2": "2",
				"k3": "3",
			},
			new: map[string]string{
				"k1": "1",
				"k2": "2",
			},
			expected: []Operation{
				{
					Op:   "replace",
					Path: "/metadata/annotations",
					Value: map[string]string{
						"k1": "1",
						"k2": "2",
						"k3": "3",
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := WithAnnotations(c.old, c.new)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
		})
	}
}
