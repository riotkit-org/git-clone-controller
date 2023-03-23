package mutation

import (
	appCtx "github.com/riotkit-org/git-clone-controller/pkg/context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

const InitContainerName = "git-checkout"

// MutatePodByInjectingInitContainer returns a new mutated pod according to set env rules
func MutatePodByInjectingInitContainer(pod *corev1.Pod, logger logrus.FieldLogger, params appCtx.Parameters) (*corev1.Pod, error) {
	nLogger := logger.WithField("mutation", "Mutating pod")
	mutatedPod := pod.DeepCopy()

	if hasGitInitContainer(pod) {
		nLogger.Infof("ResolvePod '%s' already has initContainer present", pod.ObjectMeta.Name)
		return mutatedPod, nil
	}

	injectInitContainer(mutatedPod, params.Image, params.TargetPath, params.GitRevision, params.GitUrl, params.GitToken, params.GitUsername, params.FilesOwner, params.FilesGroup, params.CleanUpWorkspace)
	return mutatedPod, nil
}

// injectInitContainer injects an initContainer
func injectInitContainer(pod *corev1.Pod, image string, path string, rev string, gitUrl string, gitToken string,
	userName string, owner string, group string, cleanUpWorkspace bool) {

	args := []string{
		"checkout",
		gitUrl,
		"--path", path,
		"--rev", rev,
		"--token", gitToken,
		"--username", userName,
	}

	if cleanUpWorkspace {
		args = append(args, "--clean-workspace")
	}

	container := corev1.Container{
		Name:       InitContainerName,
		Image:      image,
		Command:    []string{"/usr/bin/git-clone-controller"},
		Args:       args,
		WorkingDir: "/",
		// EnvFrom:    nil,
		VolumeMounts: mergeVolumeMounts(pod.Spec.Containers, path),
		// VolumeDevices:            nil,
		ImagePullPolicy: "Always",
	}

	// run container as specified user to operate on volume with given permissions
	if owner != "" && group != "" {
		logrus.Infof("Using UID=%v, GID=%v", owner, group)

		// RunAsNonRoot
		iOwner, _ := strconv.Atoi(owner)
		asNonRoot := iOwner > 0

		// RunAsUser
		iUser, _ := strconv.Atoi(owner)
		runAsUser := int64(iUser)

		// RunAsGroup
		iGroup, _ := strconv.Atoi(group)
		runAsGroup := int64(iGroup)

		// ReadOnlyRootFilesystem
		roFilesystem := false

		container.SecurityContext = &corev1.SecurityContext{
			RunAsUser:              &runAsUser,
			RunAsGroup:             &runAsGroup,
			RunAsNonRoot:           &asNonRoot,
			ReadOnlyRootFilesystem: &roFilesystem,
		}
	}

	pod.Spec.InitContainers = append(pod.Spec.InitContainers, container)
}

// mergeVolumeMounts merges volume mounts of multiple containers
func mergeVolumeMounts(containers []corev1.Container, targetPath string) []corev1.VolumeMount {
	var merged []corev1.VolumeMount
	var appendedPaths []string

	for _, container := range containers {
		for _, volume := range container.VolumeMounts {
			for _, existingPath := range appendedPaths {
				// already collected
				if existingPath == volume.MountPath {
					continue
				}
			}

			// do not collect non-related mount points at all
			if !strings.HasPrefix(targetPath, volume.MountPath) {
				continue
			}

			logrus.Infof("Collecting VolumeMount: %v", volume.String())
			appendedPaths = append(appendedPaths, volume.MountPath)
			merged = append(merged, volume)
		}
	}

	return merged
}

func hasGitInitContainer(pod *corev1.Pod) bool {
	for _, container := range pod.Spec.InitContainers {
		if container.Name == InitContainerName {
			return true
		}
	}
	return false
}
