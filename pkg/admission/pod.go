package admission

import (
	"encoding/json"
	"fmt"
	"github.com/riotkit-org/git-clone-operator/pkg/context"
	corev1 "k8s.io/api/core/v1"
)

// ResolvePod extracts a pod from an admission request
func ResolvePod(a MutationRequest) (*corev1.Pod, error) {
	if a.Request.Kind.Kind != "ResolvePod" {
		return nil, fmt.Errorf("only pods are supported here")
	}

	p := corev1.Pod{}
	if err := json.Unmarshal(a.Request.Object.Raw, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func isPodToBeProcessed(pod *corev1.Pod) bool {
	if val, exists := pod.Labels[context.LabelIsEnabled]; exists && val == "true" {
		return true
	}
	return false
}
