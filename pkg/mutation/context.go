package mutation

import (
	"github.com/pkg/errors"
	"github.com/riotkit-org/git-clone-operator/pkg/consts"
	corev1 "k8s.io/api/core/v1"
)

type Parameters struct {
	GitUrl      string
	GitUsername string
	GitToken    string
	GitRevision string
	FilesOwner  string
	FilesGroup  string
	TargetPath  string
	Image       string
}

func NewCheckoutParametersFromPod(pod *corev1.Pod, defaultImage string, defaultGitUsername string, defaultGitToken string, secretUsername string, secretGitToken string) (Parameters, error) {
	if _, exists := pod.Annotations[consts.AnnotationGitUrl]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in Pod, cannot recognize GIT url", consts.AnnotationGitUrl)
	}
	if _, exists := pod.Annotations[consts.AnnotationGitPath]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in Pod, cannot guess destination directory", consts.AnnotationGitPath)
	}
	if _, exists := pod.Annotations[consts.AnnotationFilesOwner]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in Pod, files owner id must be specified", consts.AnnotationFilesOwner)
	}
	if _, exists := pod.Annotations[consts.AnnotationFilesGroup]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in Pod, files owner group id must be specified", consts.AnnotationFilesOwner)
	}
	if _, exists := pod.Annotations[consts.AnnotationRev]; !exists {
		pod.Annotations[consts.AnnotationRev] = "main"
	}

	if secretUsername == "" {
		secretUsername = defaultGitUsername
	}
	if secretGitToken == "" {
		secretGitToken = defaultGitToken
	}
	return Parameters{
		Image:       defaultImage,
		GitUrl:      pod.Annotations[consts.AnnotationGitUrl],
		GitRevision: pod.Annotations[consts.AnnotationRev],
		GitUsername: defaultGitUsername,
		GitToken:    secretGitToken,
		TargetPath:  pod.Annotations[consts.AnnotationGitPath],
		FilesOwner:  pod.Annotations[consts.AnnotationFilesOwner],
		FilesGroup:  pod.Annotations[consts.AnnotationFilesGroup],
	}, nil
}
