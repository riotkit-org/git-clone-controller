package admission

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestReviewResponse(t *testing.T) {
	uid := types.UID("test")
	reason := "fail!"

	want := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: false,
			Result: &metav1.Status{
				Code:    418,
				Message: reason,
			},
		},
	}

	got := reviewResponse(uid, false, http.StatusTeapot, reason)
	assert.Equal(t, want, got)
}

func TestPatchReviewResponse(t *testing.T) {
	uid := types.UID("test")
	patchType := admissionv1.PatchTypeJSONPatch
	patch := []byte(`not quite a real patch`)

	want := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:       uid,
			Allowed:   true,
			PatchType: &patchType,
			Patch:     patch,
		},
	}

	got, err := patchReviewResponse(uid, patch)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}
