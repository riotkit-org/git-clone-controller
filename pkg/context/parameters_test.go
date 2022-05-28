package context_test

import (
	"github.com/riotkit-org/git-clone-operator/pkg/admission"
	"github.com/riotkit-org/git-clone-operator/pkg/context"
	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

func TestValidReplicaSetAdmissionRequest(t *testing.T) {
	json := "{\"kind\":\"Pod\",\"apiVersion\":\"v1\",\"metadata\":{\"generateName\":\"iwa-ait-wordpress-hardened-6799b7754f-\",\"creationTimestamp\":null,\"labels\":{\"app.kubernetes.io/instance\":\"app-anarchizm-info\",\"app.kubernetes.io/name\":\"wordpress-hardened\",\"pod-template-hash\":\"6799b7754f\",\"riotkit.org/git-clone-operator\":\"true\"},\"annotations\":{\"git-clone-operator/path\":\"/var/www/riotkit/wp-content/themes/iwa-theme\",\"git-clone-operator/revision\":\"master\",\"git-clone-operator/secretKey\":\"gitToken\",\"git-clone-operator/secretName\":\"git-iwa\",\"git-clone-operator/url\":\"https://git.example.org/iwa-ait/iwa-theme.git\", \"git-clone-operator/owner\": \"161\", \"git-clone-operator/group\": \"161\"},\"ownerReferences\":[{\"apiVersion\":\"apps/v1\",\"kind\":\"ReplicaSet\",\"name\":\"iwa-ait-wordpress-hardened-6799b7754f\",\"uid\":\"e89116b7-f762-47e0-a337-55eb6ac205c3\",\"controller\":true,\"blockOwnerDeletion\":true}],\"managedFields\":[{\"manager\":\"k3s\",\"operation\":\"Update\",\"apiVersion\":\"v1\",\"time\":\"2022-05-28T16:16:15Z\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:annotations\":{\".\":{},\"f:git-clone-operator/path\":{},\"f:git-clone-operator/revision\":{},\"f:git-clone-operator/secretKey\":{},\"f:git-clone-operator/secretName\":{},\"f:git-clone-operator/url\":{}},\"f:generateName\":{},\"f:labels\":{\".\":{},\"f:app.kubernetes.io/instance\":{},\"f:app.kubernetes.io/name\":{},\"f:pod-template-hash\":{},\"f:riotkit.org/git-clone-operator\":{}},\"f:ownerReferences\":{\".\":{},\"k:{\\\"uid\\\":\\\"e89116b7-f762-47e0-a337-55eb6ac205c3\\\"}\":{}}},\"f:spec\":{\"f:automountServiceAccountToken\":{},\"f:containers\":{\"k:{\\\"name\\\":\\\"app\\\"}\":{\".\":{},\"f:env\":{\".\":{},\"k:{\\\"name\\\":\\\"HEALTH_CHECK_ALLOWED_SUBNET\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"HTTPS\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"PHP_DISPLAY_ERRORS\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"PHP_ERROR_REPORTING\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"WP_PAGE_URL\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}}},\"f:envFrom\":{},\"f:image\":{},\"f:imagePullPolicy\":{},\"f:livenessProbe\":{\".\":{},\"f:failureThreshold\":{},\"f:httpGet\":{\".\":{},\"f:path\":{},\"f:port\":{},\"f:scheme\":{}},\"f:periodSeconds\":{},\"f:successThreshold\":{},\"f:timeoutSeconds\":{}},\"f:name\":{},\"f:ports\":{\".\":{},\"k:{\\\"containerPort\\\":8080,\\\"protocol\\\":\\\"TCP\\\"}\":{\".\":{},\"f:containerPort\":{},\"f:name\":{},\"f:protocol\":{}}},\"f:readinessProbe\":{\".\":{},\"f:failureThreshold\":{},\"f:httpGet\":{\".\":{},\"f:path\":{},\"f:port\":{},\"f:scheme\":{}},\"f:periodSeconds\":{},\"f:successThreshold\":{},\"f:timeoutSeconds\":{}},\"f:resources\":{\".\":{},\"f:limits\":{\".\":{},\"f:cpu\":{},\"f:memory\":{}},\"f:requests\":{\".\":{},\"f:cpu\":{},\"f:memory\":{}}},\"f:securityContext\":{\".\":{},\"f:allowPrivilegeEscalation\":{}},\"f:startupProbe\":{\".\":{},\"f:failureThreshold\":{},\"f:httpGet\":{\".\":{},\"f:path\":{},\"f:port\":{},\"f:scheme\":{}},\"f:periodSeconds\":{},\"f:successThreshold\":{},\"f:timeoutSeconds\":{}},\"f:terminationMessagePath\":{},\"f:terminationMessagePolicy\":{},\"f:volumeMounts\":{\".\":{},\"k:{\\\"mountPath\\\":\\\"/var/www/riotkit/wp-content\\\"}\":{\".\":{},\"f:mountPath\":{},\"f:name\":{}}}},\"k:{\\\"name\\\":\\\"waf-proxy\\\"}\":{\".\":{},\"f:env\":{\".\":{},\"k:{\\\"name\\\":\\\"CADDY_PORT\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"DEBUG\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"ENABLE_CORAZA_WAF\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"ENABLE_CRS\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"ENABLE_RATE_LIMITER\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"ENABLE_RULE_WORDPRESS\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"RATE_LIMIT_EVENTS\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"RATE_LIMIT_WINDOW\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}},\"k:{\\\"name\\\":\\\"UPSTREAM_0\\\"}\":{\".\":{},\"f:name\":{},\"f:value\":{}}},\"f:image\":{},\"f:imagePullPolicy\":{},\"f:livenessProbe\":{\".\":{},\"f:failureThreshold\":{},\"f:httpGet\":{\".\":{},\"f:path\":{},\"f:port\":{},\"f:scheme\":{}},\"f:periodSeconds\":{},\"f:successThreshold\":{},\"f:timeoutSeconds\":{}},\"f:name\":{},\"f:ports\":{\".\":{},\"k:{\\\"containerPort\\\":2019,\\\"protocol\\\":\\\"TCP\\\"}\":{\".\":{},\"f:containerPort\":{},\"f:name\":{},\"f:protocol\":{}},\"k:{\\\"containerPort\\\":8081,\\\"protocol\\\":\\\"TCP\\\"}\":{\".\":{},\"f:containerPort\":{},\"f:name\":{},\"f:protocol\":{}},\"k:{\\\"containerPort\\\":8090,\\\"protocol\\\":\\\"TCP\\\"}\":{\".\":{},\"f:containerPort\":{},\"f:name\":{},\"f:protocol\":{}}},\"f:resources\":{},\"f:terminationMessagePath\":{},\"f:terminationMessagePolicy\":{},\"f:volumeMounts\":{\".\":{},\"k:{\\\"mountPath\\\":\\\"/etc/caddy/rules/custom.conf\\\"}\":{\".\":{},\"f:mountPath\":{},\"f:name\":{},\"f:subPath\":{}}}}},\"f:dnsPolicy\":{},\"f:enableServiceLinks\":{},\"f:nodeSelector\":{},\"f:restartPolicy\":{},\"f:schedulerName\":{},\"f:securityContext\":{\".\":{},\"f:fsGroup\":{},\"f:runAsGroup\":{},\"f:runAsUser\":{}},\"f:terminationGracePeriodSeconds\":{},\"f:volumes\":{\".\":{},\"k:{\\\"name\\\":\\\"waf-custom-config\\\"}\":{\".\":{},\"f:configMap\":{\".\":{},\"f:defaultMode\":{},\"f:name\":{}},\"f:name\":{}},\"k:{\\\"name\\\":\\\"wp-content\\\"}\":{\".\":{},\"f:name\":{},\"f:persistentVolumeClaim\":{\".\":{},\"f:claimName\":{}}}}}}}]},\"spec\":{\"volumes\":[{\"name\":\"wp-content\",\"persistentVolumeClaim\":{\"claimName\":\"wp-content\"}},{\"name\":\"waf-custom-config\",\"configMap\":{\"name\":\"iwa-ait-wordpress-hardened-waf-custom-config\",\"defaultMode\":420}}],\"containers\":[{\"name\":\"waf-proxy\",\"image\":\"ghcr.io/riotkit-org/waf-proxy:snapshot\",\"ports\":[{\"name\":\"http-waf\",\"containerPort\":8090,\"protocol\":\"TCP\"},{\"name\":\"waf-metrics\",\"containerPort\":2019,\"protocol\":\"TCP\"},{\"name\":\"waf-healthcheck\",\"containerPort\":8081,\"protocol\":\"TCP\"}],\"env\":[{\"name\":\"CADDY_PORT\",\"value\":\"8090\"},{\"name\":\"UPSTREAM_0\",\"value\":\"{\\\"pass_to\\\": \\\"http://127.0.0.1:8080\\\", \\\"hostname\\\": \\\"iwa-ait.org\\\"}\"},{\"name\":\"DEBUG\",\"value\":\"true\"},{\"name\":\"ENABLE_CORAZA_WAF\",\"value\":\"false\"},{\"name\":\"ENABLE_CRS\",\"value\":\"true\"},{\"name\":\"ENABLE_RATE_LIMITER\",\"value\":\"true\"},{\"name\":\"ENABLE_RULE_WORDPRESS\",\"value\":\"true\"},{\"name\":\"RATE_LIMIT_EVENTS\",\"value\":\"30\"},{\"name\":\"RATE_LIMIT_WINDOW\",\"value\":\"5s\"}],\"resources\":{},\"volumeMounts\":[{\"name\":\"waf-custom-config\",\"mountPath\":\"/etc/caddy/rules/custom.conf\",\"subPath\":\"custom.conf\"}],\"livenessProbe\":{\"httpGet\":{\"path\":\"/\",\"port\":\"waf-healthcheck\",\"scheme\":\"HTTP\"},\"timeoutSeconds\":1,\"periodSeconds\":60,\"successThreshold\":1,\"failureThreshold\":2},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\",\"imagePullPolicy\":\"Always\"},{\"name\":\"app\",\"image\":\"ghcr.io/riotkit-org/wordpress-hardened:master\",\"ports\":[{\"name\":\"http\",\"containerPort\":8080,\"protocol\":\"TCP\"}],\"envFrom\":[{\"secretRef\":{\"name\":\"wordpress-anarchizm-info\"}}],\"env\":[{\"name\":\"HEALTH_CHECK_ALLOWED_SUBNET\",\"value\":\"10.0.0.0/8\"},{\"name\":\"HTTPS\",\"value\":\"on\"},{\"name\":\"PHP_DISPLAY_ERRORS\",\"value\":\"On\"},{\"name\":\"PHP_ERROR_REPORTING\",\"value\":\"E_ALL\"},{\"name\":\"WP_PAGE_URL\",\"value\":\"https://anarchizm-info.c1.riotkit.org\"}],\"resources\":{\"limits\":{\"cpu\":\"100m\",\"memory\":\"128Mi\"},\"requests\":{\"cpu\":\"0\",\"memory\":\"50Mi\"}},\"volumeMounts\":[{\"name\":\"wp-content\",\"mountPath\":\"/var/www/riotkit/wp-content\"}],\"livenessProbe\":{\"httpGet\":{\"path\":\"/liveness.php\",\"port\":\"http\",\"scheme\":\"HTTP\"},\"timeoutSeconds\":1,\"periodSeconds\":60,\"successThreshold\":1,\"failureThreshold\":2},\"readinessProbe\":{\"httpGet\":{\"path\":\"/readiness.php\",\"port\":\"http\",\"scheme\":\"HTTP\"},\"timeoutSeconds\":1,\"periodSeconds\":60,\"successThreshold\":1,\"failureThreshold\":2},\"startupProbe\":{\"httpGet\":{\"path\":\"/liveness.php\",\"port\":\"http\",\"scheme\":\"HTTP\"},\"timeoutSeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"failureThreshold\":10},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\",\"imagePullPolicy\":\"Always\",\"securityContext\":{\"allowPrivilegeEscalation\":false}}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":5,\"dnsPolicy\":\"ClusterFirst\",\"nodeSelector\":{\"node-instance-title\":\"compute-2\",\"node-type\":\"compute\"},\"serviceAccountName\":\"default\",\"serviceAccount\":\"default\",\"automountServiceAccountToken\":false,\"securityContext\":{\"runAsUser\":65161,\"runAsGroup\":65161,\"fsGroup\":65161},\"schedulerName\":\"default-scheduler\",\"tolerations\":[{\"key\":\"node.kubernetes.io/not-ready\",\"operator\":\"Exists\",\"effect\":\"NoExecute\",\"tolerationSeconds\":300},{\"key\":\"node.kubernetes.io/unreachable\",\"operator\":\"Exists\",\"effect\":\"NoExecute\",\"tolerationSeconds\":300}],\"priority\":0,\"enableServiceLinks\":true,\"preemptionPolicy\":\"PreemptLowerPriority\"},\"status\":{}}"

	req := admission.MutationRequest{Request: &admissionv1.AdmissionRequest{Object: runtime.RawExtension{Raw: []byte(json)}}}
	req.Request.Kind.Kind = "Pod"
	req.Request.Namespace = "backups"
	pod, _ := admission.ResolvePod(req)

	parameters, err := context.NewCheckoutParametersFromPod(pod, "image", "__token__", "riotkit", "riotkit-org", "token")

	assert.Nil(t, err)
	assert.Equal(t, "image", parameters.Image)
	assert.Equal(t, "token", parameters.GitToken)          // not a default
	assert.Equal(t, "riotkit-org", parameters.GitUsername) // not a default
}
