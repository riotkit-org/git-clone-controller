package mutation

import (
	appCtx "github.com/riotkit-org/git-clone-operator/pkg/context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const InitContainerName = "git-checkout"

// MutatePodByInjectingInitContainer returns a new mutated pod according to set env rules
func MutatePodByInjectingInitContainer(pod *corev1.Pod, logger logrus.FieldLogger, params appCtx.Parameters) (*corev1.Pod, error) {
	nLogger := logger.WithField("mutation", "Mutating pod")
	mutatedPod := pod.DeepCopy()

	if hasGitInitContainer(pod) {
		nLogger.Infof("ResolvePod '%s' already has initContainer present", pod.Name)
		return mutatedPod, nil
	}

	injectInitContainer(mutatedPod, params.Image, params.TargetPath, params.GitRevision, params.GitUrl, params.GitToken, params.GitUsername)
	return mutatedPod, nil
}

// injectInitContainer injects an initContainer
func injectInitContainer(pod *corev1.Pod, image string, path string, rev string, gitUrl string, gitToken string, userName string) {
	pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
		Name:       InitContainerName,
		Image:      image,
		Command:    []string{"/usr/bin/git-clone-operator"},
		Args:       []string{"checkout", "--path", path, "--rev", rev, "--url", gitUrl, "--token", gitToken, "--user", userName},
		WorkingDir: path,
		//EnvFrom:    nil,
		//VolumeMounts:             nil,
		//VolumeDevices:            nil,
		ImagePullPolicy: "Always",
	})
}

func hasGitInitContainer(pod *corev1.Pod) bool {
	for _, container := range pod.Spec.InitContainers {
		if container.Name == InitContainerName {
			return true
		}
	}
	return false
}
