package response

import (
	"encoding/json"
	"mutating-trace-admission-controller/pkg/util/patch"

	"github.com/golang/glog"
	admissionv1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildPodPatch(raw []byte, newAnnotations map[string]string) *admissionv1.AdmissionResponse {
	var pod corev1.Pod
	err := json.Unmarshal(raw, &pod)
	if err != nil {
		glog.Errorf("unmarshal pod raw failed: %v", err)
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchBytes, err := patch.Encode(patch.WithAnnotations(pod.GetAnnotations(), newAnnotations))
	if err != nil {
		glog.Errorf("encode pod patch failed: %v", err)
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
