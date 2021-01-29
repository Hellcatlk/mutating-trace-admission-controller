package response

import (
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1beta1"

	"mutating-trace-admission-controller/pkg/config"
	"mutating-trace-admission-controller/pkg/util/trace"

	apitrace "go.opentelemetry.io/otel/api/trace"
)

// Build the response to inject the trace context into received object
func Build(r *http.Request, ar *admissionv1.AdmissionReview) (response *admissionv1.AdmissionResponse) {
	fmt.Println("-------------------------------------")
	fmt.Println(r.Header)
	fmt.Println(ar.Request.Operation)
	fmt.Println(ar.Request.Kind.Kind)
	fmt.Println("-------------------------------------")

	spanContext := trace.SpanContextFromRequestHeader(r)
	// Build the annotations to patch
	newAnnotations, err := buildAnnotations(spanContext)
	if len(newAnnotations) == 0 || err != nil {
		return &admissionv1.AdmissionResponse{
			UID:     ar.Request.UID,
			Allowed: true,
		}
	}

	switch ar.Request.Kind.Kind {
	case "Deployment":
		response = buildDeploymentPatch(ar.Request.Object.Raw, newAnnotations)
	case "DeamonSet":
		response = buildDeamonSetPatch(ar.Request.Object.Raw, newAnnotations)
	case "StatefulSet":
		response = buildStatefulSetPatch(ar.Request.Object.Raw, newAnnotations)
	case "ReplicaSet":
		response = buildReplicaSetPatch(ar.Request.Object.Raw, newAnnotations)
	case "Pod":
		response = buildPodPatch(ar.Request.Object.Raw, newAnnotations)
	default:
		response = &admissionv1.AdmissionResponse{
			Allowed: true,
		}
	}
	response.UID = ar.Request.UID

	return
}

// buildAnnotations create a annotation with initTraceID and span
func buildAnnotations(spanContext apitrace.SpanContext) (map[string]string, error) {
	encodedSpanContext, err := trace.EncodedSpanContext(spanContext)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		config.Get().Trace.SpanContextAnnotationKey: encodedSpanContext,
	}, nil
}
