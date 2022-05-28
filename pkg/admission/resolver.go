package admission

import (
	goCtx "context"
	"github.com/pkg/errors"
	"github.com/riotkit-org/git-clone-operator/pkg/context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// resolveSecretForPod Finds a `kind: Secret` using information from ResolvePod's annotations and extracts secrets from that secret by specified keys
func resolveSecretForPod(ctx goCtx.Context, client kubernetes.Interface, pod *corev1.Pod) (string, string, error) {
	// checking required annotations
	if val, exists := pod.Labels[context.AnnotationSecretName]; !exists || val == "" {
		logrus.Infof("No annotation '%s' defined for Pod '%s/%s'", context.AnnotationSecretName, pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
		return "", "", nil
	}
	if val, exists := pod.Labels[context.AnnotationSecretTokenKey]; !exists || val == "" {
		logrus.Infof("No annotation '%s' defined for Pod '%s/%s'", context.AnnotationSecretTokenKey, pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
		return "", "", nil
	}
	if val, exists := pod.Labels[context.AnnotationSecretUserKey]; !exists || val == "" {
		logrus.Infof("No annotation '%s' defined for Pod '%s/%s'", context.AnnotationSecretUserKey, pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
		return "", "", nil
	}

	// fetching `kind: Secret` from API
	secretName := pod.Labels[context.AnnotationSecretName]
	secret, err := client.CoreV1().Secrets(pod.ObjectMeta.Namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return "", "", errors.Wrapf(err, "Cannot fetch secret for ResolvePod annotated with '%s=%s'", context.AnnotationSecretName, pod.Labels[context.AnnotationSecretName])
	}

	// extracting data from `kind: Secret`
	token, tokenDefined := secret.Data[context.AnnotationSecretTokenKey]
	if !tokenDefined {
		return "", "", errors.Errorf("The secret '%s' does not contain key '%s'", secretName, secret.Data[context.AnnotationSecretTokenKey])
	}
	username, usernameDefined := secret.Data[context.AnnotationSecretTokenKey]
	if !usernameDefined {
		return "", "", errors.Errorf("The secret '%s' does not contain key '%s'", secretName, secret.Data[context.AnnotationSecretUserKey])
	}

	return string(username), string(token), nil
}
