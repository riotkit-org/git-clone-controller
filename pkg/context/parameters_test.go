package context_test

import (
	"github.com/riotkit-org/git-clone-operator/pkg/context"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

type parameterTestVariant struct {
	expectedErr   string
	expectedUser  string
	expectedToken string

	annotations map[string]string

	defaultImage       string
	defaultGitUsername string
	defaultGitToken    string
	secretUsername     string
	secretGitToken     string
}

func TestNewCheckoutParametersFromPod(t *testing.T) {
	var variants = []parameterTestVariant{
		// Successful case
		{
			expectedErr:   "",
			expectedUser:  "custom-user",
			expectedToken: "custom-token",

			annotations: map[string]string{
				"git-clone-operator/revision":   "main",
				"git-clone-operator/url":        "https://github.com/jenkins-x/go-scm",
				"git-clone-operator/path":       "/workspace/source",
				"git-clone-operator/owner":      "1000",
				"git-clone-operator/group":      "1000",
				"git-clone-operator/secretName": "git-secrets",
				"git-clone-operator/secretKey":  "jenkins-x",
			},

			defaultImage:       "ghcr.io/riotkit-org/backup-repository",
			defaultGitUsername: "default-user",
			defaultGitToken:    "default-token",

			secretUsername: "custom-user",
			secretGitToken: "custom-token",
		},

		// Missing url
		{
			expectedErr:   "cannot recognize GIT url",
			expectedUser:  "custom-user",
			expectedToken: "custom-token",

			annotations: map[string]string{
				"git-clone-operator/revision":   "main",
				"git-clone-operator/url":        "", // TEST: MISSING
				"git-clone-operator/path":       "/workspace/source",
				"git-clone-operator/owner":      "1000",
				"git-clone-operator/group":      "1000",
				"git-clone-operator/secretName": "git-secrets",
				"git-clone-operator/secretKey":  "jenkins-x",
			},

			defaultImage:       "ghcr.io/riotkit-org/backup-repository",
			defaultGitUsername: "default-user",
			defaultGitToken:    "default-token",

			secretUsername: "custom-user",
			secretGitToken: "custom-token",
		},

		// Missing PATH
		{
			expectedErr:   "cannot guess destination directory",
			expectedUser:  "custom-user",
			expectedToken: "custom-token",

			annotations: map[string]string{
				"git-clone-operator/revision": "main",
				"git-clone-operator/url":      "https://github.com/jenkins-x/go-scm",
				// path is missing
				"git-clone-operator/owner":      "1000",
				"git-clone-operator/group":      "1000",
				"git-clone-operator/secretName": "git-secrets",
				"git-clone-operator/secretKey":  "jenkins-x",
			},

			defaultImage:       "ghcr.io/riotkit-org/backup-repository",
			defaultGitUsername: "default-user",
			defaultGitToken:    "default-token",

			secretUsername: "custom-user",
			secretGitToken: "custom-token",
		},

		// Owner id is missing
		{
			expectedErr:   "files owner id must be specified",
			expectedUser:  "custom-user",
			expectedToken: "custom-token",

			annotations: map[string]string{
				"git-clone-operator/revision":   "main",
				"git-clone-operator/url":        "https://github.com/jenkins-x/go-scm",
				"git-clone-operator/path":       "/workspace/source",
				"git-clone-operator/owner":      "", // MISSING
				"git-clone-operator/group":      "1000",
				"git-clone-operator/secretName": "git-secrets",
				"git-clone-operator/secretKey":  "jenkins-x",
			},

			defaultImage:       "ghcr.io/riotkit-org/backup-repository",
			defaultGitUsername: "default-user",
			defaultGitToken:    "default-token",

			secretUsername: "custom-user",
			secretGitToken: "custom-token",
		},

		// Missing group id
		{
			expectedErr:   "files owner group id must be specified",
			expectedUser:  "custom-user",
			expectedToken: "custom-token",

			annotations: map[string]string{
				"git-clone-operator/revision": "main",
				"git-clone-operator/url":      "https://github.com/jenkins-x/go-scm",
				"git-clone-operator/path":     "/workspace/source",
				"git-clone-operator/owner":    "1000",
				// group is missing
				"git-clone-operator/secretName": "git-secrets",
				"git-clone-operator/secretKey":  "jenkins-x",
			},

			defaultImage:       "ghcr.io/riotkit-org/backup-repository",
			defaultGitUsername: "default-user",
			defaultGitToken:    "default-token",

			secretUsername: "custom-user",
			secretGitToken: "custom-token",
		},
	}

	for _, variant := range variants {
		pod := v1.Pod{}
		pod.SetAnnotations(variant.annotations)

		params, err := context.NewCheckoutParametersFromPod(&pod, variant.defaultImage, variant.defaultGitUsername, variant.defaultGitToken, variant.secretUsername, variant.secretGitToken)

		if variant.expectedErr == "" {
			assert.Nil(t, err)
			assert.Equal(t, variant.expectedUser, params.GitUsername)
			assert.Equal(t, variant.expectedToken, params.GitToken)
		} else {
			assert.Contains(t, err.Error(), variant.expectedErr)
		}
	}
}
