package response

import (
	"mutating-trace-admission-controller/pkg/config"
	"reflect"
	"testing"

	apitrace "go.opentelemetry.io/otel/api/trace"
)

func TestBuildAnnotations(t *testing.T) {
	cases := []struct {
		name          string
		spanContext   apitrace.SpanContext
		expected      map[string]string
		expectedError bool
	}{
		{
			name: "both",
			spanContext: apitrace.SpanContext{
				TraceID:    [16]byte{1, 2, 3},
				SpanID:     [8]byte{4, 5},
				TraceFlags: 1,
			},
			expected: map[string]string{
				config.Get().Trace.SpanContextAnnotationKey: "AQIDAAAAAAAAAAAAAAAAAAQFAAAAAAAAAQ==",
			},
		},
		{
			name:          "only init trace id",
			expectedError: true,
		},
		{
			name: "only span context",
			spanContext: apitrace.SpanContext{
				TraceID:    [16]byte{1, 2, 3},
				SpanID:     [8]byte{4, 5},
				TraceFlags: 1,
			},
			expected: map[string]string{
				config.Get().Trace.SpanContextAnnotationKey: "AQIDAAAAAAAAAAAAAAAAAAQFAAAAAAAAAQ==",
			},
		},
		{
			name:          "empty",
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := buildAnnotations(c.spanContext)
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
		})
	}
}
