package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	logclient "github.com/platform9/fluentd-operator/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var mode string
var logLevel string
var dataNs string
var dataSrc string
var dataStore string
var timeout int

const (
	defaultMode      = "k8s"
	defaultLogLevel  = "INFO"
	defaultDataNs    = "pf9-operators"
	defaultDataSrc   = "pf9-log"
	defaultObjNs     = "logging"
	defaultDataStore = "elasticsearch"
	defaultTimeout   = 10
)

var (
	runtimeClassGVR = schema.GroupVersionResource{
		Group:    "logging.pf9.io",
		Version:  "v1alpha1",
		Resource: "outputs",
	}
)

// esParams stores elasticsearch data
type esParams struct {
	Name       string
	Namespace  string
	Deployment string
	Port       uint16
}

// Main starts it all
func Main() int {
	log.SetFormatter(&log.JSONFormatter{})

	log.Print("Loading client config")
	config := getConfig()

	log.Print("Loading client")
	apiClient, err := apixv1beta1client.NewForConfig(config)
	errExit("Failed to create client", err)

	waitForCRDs(apiClient)
	log.Print("Found output CRD")

	cs, err := kubernetes.NewForConfig(config)
	errExit("Failed to create core clientset", err)

	log.Print("Creating logging operator client")
	lc, err := logclient.NewForConfig(config)
	errExit("Failed to create logging operator client", err)

	log.Print("Configuring default backend datastore")
	configureDataStore(cs.CoreV1(), lc)

	log.Print("Creating Output CRs")
	createCrs(cs.CoreV1(), lc)

	log.Print("Done!")
	return 0
}

func getConfig() *rest.Config {
	var cfg *rest.Config
	var err error
	switch viper.GetString("mode") {
	case "standalone":
		cfg, err = getByKubeCfg()
	case "k8s":
		cfg, err = getInCluster()
	default:
		errExit("unsupported mode", os.ErrInvalid)
	}

	if err != nil {
		errExit("getting kubeconfig", err)
	}

	return cfg
}

func getInCluster() (*rest.Config, error) {
	return rest.InClusterConfig()
}

func getByKubeCfg() (*rest.Config, error) {
	defaultKubeCfg := path.Join(os.Getenv("HOME"), ".kube", "config")
	if os.Getenv("KUBECONFIG") != "" {
		defaultKubeCfg = os.Getenv("KUBECONFIG")
	}

	return clientcmd.BuildConfigFromFlags("", defaultKubeCfg)
}

func waitForCRDs(apixClient apixv1beta1client.ApiextensionsV1beta1Interface) {
	timeoutDuration := time.After(time.Duration(timeout) * time.Minute)
	tickDuration := time.Tick(10 * time.Second)

	for {
		select {
		case <-timeoutDuration:
			errExit("waiting for CRD's", fmt.Errorf("Timed out waiting for CRD's to come up"))
			return
		case <-tickDuration:
			_, err := apixClient.CustomResourceDefinitions().Get("outputs.logging.pf9.io", metav1.GetOptions{})
			if err == nil {
				return
			}
		}
	}
}

func configureDataStore(coreClient corev1.CoreV1Interface, lc logclient.LoggingV1alpha1Interface) {
	object := &loggingv1alpha1.Output{}

	// Check default datastore to use and set values
	if dataStore == "elasticsearch" {
		// Creating Output object for elastic search
		params := &esParams{
			Name:       "es-objstore",
			Namespace:  dataNs,
			Deployment: "elasticsearch",
			Port:       9200,
		}
		object = params.getESOutputObject()
	}

	_, err := lc.LoggingV1alpha1().Outputs().Create(object)
	if err != nil {
		errExit("while creating default Output object", err)
	}
}

func createCrs(coreClient corev1.CoreV1Interface, lc logclient.LoggingV1alpha1Interface) {
	// Read the secret, create struct
	sec, err := coreClient.Secrets(dataNs).Get(dataSrc, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// Assumptions: Before CRDs show up, secret containing user-data is already created.
		log.Print("Secret was not found, assuming no customizations needed")
		os.Exit(0)
	}

	errExit("while querying data secret", err)

	v, ok := sec.Data["user-data"]
	if !ok {
		errExit("Cannot find user-data in secret", os.ErrNotExist)
	}

	var outputs []loggingv1alpha1.Output

	errExit("while parsing user-data", json.Unmarshal(v, &outputs))

	for _, o := range outputs {
		_, err = lc.LoggingV1alpha1().Outputs().Create(&o)
		if errors.IsAlreadyExists(err) {
			log.Printf("Output %s already exists in %s, skipping", o.Name, o.Namespace)
			continue
		} else {
			errExit("while creating output object", err)
		}
	}
}

func errExit(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %#v", msg, err)
	}
}

func (p *esParams) getESOutputObject() *loggingv1alpha1.Output {
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", p.Deployment, p.Namespace, p.Port)
	return &loggingv1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},
		Spec: loggingv1alpha1.OutputSpec{
			Type: "elasticsearch",
			Params: []loggingv1alpha1.Param{
				{
					Name:  "url",
					Value: url,
				},
			},
		},
	}
}

func buildCmd() *cobra.Command {
	cobra.OnInitialize(initCfg)
	rootCmd := &cobra.Command{
		Use:   "helper",
		Short: "Helper creates initial log output custom resources",
		Long:  "Helper checks for user customizations and creates appropriate log output CRs as needed.",
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(Main())
		},
	}

	pf := rootCmd.PersistentFlags()
	pf.StringVar(&mode, "mode", defaultMode, "Operational mode: k8s or standalone")
	viper.BindPFlag("mode", pf.Lookup("mode"))

	pf.StringVar(&logLevel, "log-level", defaultLogLevel, "Log level: DEBUG, INFO, WARN or FATAL")
	viper.BindPFlag("log-level", pf.Lookup("log-level"))

	pf.StringVar(&dataNs, "datasource-ns", defaultDataNs, "Namespace of user-data")
	viper.BindPFlag("datasource-ns", pf.Lookup("datasource-ns"))

	pf.StringVar(&dataSrc, "datasource", defaultDataSrc, "Name of secret for user config data")
	viper.BindPFlag("datasource", pf.Lookup("datasource"))

	pf.StringVar(&dataStore, "datastore", defaultDataStore, "Name of the backend datastore")
	viper.BindPFlag("datastore", pf.Lookup("datastore"))

	pf.IntVar(&timeout, "timeout", defaultTimeout, "Wait period for the CRD's")
	viper.BindPFlag("timeout", pf.Lookup("timeout"))

	return rootCmd
}

func initCfg() {
	mode := viper.GetString("mode")
	if mode != "k8s" && mode != "standalone" {
		fmt.Fprintf(os.Stderr, "mode %s is invalid", mode)
		os.Exit(1)
	}

	lvl := viper.GetString("log-level")
	switch lvl {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	default:
		fmt.Fprintf(os.Stderr, "log level %s is invalid", lvl)
	}
}

func main() {
	cmd := buildCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}
