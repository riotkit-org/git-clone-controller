package admission

import (
	"encoding/json"
	"fmt"
	"github.com/riotkit-org/git-clone-controller/pkg/context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// ResolvePod extracts a pod from an admission request
func ResolvePod(a MutationRequest) (*corev1.Pod, error) {
	if a.Request.Kind.Kind != "Pod" {
		return nil, fmt.Errorf("only pods are supported here, got request type: %v", a.Request.Kind.Kind)
	}

	if a.IsDebugLevel {
		logrus.Printf("Processing request: %v", string(a.Request.Object.Raw))
	}

	p := corev1.Pod{}
	if err := json.Unmarshal(a.Request.Object.Raw, &p); err != nil {
		return nil, err
	}

	// fix: Missing namespace in case of scoped call by controllers like ReplicaSet
	if p.ObjectMeta.Namespace == "" && a.Request.Namespace != "" {
		p.ObjectMeta.Namespace = a.Request.Namespace
	}

	return &p, nil
}

func isPodToBeProcessed(pod *corev1.Pod) bool {
	if val, exists := pod.Labels[context.LabelIsEnabled]; exists && val == "true" {
		return true
	}
	return false
}
