/*
Copyright 2019 Platform9 Systems, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"flag"
)

const (
	defaultLogNs          = "logging"
	defaultFluentSvcAct   = "fluent"
	defaultFluentbitImage = "fluent/fluent-bit:1.0.6"
	defaultFluentdImage   = "platform9/fluentd:v1.0"
	defaultCfgDir         = "etc/conf"
	defaultFwdPort        = 62073
	defaultReloadPort     = 45550
	defaultReloadHost     = "fluentd.logging.svc.cluster.local"
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
	// ReloadPort is fluentd port used to reload fluentd config
	ReloadPort = flag.Int("reload-port", defaultReloadPort, "Fluentd config reload port")
	// ReloadHost refers to fluentd reload webhook
	ReloadHost = flag.String("reload-host", defaultReloadHost, "Fluentd reload host")
)
