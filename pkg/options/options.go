package options

import (
	"flag"
)

const (
	defaultLogNs          = "logging"
	defaultFluentSvcAct   = "fluent-bit"
	defaultFluentbitImage = "fluent/fluent-bit:1.0.6"
	defaultCfgDir         = "etc"
)

var (
	LogNs          = flag.String("log-ns", defaultLogNs, "Namespace for running fluent-bit and fluentd")
	SvcAcct        = flag.String("svc-acct", defaultFluentSvcAct, "Service account to use for fluentd and fluentbit")
	FluentbitImage = flag.String("fluentbit-image", defaultFluentbitImage, "Fluentbit image")
	CfgDir         = flag.String("cfg-dir", defaultCfgDir, "Config directory")
)
