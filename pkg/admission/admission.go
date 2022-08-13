// Package admission handles kubernetes admissions,
// it takes admission requests and returns admission reviews;
// for example, to mutate or validate pods
package admission

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	appContext "github.com/riotkit-org/git-clone-controller/pkg/context"
	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"net/http"

	"github.com/riotkit-org/git-clone-controller/pkg/mutation"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MutationRequest is a container for admission logic
type MutationRequest struct {
	Logger       *logrus.Entry
	IsDebugLevel bool
	Request      *admissionv1.AdmissionRequest

	DefaultImage       string
	DefaultGitUsername string
	DefaultGitToken    string

	Client kubernetes.Interface
}

// ProcessAdmissionRequest takes an admission request and mutates the pod within,
// it returns an admission review with mutations as a json patch (if any)
func (a MutationRequest) ProcessAdmissionRequest() (*admissionv1.AdmissionReview, error) {
	pod, err := ResolvePod(a)
	if err != nil {
		e := fmt.Sprintf("could not parse pod in admission review request: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	// validate
	if !isPodToBeProcessed(pod) {
		return reviewResponse(a.Request.UID, true, http.StatusOK, ""), nil
	}
	gitUserName, gitToken, secretErr := resolveSecretForPod(context.TODO(), a.Client, pod)
	if secretErr != nil {
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, errors.Wrap(secretErr, "git-clone-controller: Missing `kind: Secret` for annotated Pod").Error()), err
	}

	// glue parameters together
	parameters, paramsErr := appContext.NewCheckoutParametersFromPod(pod, a.DefaultImage, a.DefaultGitUsername, a.DefaultGitToken, gitUserName, gitToken)
	if paramsErr != nil {
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, errors.Wrap(paramsErr, "git-clone-controller: Cannot parse Pod labels/annotations").Error()), paramsErr
	}

	// create a patch
	patch, err := a.CreatePodPatch(pod, parameters)
	if err != nil {
		e := fmt.Sprintf("could not mutate pod: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	return patchReviewResponse(a.Request.UID, patch)
}

// CreatePodPatch returns a json patch containing all the mutations needed for
// a given pod
func (a MutationRequest) CreatePodPatch(pod *corev1.Pod, params appContext.Parameters) ([]byte, error) {
	var podName string
	if pod.ObjectMeta.Name != "" {
		podName = pod.ObjectMeta.Name
	} else {
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", podName)

	mutatedPod, mutateErr := mutation.MutatePodByInjectingInitContainer(pod.DeepCopy(), log, params)
	if mutateErr != nil {
		return nil, errors.Wrap(mutateErr, "Cannot mutate pod")
	}

	// generate json patch
	patch, err := jsondiff.Compare(pod, mutatedPod)
	if err != nil {
		return nil, err
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return patchBytes, nil
}

// reviewResponse sends a review response without returning a patch
func reviewResponse(uid types.UID, allowed bool, httpCode int32, reason string) *admissionv1.AdmissionReview {
	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: allowed,
			Result: &metav1.Status{
				Code:    httpCode,
				Message: reason,
			},
		},
	}
}

// patchReviewResponse builds an admission review with given json patch
func patchReviewResponse(uid types.UID, patch []byte) (*admissionv1.AdmissionReview, error) {
	patchType := admissionv1.PatchTypeJSONPatch

	return &admissionv1.AdmissionReview{
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
	}, nil
}
