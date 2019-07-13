package resources

// Resource is an interface to be implemented by all fluent-operator constructs. It is used to build configuration for fluentd
type Resource interface {
	Render() []byte
}
