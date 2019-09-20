package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	logclient "github.com/platform9/fluentd-operator/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
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

const (
	defaultMode     = "standalone"
	defaultLogLevel = "INFO"
	defaultDataNs   = "pf9-operators"
	defaultDataSrc  = "pf9-log"
	defaultObjNs    = "logging"
)

var (
	runtimeClassGVR = schema.GroupVersionResource{
		Group:    "logging.pf9.io",
		Version:  "v1alpha1",
		Resource: "outputs",
	}
)

// Main starts it all
func Main() int {
	log.SetFormatter(&log.JSONFormatter{})

	log.Print("Loading client config")
	config := getConfig()

	log.Print("Loading client")
	apiClient, err := apixv1beta1client.NewForConfig(config)
	errExit("Failed to create client", err)

	checkCRDExists(apiClient)
	log.Print("Found output CRD")

	cs, err := kubernetes.NewForConfig(config)
	errExit("Failed to create core clientset", err)

	log.Print("Creating logging operator client")
	lc, err := logclient.NewForConfig(config)
	errExit("Failed to create logging operator client", err)

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

func checkCRDExists(apixClient apixv1beta1client.ApiextensionsV1beta1Interface) {
	_, err := apixClient.CustomResourceDefinitions().Get("outputs.logging.pf9.io", metav1.GetOptions{})

	errExit("Error while querying output CRD", err)

}

func createCrs(coreClient corev1.CoreV1Interface, lc logclient.LoggingV1alpha1Interface) {
	// Read the secret, create struct
	sec, err := coreClient.Secrets(dataNs).Get(dataSrc, metav1.GetOptions{})
	errExit("while querying data secret", err)

	v, ok := sec.Data["user-data"]
	if !ok {
		errExit("Cannot find user-data in secret", os.ErrNotExist)
	}

	var outputs []loggingv1alpha1.Output

	errExit("while parsing user-data", json.Unmarshal(v, &outputs))

	for _, o := range outputs {
		_, err = lc.LoggingV1alpha1().Outputs().Create(&o)
		errExit("while creating output object", err)
	}
}

func errExit(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %#v", msg, err)
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
