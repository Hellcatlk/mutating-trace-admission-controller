package response

import (
	"encoding/json"
	"mutating-trace-admission-controller/pkg/util/patch"

	"github.com/golang/glog"
	admissionv1 "k8s.io/api/admission/v1beta1"
	autosalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildScalePatch(raw []byte, newAnnotations map[string]string) *admissionv1.AdmissionResponse {
	var scale autosalingv1.Scale
	err := json.Unmarshal(raw, &scale)
	if err != nil {
		glog.Errorf("unmarshal scale raw failed: %v", err)
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchBytes, err := patch.Encode(patch.WithAnnotations(scale.Annotations, newAnnotations))
	if err != nil {
		glog.Errorf("encode scale patch failed: %v", err)
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
