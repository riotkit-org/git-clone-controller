package mutation

import (
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
)

// Mutator is a container for mutation
type Mutator struct {
	Logger *logrus.Entry
	params *Parameters
}

// NewMutator returns an initialised instance of Mutator
func NewMutator(logger *logrus.Entry, params *Parameters) *Mutator {
	return &Mutator{Logger: logger, params: params}
}

// MutatePodPatch returns a json patch containing all the mutations needed for
// a given pod
func (m *Mutator) MutatePodPatch(pod *corev1.Pod) ([]byte, error) {
	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else {
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", podName)

	mutator := initContainerInjector{
		Logger:      log,
		image:       m.params.Image,
		path:        m.params.TargetPath,
		rev:         m.params.GitRevision,
		gitUrl:      m.params.GitUrl,
		gitToken:    m.params.GitToken,
		gitUsername: m.params.GitUsername,
	}
	mutatedPod, mutateErr := mutator.Mutate(pod.DeepCopy())
	if mutateErr != nil {
		return nil, errors.Wrap(mutateErr, "Cannot mutate pod")
	}

	// generate json patch
	patch, err := jsondiff.Compare(pod, mutatedPod)
	if err != nil {
		return nil, err
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return patchb, nil
}
