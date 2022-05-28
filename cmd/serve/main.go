package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/riotkit-org/git-clone-operator/pkg/admission"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
)

type Command struct {
	LogLevel string
	TLS      bool
	LogJSON  bool

	DefaultImage       string
	DefaultGitUsername string
	DefaultGitToken    string

	client *kubernetes.Clientset
}

func (c *Command) Run() error {
	c.setLogger()
	c.client = initClient()

	// handle our core application
	http.HandleFunc("/mutate-pods", c.ServeMutatePods)
	http.HandleFunc("/health", c.ServeHealth)

	// start the server
	// listens to clear text http on port 8080 unless TLS env var is set to "true"
	if c.TLS {
		cert := "/etc/admission-webhook/tls/tls.crt"
		key := "/etc/admission-webhook/tls/tls.key"
		logrus.Print("Listening on port 4443...")
		return http.ListenAndServeTLS(":4443", cert, key, nil)
	} else {
		logrus.Print("Listening on port 8080...")
		return http.ListenAndServe(":8080", nil)
	}
}

// ServeHealth returns 200 when things are good
func (c *Command) ServeHealth(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("uri", r.RequestURI).Debug("healthy")
	fmt.Fprint(w, "OK")
}

// ServeMutatePods returns an admission review with pod mutations as a json patch
// in the review response
func (c *Command) ServeMutatePods(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("uri", r.RequestURI)
	logger.Debug("received mutation request")

	in, err := c.parseRequest(*r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	adm := admission.MutationRequest{
		Logger:  logger,
		Request: in.Request,

		DefaultImage:       c.DefaultImage,
		DefaultGitUsername: c.DefaultGitUsername,
		DefaultGitToken:    c.DefaultGitToken,

		Client: c.client,
	}

	out, err := adm.ProcessAdmissionRequest()
	if err != nil {
		e := fmt.Sprintf("could not generate admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jout, err := json.Marshal(out)
	if err != nil {
		e := fmt.Sprintf("could not parse admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	logger.Debug("sending response")
	logger.Debugf("%s", jout)
	fmt.Fprintf(w, "%s", jout)
}

// setLogger sets the logger using env vars, it defaults to text logs on
// debug level unless otherwise specified
func (c *Command) setLogger() {
	lvl, parseErr := logrus.ParseLevel(c.LogLevel)
	if parseErr != nil {
		logrus.Fatalf("Cannot parse log level: %v", parseErr)
	}
	logrus.SetLevel(lvl)

	if c.LogJSON {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

// parseRequest extracts an AdmissionReview from an http.Request if possible
func (c *Command) parseRequest(r http.Request) (*admissionv1.AdmissionReview, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	if len(body) == 0 {
		return nil, fmt.Errorf("admission request body is empty")
	}

	var a admissionv1.AdmissionReview

	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	return &a, nil
}

func initClient() *kubernetes.Clientset {
	kubeConfig := os.Getenv("HOME") + "/.kube/config"
	if os.Getenv("KUBECONFIG") != "" {
		kubeConfig = os.Getenv("KUBECONFIG")
	}
	if _, err := os.Stat(kubeConfig); errors.Is(err, os.ErrNotExist) {
		kubeConfig = ""
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientSet
}
