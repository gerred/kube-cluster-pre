// Copyright 2015 The kube-cluster Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
