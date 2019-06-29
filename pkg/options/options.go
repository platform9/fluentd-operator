package options

import (
	"flag"
)

const (
	defaultLogNs        = "logging"
	defaultFluentSvcAct = "fluent-bit"
)

var LogNs string
var SvcAcct string

// AddFlags adds command line flags for the operator
func AddFlags() {
	flag.StringVar(&LogNs, "log-ns", defaultLogNs, "Namespace for running fluent-bit and fluentd")
	flag.StringVar(&SvcAcct, "svc-acct", defaultFluentSvcAct, "Service account to use for fluentd and fluentbit")
}
