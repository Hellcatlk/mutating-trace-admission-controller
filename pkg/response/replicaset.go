package response

import (
	"encoding/json"
	"mutating-trace-admission-controller/pkg/util/patch"

	"github.com/golang/glog"
	admissionv1 "k8s.io/api/admission/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildReplicaSetPatch(raw []byte, newAnnotations map[string]string) *admissionv1.AdmissionResponse {
	var replicaSet appv1.ReplicaSet
	err := json.Unmarshal(raw, &replicaSet)
	if err != nil {
		glog.Errorf("unmarshal replicaset raw failed: %v", err)
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	// FIXME: use temporary measures to avoid bugs(the infinite loop of replicaset when update deployment)
	if replicaSet.OwnerReferences != nil {
		return &admissionv1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := patch.Encode(patch.WithAnnotations(replicaSet.Annotations, newAnnotations))
	if err != nil {
		glog.Errorf("encode replicaset patch failed: %v", err)
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
