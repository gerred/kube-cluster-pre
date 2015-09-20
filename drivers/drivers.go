package drivers

// The Driver interface is the set of functions required to start, manage, and verify a Kubernetes cluster. Inspired by the steps used in kube-up.sh, with extensions for management and scaling of clusters.
type Driver interface {
	GenerateCerts()
	GetTokens()
	ProvisionMaster()
	ConfigureMaster()
	ProvisionNode()
	ConfigureNode()
}
