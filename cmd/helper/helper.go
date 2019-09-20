package main

import (
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var mode string
var logLevel string

const (
	defaultMode     = "standalone"
	defaultLogLevel = "INFO"
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

	log.Print("Loading dynamic client")
	_, err := dynamic.NewForConfig(config)
	errExit("Failed to create client", err)

	checkCRDExists(config)
	log.Print("Found output CRD")

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

func checkCRDExists(config *rest.Config) {
	apixClient, err := apixv1beta1client.NewForConfig(config)
	errExit("Failed to load apiextensions client", err)

	_, err = apixClient.CustomResourceDefinitions().Get("outputs.logging.pf9.io", metav1.GetOptions{})

	if err != nil {
		errExit("Error while querying output CRD", err)
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
