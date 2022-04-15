package context

import (
	"github.com/pkg/errors"
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
	if _, exists := pod.Annotations[AnnotationGitUrl]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in ResolvePod, cannot recognize GIT url", AnnotationGitUrl)
	}
	if _, exists := pod.Annotations[AnnotationGitPath]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in ResolvePod, cannot guess destination directory", AnnotationGitPath)
	}
	if _, exists := pod.Annotations[AnnotationFilesOwner]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in ResolvePod, files owner id must be specified", AnnotationFilesOwner)
	}
	if _, exists := pod.Annotations[AnnotationFilesGroup]; !exists {
		return Parameters{}, errors.Errorf("Label '%s' not found in ResolvePod, files owner group id must be specified", AnnotationFilesOwner)
	}
	if _, exists := pod.Annotations[AnnotationRev]; !exists {
		pod.Annotations[AnnotationRev] = "main"
	}

	if secretUsername == "" {
		secretUsername = defaultGitUsername
	}
	if secretGitToken == "" {
		secretGitToken = defaultGitToken
	}
	return Parameters{
		Image:       defaultImage,
		GitUrl:      pod.Annotations[AnnotationGitUrl],
		GitRevision: pod.Annotations[AnnotationRev],
		GitUsername: defaultGitUsername,
		GitToken:    secretGitToken,
		TargetPath:  pod.Annotations[AnnotationGitPath],
		FilesOwner:  pod.Annotations[AnnotationFilesOwner],
		FilesGroup:  pod.Annotations[AnnotationFilesGroup],
	}, nil
}
