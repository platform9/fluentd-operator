package options

import (
	"flag"
)

const (
	defaultLogNs          = "logging"
	defaultFluentSvcAct   = "fluent"
	defaultFluentbitImage = "fluent/fluent-bit:1.0.6"
	defaultFluentdImage   = "platform9/fluentd:v1.0"
	defaultCfgDir         = "etc"
	defaultFwdPort        = 62073
	defaultReloadPort     = 45550
)

var (
	// LogNs is the namespace name for running fluent components
	LogNs = flag.String("log-ns", defaultLogNs, "Namespace for running fluent-bit and fluentd")
	// SvcAcct is the service account for running fluent components
	SvcAcct = flag.String("svc-acct", defaultFluentSvcAct, "Service account to use for fluentd and fluentbit")
	// FluentbitImage points to container image for running fluentbit
	FluentbitImage = flag.String("fluentbit-image", defaultFluentbitImage, "Fluentbit image")
	// FluentdImage points to container image for running fluentd
	FluentdImage = flag.String("fluentd-image", defaultFluentdImage, "Fluentd image")
	// CfgDir is the directory local to operator, which contains initial configuration of fluentd and fluentbit
	CfgDir = flag.String("cfg-dir", defaultCfgDir, "Config directory")
	// ForwardPort is fluentd port to which fluent-bit forwards logs
	ForwardPort = flag.Int("fwd-port", defaultFwdPort, "Forwarding port for fluentd")
	// ReloadPort is fluentd rpc port which can reload config
	ReloadPort = flag.Int("reload-port", defaultReloadPort, "RPC port for fluentd")
)
