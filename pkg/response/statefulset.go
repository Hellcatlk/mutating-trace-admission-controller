package response

import (
	"encoding/json"
	"mutating-trace-admission-controller/pkg/util/patch"
	"mutating-trace-admission-controller/pkg/util/trace"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildStatefulSetPatch(raw []byte, op v1beta1.Operation) *v1beta1.AdmissionResponse {
	var statefulSet appv1.StatefulSet
	err := json.Unmarshal(raw, &statefulSet)
	if err != nil {
		glog.Errorf("unmarshal statefulset raw failed: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	// create or update span context
	spanContext, err := trace.DecodeSpanContext(statefulSet.GetAnnotations()[spanContextAnnotationKey])
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}
	spanContext = trace.StartSpan(spanContext)

	// create initial trace id
	inititalTraceID := ""
	if op == v1beta1.Create {
		if statefulSet.GetAnnotations()[initialTraceIDAnnotationKey] != "" {
			inititalTraceID = statefulSet.GetAnnotations()[initialTraceIDAnnotationKey]
		} else {
			inititalTraceID = spanContext.TraceID.String()
		}
	}

	// create patch annotations
	patchAnnotations, err := buildAnnotations(inititalTraceID, spanContext)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchBytes, err := patch.EncodePatch(patch.BuildAnnotationsPatch(statefulSet.Annotations, patchAnnotations))
	if err != nil {
		glog.Errorf("encode statefulset patch failed: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
