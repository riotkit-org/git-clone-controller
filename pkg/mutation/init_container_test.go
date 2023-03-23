package mutation_test

import (
	"github.com/riotkit-org/git-clone-controller/pkg/context"
	"github.com/riotkit-org/git-clone-controller/pkg/mutation"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"testing"
)

var exampleSpec = `apiVersion: v1
kind: Pod
metadata:
    name: "mutual-aid"
    namespace: anarchism
    # labels are provided with context.Parameters{}
spec:
    restartPolicy: Never
    automountServiceAccountToken: false
    containers:
        - command:
              - /bin/sh
              - "-c"
              - "find /workspace/source; ls -la /workspace/source"
          image: busybox:latest
          name: test
          volumeMounts:
              - mountPath: /workspace/source
                name: workspace
    volumes:
        - name: workspace
          emptyDir: {}
    securityContext:
        fsGroup: 1000
`

func TestMutatePodByInjectingInitContainer(t *testing.T) {
	examplePod := &corev1.Pod{}
	if err := yaml.Unmarshal([]byte(exampleSpec), &examplePod); err != nil {
		logrus.Fatal(err)
	}

	params := context.Parameters{
		GitUrl:      "https://github.com/riotkit-org/backup-repository",
		GitUsername: "",
		GitToken:    "",
		GitRevision: "main",
		FilesOwner:  "1000",
		FilesGroup:  "1001",
		TargetPath:  "/workspace/git",
		Image:       "ghcr.io/peter/kropotkin",
	}

	m, err := mutation.MutatePodByInjectingInitContainer(examplePod, &logrus.Logger{}, params)

	assert.Nil(t, err)
	assert.Len(t, m.Spec.InitContainers, 1, "Expected that one initContainer will be added")

	// basic things
	assert.Equal(t, "git-checkout", m.Spec.InitContainers[0].Name)
	assert.Equal(t, "ghcr.io/peter/kropotkin", m.Spec.InitContainers[0].Image)

	// this may fail time-to-time if commandline will be changed
	assert.Equal(t, []string{"checkout", "https://github.com/riotkit-org/backup-repository", "--path", "/workspace/git", "--rev", "main", "--token", "", "--username", "", "--clean-remotes"}, m.Spec.InitContainers[0].Args)

	// security context
	runAsRoot := true
	runAsUser := int64(1000)
	runAsGroup := int64(1001)
	assert.Equal(t, &runAsRoot, m.Spec.InitContainers[0].SecurityContext.RunAsNonRoot)
	assert.Equal(t, &runAsUser, m.Spec.InitContainers[0].SecurityContext.RunAsUser)
	assert.Equal(t, &runAsGroup, m.Spec.InitContainers[0].SecurityContext.RunAsGroup)
}

func TestMutatePodByInjectingInitContainer_WithoutSecurityContext(t *testing.T) {
	examplePod := &corev1.Pod{}
	if err := yaml.Unmarshal([]byte(exampleSpec), &examplePod); err != nil {
		logrus.Fatal(err)
	}

	params := context.Parameters{
		GitUrl:      "https://github.com/riotkit-org/backup-repository",
		GitUsername: "",
		GitToken:    "",
		GitRevision: "main",
		FilesOwner:  "",
		FilesGroup:  "",
		TargetPath:  "/workspace/git",
		Image:       "ghcr.io/peter/kropotkin",
	}

	m, err := mutation.MutatePodByInjectingInitContainer(examplePod, &logrus.Logger{}, params)

	assert.Nil(t, err)
	assert.Len(t, m.Spec.InitContainers, 1, "Expected that one initContainer will be added")

	// security context
	assert.Nil(t, m.Spec.InitContainers[0].SecurityContext)
}
