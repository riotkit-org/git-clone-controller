package admission

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestResolvingWithAllValidFields(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "my-secret-name", Namespace: "default"},
		Data: map[string][]byte{
			"username": []byte("hello"),
			"password": []byte("riotkit"),
		},
		Type: "opaque",
	})

	pod := corev1.Pod{}
	pod.Namespace = "default"
	pod.Annotations = map[string]string{
		"git-clone-operator/secretName": "my-secret-name",
	}
	pod.Annotations["git-clone-operator/secretTokenKey"] = "password"
	pod.Annotations["git-clone-operator/secretUsernameKey"] = "username"
	pod.Annotations["git-clone-operator/revision"] = "HEAD"
	pod.Annotations["git-clone-operator/group"] = "161"
	pod.Annotations["git-clone-operator/owner"] = "161"
	pod.Annotations["git-clone-operator/path"] = "/var/www/riotkit"
	pod.Annotations["git-clone-operator/url"] = "https://github.com/riotkit-org/git-clone-operator"

	returnedUsername, returnedPassword, err := resolveSecretForPod(context.TODO(), client, &pod)

	assert.Nil(t, err)
	assert.Equal(t, "hello", returnedUsername)
	assert.Equal(t, "riotkit", returnedPassword)
}

func TestResolvingWithMissingPasswordInKindSecret(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "my-secret-name", Namespace: "default"},
		Data: map[string][]byte{
			"username": []byte("hello"),
			// "password": []byte("riotkit"),  // this one is MISSING
		},
		Type: "opaque",
	})

	pod := corev1.Pod{}
	pod.Namespace = "default"
	pod.Annotations = map[string]string{
		"git-clone-operator/secretName": "my-secret-name",
	}
	pod.Annotations["git-clone-operator/secretTokenKey"] = "password"
	pod.Annotations["git-clone-operator/secretUsernameKey"] = "username"
	pod.Annotations["git-clone-operator/revision"] = "HEAD"
	pod.Annotations["git-clone-operator/group"] = "161"
	pod.Annotations["git-clone-operator/owner"] = "161"
	pod.Annotations["git-clone-operator/path"] = "/var/www/riotkit"
	pod.Annotations["git-clone-operator/url"] = "https://github.com/riotkit-org/git-clone-operator"

	_, _, err := resolveSecretForPod(context.TODO(), client, &pod)

	assert.Equal(t, "The secret 'my-secret-name' does not contain key 'password'", err.Error())
}

func TestResolvingWithMissingUsernameInKindSecret(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "my-secret-name", Namespace: "default"},
		Data: map[string][]byte{
			// "username": []byte("hello"),  // this one was defined in annotations, but missing in `kind: Secret`
			"password": []byte("riotkit"),
		},
		Type: "opaque",
	})

	pod := corev1.Pod{}
	pod.Namespace = "default"
	pod.Annotations = map[string]string{
		"git-clone-operator/secretName": "my-secret-name",
	}
	pod.Annotations["git-clone-operator/secretTokenKey"] = "password"
	pod.Annotations["git-clone-operator/secretUsernameKey"] = "username" // here we make the "username" field in `kind: Secret` mandatory
	pod.Annotations["git-clone-operator/revision"] = "HEAD"
	pod.Annotations["git-clone-operator/group"] = "161"
	pod.Annotations["git-clone-operator/owner"] = "161"
	pod.Annotations["git-clone-operator/path"] = "/var/www/riotkit"
	pod.Annotations["git-clone-operator/url"] = "https://github.com/riotkit-org/git-clone-operator"

	_, _, err := resolveSecretForPod(context.TODO(), client, &pod)

	assert.Equal(t, "The secret 'my-secret-name' does not contain key 'username', while the annotation on Pod specifies that key", err.Error())
}

func TestResolvingWithUsernameIsNotMandatory(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "my-secret-name", Namespace: "default"},
		Data: map[string][]byte{
			// "username": []byte("hello"),  // this one is not mandatory, when `git-clone-operator/secretUsernameKey` annotation was not defined
			"password": []byte("riotkit"),
		},
		Type: "opaque",
	})

	pod := corev1.Pod{}
	pod.Namespace = "default"
	pod.Annotations = map[string]string{
		"git-clone-operator/secretName": "my-secret-name",
	}
	pod.Annotations["git-clone-operator/secretTokenKey"] = "password"
	pod.Annotations["git-clone-operator/revision"] = "HEAD"
	pod.Annotations["git-clone-operator/group"] = "161"
	pod.Annotations["git-clone-operator/owner"] = "161"
	pod.Annotations["git-clone-operator/path"] = "/var/www/riotkit"
	pod.Annotations["git-clone-operator/url"] = "https://github.com/riotkit-org/git-clone-operator"

	_, _, err := resolveSecretForPod(context.TODO(), client, &pod)

	assert.Nil(t, err)
}
