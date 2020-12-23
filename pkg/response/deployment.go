package response

import (
	"encoding/json"
	"mutating-trace-admission-controller/pkg/util/patch"

	"github.com/golang/glog"
	admissionv1 "k8s.io/api/admission/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildDeploymentPatch(raw []byte, newAnnotations map[string]string) *admissionv1.AdmissionResponse {
	var deployment appv1.Deployment
	err := json.Unmarshal(raw, &deployment)
	if err != nil {
		glog.Errorf("unmarshal deployment raw failed: %v", err)
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchBytes, err := patch.Encode(patch.WithAnnotations(deployment.Annotations, newAnnotations))
	if err != nil {
		glog.Errorf("encode deployment patch failed: %v", err)
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
