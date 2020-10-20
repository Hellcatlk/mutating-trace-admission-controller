package response

import (
	"net/http"

	"k8s.io/api/admission/v1beta1"

	"mutating-trace-admission-controller/pkg/util/trace"

	apitrace "go.opentelemetry.io/otel/api/trace"
)

// avoid use char `/` in string
const initialTraceIDAnnotationKey string = "trace.kubernetes.io.initial"

// avoid use char `/` in string
const spanContextAnnotationKey string = "trace.kubernetes.io.span.context"

// Build build the response to inject the trace context into received object
func Build(r *http.Request, ar *v1beta1.AdmissionReview) (response *v1beta1.AdmissionResponse) {
	switch ar.Request.Kind.Kind {
	case "Deployment":
		response = buildDeploymentPatch(ar.Request.Object.Raw, ar.Request.Operation)
	case "DeamonSet":
		response = buildDeamonSetPatch(ar.Request.Object.Raw, ar.Request.Operation)
	case "StatefulSet":
		response = buildStatefulSetPatch(ar.Request.Object.Raw, ar.Request.Operation)
	case "ReplicaSet":
		response = buildReplicaSetPatch(ar.Request.Object.Raw, ar.Request.Operation)
	case "Pod":
		response = buildPodPatch(ar.Request.Object.Raw, ar.Request.Operation)
	default:
		response = &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}
	response.UID = ar.Request.UID

	return
}

// buildAnnotations create a annotation with initTraceID and span
func buildAnnotations(initTraceID string, spanContext apitrace.SpanContext) (map[string]string, error) {
	encodedSpanContext, err := trace.EncodedSpanContext(spanContext)
	if err != nil {
		return nil, err
	}
	if initTraceID == "" {
		return map[string]string{
			spanContextAnnotationKey: encodedSpanContext,
		}, nil
	}
	return map[string]string{
		initialTraceIDAnnotationKey: initTraceID,
		spanContextAnnotationKey:    encodedSpanContext,
	}, nil
}
