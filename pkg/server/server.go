package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"mutating-trace-admission-controller/pkg/response"

	"github.com/golang/glog"
	admissionv1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionv1.AddToScheme(runtimeScheme)
	_ = v1.AddToScheme(runtimeScheme)
}

// WebhookServer is ...
type WebhookServer struct {
	Server *http.Server
}

// Serve http handler
func (whsvr *WebhookServer) Serve(w http.ResponseWriter, r *http.Request) {
	// Verify the content type is accurate
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("read request body failed: %v", err)
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	// Decode request body
	ar := admissionv1.AdmissionReview{}
	_, _, err = deserializer.Decode(body, nil, &ar)
	if err != nil {
		glog.Errorf("decode request body failed: %v", err)
		http.Error(w, fmt.Sprintf("could not decode response: %v", err), http.StatusBadRequest)
		return
	}

	// Build response
	admissionReview := admissionv1.AdmissionReview{}
	admissionReview.Response = response.Build(r, &ar)

	// Marshal respson
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		glog.Errorf("marshal respson failed: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}

	// Return respson
	_, err = w.Write(resp)
	if err != nil {
		glog.Errorf("write respson failed: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}

	return
}
