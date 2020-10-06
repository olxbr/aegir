package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"net/http"

	notifications "github.com/grupozap/aegir/internal/pkg/notifications/slack"
	"github.com/grupozap/aegir/internal/pkg/rules"
	"github.com/grupozap/aegir/internal/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	jsonContentType = `application/json`
)

var rulesFile string
var slackToken string
var listenPort string
var tlsCertPath string
var tlsKeyPath string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs Aegir's admission controller.",
	Run:   serve,
}

func init() {
	RootCmd.AddCommand(serverCmd)
	cobra.OnInitialize(initConfig)
	serverCmd.PersistentFlags().StringVar(&tlsCertPath, "tls-cert-file", "", "Path to TLS certificate file")
	serverCmd.PersistentFlags().StringVar(&tlsKeyPath, "tls-key-file", "", "Path to TLS key file")
	serverCmd.PersistentFlags().StringVar(&rulesFile, "rules-file", "", "File that contains the rules that will be applied for the Kubernetes resources.")
	serverCmd.PersistentFlags().StringVar(&slackToken, "slack-token", "", "Slack API Token to enable Aegir notifications")
	serverCmd.PersistentFlags().StringVar(&listenPort, "port", "8443", "TCP port that connections will be listen.")
}

func initConfig() {
	if rulesFile == "" {
		log.Fatalf("You must provide a rules file valid path. Eg: %s %s --rules-file=/path/to/file/rules.yaml\n", RootCmd.Name(), serverCmd.Name())
	}
}

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	skippedNamespaces, _  = utils.GetEnvAsSlice("SKIP_NAMESPACES", ",")
)

type validationFunc func(*v1beta1.AdmissionRequest) []*utils.Violation

func validateRules(req *v1beta1.AdmissionRequest) []*utils.Violation {
	raw := req.Object.Raw
	rsc := rules.Resource{}
	json.Unmarshal(raw, &rsc)
	var violationsSlice []*utils.Violation
	for _, rule := range rules.GetRules(req.Namespace, req.Kind.Kind) {
		//Skip rule if namespace is inside SKIP_NAMESPACES environment variable
		if rule.Namespace == "*" && utils.Include(skippedNamespaces, req.Namespace) {
			continue
		}
		for _, ruledef := range rule.RulesDefinitions {
			violations := ruledef.GetViolations(string(raw))
			for _, violated := range violations {
				violated.SlackChannel = rule.SlackNotificationChannel
				violated.RuleName = rule.Name
				violationsSlice = append(violationsSlice, violated)
			}
		}
	}
	return violationsSlice
}

func printValidationErrors(v []*utils.Violation) string {
	sb := strings.Builder{}
	for _, violation := range v {
		m := fmt.Sprintf("\trule name: '%s', field: '%s', description: '%s', message: %s\n", violation.RuleName, violation.JSONPath, violation.Description, violation.Message)
		sb.WriteString(m)
	}
	defer sb.Reset()
	return strings.TrimRight(sb.String(), "\n")
}

func handleAdmissionRequest(w http.ResponseWriter, r *http.Request, v validationFunc) ([]byte, error) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("only POST methods are allowed")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("could not read request body")
	}

	if contentType := r.Header.Get("Content-Type"); contentType != jsonContentType {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("Unsupported content type %s, only %s is supported", contentType, jsonContentType)
	}

	var admissionReviewReq v1beta1.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("could not deserialize the request into an admission review: %q", err)
	} else if admissionReviewReq.Request == nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New("malformed admission request: request is nil")
	}

	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID: admissionReviewReq.Request.UID,
		},
	}

	violatedRules := validateRules(admissionReviewReq.Request)
	if len(violatedRules) > 0 {
		admissionReviewResponse.Response = &v1beta1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: fmt.Sprintf("We found violations in your request. The following rules were violated: \n %s", printValidationErrors(violatedRules)),
				Code:    http.StatusForbidden,
			},
		}
		var msg notifications.NotificationMessage
		for _, violation := range violatedRules {
			msg.Message = fmt.Sprintf("Rule name: *%s*\n Rule Description: *%s*\n", violation.RuleName, violation.Description)
			msg.ResourceType = admissionReviewReq.Request.Kind.Kind
			msg.ResourceNamespace = admissionReviewReq.Request.Namespace
			go notifications.NotifyViolation(msg, slackToken, violation.SlackChannel, "#FD0D0D")
		}
	} else if len(violatedRules) == 0 {
		fmt.Printf("There was no violations!")
		admissionReviewResponse.Response.Allowed = true
	}
	response, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %q", err)
	}
	return response, nil
}

func serveAdmitFunc(w http.ResponseWriter, r *http.Request, v validationFunc) {
	log.Print("Handling webhook request ...")

	var writeErr error
	if bytes, err := handleAdmissionRequest(w, r, v); err != nil {
		log.Printf("Error handling webhook request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr = w.Write([]byte(err.Error()))
	} else {
		log.Print("Webhook request handled successfully")
		_, writeErr = w.Write(bytes)
	}

	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}

func admitFuncHandler(v validationFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveAdmitFunc(w, r, v)
	})
}

func serve(cmd *cobra.Command, args []string) {
	tlsCert, err := filepath.Abs(tlsCertPath)
	if err != nil {
		panic(err)
	}
	tlsKey, err := filepath.Abs(tlsKeyPath)
	if err != nil {
		panic(err)
	}
	rl := rules.RulesLoader(rulesFile)
	rules.BuildRuleStore(&rl)
	mux := http.NewServeMux()

	// Dummy endpoint for livenessProbes
	up := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "UP\n")
	}
	mux.HandleFunc("/healthcheck", up)
	mux.Handle("/admission", admitFuncHandler(validateRules))
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    fmt.Sprintf(":%s", listenPort),
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(tlsCert, tlsKey))
}
