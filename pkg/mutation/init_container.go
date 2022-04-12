package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const InitContainerName = "git-checkout"

// initContainerInjector is a container for the mutation injecting environment vars
type initContainerInjector struct {
	Logger   logrus.FieldLogger
	image    string
	path     string
	rev      string
	gitUrl   string
	gitToken string
	userName string
}

// Mutate returns a new mutated pod according to set env rules
func (se initContainerInjector) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation", "Mutating pod")
	mutatedPod := pod.DeepCopy()

	if hasGitInitContainer(pod) {
		se.Logger.Infof("Pod '%s' already has initContainer present", pod.Name)
		return mutatedPod, nil
	}

	injectInitContainer(mutatedPod, se.image, se.path, se.rev, se.gitUrl, se.gitToken, se.userName)
	return mutatedPod, nil
}

// injectInitContainer injects an initContainer
func injectInitContainer(pod *corev1.Pod, image string, path string, rev string, gitUrl string, gitToken string, userName string) {
	pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
		Name:       InitContainerName,
		Image:      image,
		Command:    []string{"/usr/bin/git-clone-operator"},
		Args:       []string{"--path", path, "--rev", rev, "--url", gitUrl, "--token", gitToken, "--user", userName},
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
