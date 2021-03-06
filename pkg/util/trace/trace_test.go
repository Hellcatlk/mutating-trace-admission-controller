package trace

import (
	"reflect"
	"testing"

	apitrace "go.opentelemetry.io/otel/api/trace"
)

func TestEncodeSpanContext(t *testing.T) {
	cases := []struct {
		name          string
		SpanContext   apitrace.SpanContext
		expected      string
		expectedError bool
	}{
		{
			name: "normal",
			SpanContext: apitrace.SpanContext{
				TraceID:    [16]byte{1, 2, 3},
				SpanID:     [8]byte{4, 5},
				TraceFlags: 1,
			},
			expected: "AQIDAAAAAAAAAAAAAAAAAAQFAAAAAAAAAQ==",
		},
		{
			name:          "empty",
			SpanContext:   apitrace.EmptySpanContext(),
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := EncodedSpanContext(c.SpanContext)
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("expected: %+v,got: %+v", c.expected, got)
			}
		})
	}
}

func TestDecodeSpanContext(t *testing.T) {
	cases := []struct {
		encodedSpanContext string
		expected           apitrace.SpanContext
	}{

		{
			encodedSpanContext: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
			expected:           apitrace.EmptySpanContext(),
		},
		{
			encodedSpanContext: "AQIDAAAAAAAAAAAAAAAAAAQFAAAAAAAAAQ==",
			expected: apitrace.SpanContext{
				TraceID:    [16]byte{1, 2, 3},
				SpanID:     [8]byte{4, 5},
				TraceFlags: 1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.encodedSpanContext, func(t *testing.T) {
			got, err := DecodeSpanContext(c.encodedSpanContext)
			if err != nil {
				t.Errorf("got unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("expected: %+v,got: %+v", c.expected, got)
			}
		})
	}
}
