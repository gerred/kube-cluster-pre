package drivers

// The Driver interface is the set of functions required to start and verify a Kubernetes cluster. Inspired by the steps used in kube-up.sh
type Driver interface {
	GenerateCerts()
	GetTokens()
	ProvisionMaster()
	ConfigureMaster()
	ProvisionNode()
	ConfigureNode()
}
