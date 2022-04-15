// Package admission handles kubernetes admissions,
// it takes admission requests and returns admission reviews;
// for example, to mutate or validate pods
package admission

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/riotkit-org/git-clone-operator/pkg/consts"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"net/http"

	"github.com/riotkit-org/git-clone-operator/pkg/mutation"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProcessingService is a container for admission business
type ProcessingService struct {
	Logger  *logrus.Entry
	Request *admissionv1.AdmissionRequest

	DefaultImage       string
	DefaultGitUsername string
	DefaultGitToken    string

	Client *kubernetes.Clientset
}

// MutatePodReview takes an admission request and mutates the pod within,
// it returns an admission review with mutations as a json patch (if any)
func (a ProcessingService) MutatePodReview() (*admissionv1.AdmissionReview, error) {
	pod, err := a.Pod()
	if err != nil {
		e := fmt.Sprintf("could not parse pod in admission review request: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	if !a.isPodToBeProcessed(pod) {
		return patchReviewResponse(a.Request.UID, []byte("{}"))
	}
	gitUserName, gitToken, secretErr := a.resolveSecretForPod(context.TODO(), pod)
	if secretErr != nil {
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, errors.Wrap(secretErr, "git-clone-operator: Missing `kind: Secret` for annotated Pod").Error()), err
	}

	parameters, paramsErr := mutation.NewCheckoutParametersFromPod(pod, a.DefaultImage, a.DefaultGitUsername, a.DefaultGitToken, gitUserName, gitToken)
	if paramsErr != nil {
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, errors.Wrap(secretErr, "git-clone-operator: Cannot parse Pod labels/annotations").Error()), err
	}

	m := mutation.NewMutator(a.Logger, &parameters)
	patch, err := m.MutatePodPatch(pod)
	if err != nil {
		e := fmt.Sprintf("could not mutate pod: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	return patchReviewResponse(a.Request.UID, patch)
}

func (a ProcessingService) isPodToBeProcessed(pod *corev1.Pod) bool {
	if val, exists := pod.Labels[consts.LabelIsEnabled]; exists && val == "true" {
		return true
	}
	return false
}

// Pod extracts a pod from an admission request
func (a ProcessingService) Pod() (*corev1.Pod, error) {
	if a.Request.Kind.Kind != "Pod" {
		return nil, fmt.Errorf("only pods are supported here")
	}

	p := corev1.Pod{}
	if err := json.Unmarshal(a.Request.Object.Raw, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// resolveSecretForPod Finds a `kind: Secret` using information from Pod's annotations and extracts secrets from that secret by specified keys
func (a ProcessingService) resolveSecretForPod(ctx context.Context, pod *corev1.Pod) (string, string, error) {
	// checking required annotations
	if val, exists := pod.Labels[consts.AnnotationSecretName]; !exists || val == "" {
		logrus.Infof("No label '%s' defined for Pod '%s/%s'", consts.AnnotationSecretName, pod.Namespace, pod.Name)
		return "", "", nil
	}
	if val, exists := pod.Labels[consts.AnnotationSecretTokenKey]; !exists || val == "" {
		logrus.Infof("No label '%s' defined for Pod '%s/%s'", consts.AnnotationSecretTokenKey, pod.Namespace, pod.Name)
		return "", "", nil
	}
	if val, exists := pod.Labels[consts.AnnotationSecretUserKey]; !exists || val == "" {
		logrus.Infof("No label '%s' defined for Pod '%s/%s'", consts.AnnotationSecretUserKey, pod.Namespace, pod.Name)
		return "", "", nil
	}

	// fetching `kind: Secret` from API
	secretName := pod.Labels[consts.AnnotationSecretName]
	secret, err := a.Client.CoreV1().Secrets(pod.Namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return "", "", errors.Wrapf(err, "Cannot fetch secret for Pod annotated with '%s=%s'", consts.AnnotationSecretName, pod.Labels[consts.AnnotationSecretName])
	}

	// extracting data from `kind: Secret`
	token, tokenDefined := secret.Data[consts.AnnotationSecretTokenKey]
	if !tokenDefined {
		return "", "", errors.Errorf("The secret '%s' does not contain key '%s'", secretName, secret.Data[consts.AnnotationSecretTokenKey])
	}
	username, usernameDefined := secret.Data[consts.AnnotationSecretTokenKey]
	if !usernameDefined {
		return "", "", errors.Errorf("The secret '%s' does not contain key '%s'", secretName, secret.Data[consts.AnnotationSecretUserKey])
	}

	return string(username), string(token), nil
}

// reviewResponse TODO: godoc
func reviewResponse(uid types.UID, allowed bool, httpCode int32,
	reason string) *admissionv1.AdmissionReview {
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
