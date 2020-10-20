package trace

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagators"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var propagator otel.TextMapPropagator

const initialTraceIDBaggageKey label.Key = "Initial-Trace-Id"

func init() {
	propagator = otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{})
}

// InitTracer ...
func InitTracer() func() {
	var err error
	exp, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		glog.Fatalf("failed to initialize stdout exporter %v\n", err)
		return nil
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(
			sdktrace.Config{
				DefaultSampler: sdktrace.NeverSample(),
			},
		),
		sdktrace.WithSpanProcessor(bsp),
	)
	global.SetTracerProvider(tp)
	return bsp.Shutdown
}

// StartSpan ...
func StartSpan(spanContext apitrace.SpanContext) apitrace.SpanContext {
	// update span
	ctx := apitrace.ContextWithSpan(
		context.Background(),
		traceSpan{
			spanContext: spanContext,
		},
	)
	tracer := global.Tracer("Log enhancement")
	_, span := tracer.Start(ctx, "")
	return span.SpanContext()
}

// EncodedSpanContext encode span to string
func EncodedSpanContext(spanContext apitrace.SpanContext) (string, error) {
	if reflect.DeepEqual(spanContext, apitrace.SpanContext{}) {
		return "", fmt.Errorf("span context is nil")
	}
	// encode to byte
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, spanContext)
	if err != nil {
		return "", err
	}
	// encode to string
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

// DecodeSpanContext decode encodedSpanContext to spanContext
func DecodeSpanContext(encodedSpanContext string) (apitrace.SpanContext, error) {
	if encodedSpanContext == "" {
		return apitrace.EmptySpanContext(), nil
	}
	// decode to byte
	byteList := make([]byte, base64.StdEncoding.DecodedLen(len(encodedSpanContext)))
	l, err := base64.StdEncoding.Decode(byteList, []byte(encodedSpanContext))
	if err != nil {
		return apitrace.EmptySpanContext(), err
	}
	byteList = byteList[:l]
	// decode to span context
	buffer := bytes.NewBuffer(byteList)
	spanContext := apitrace.SpanContext{}
	err = binary.Read(buffer, binary.LittleEndian, &spanContext)
	if err != nil {
		return apitrace.EmptySpanContext(), err
	}
	return spanContext, nil
}
