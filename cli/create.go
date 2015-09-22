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

package cli

import (
	"fmt"
	"log"

	"github.com/gerred/kube-cluster/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gerred/kube-cluster/cluster"
	"github.com/gerred/kube-cluster/drivers"
)

var providerName string

var createEnvCmd = &cobra.Command{
	Use:   "create-env [name]",
	Short: "Create a Kubernetes cluster with the given name and provider options",
	Run:   CreateCluster,
}

func init() {
	createEnvCmd.Flags().StringVarP(&providerName, "provider", "p", "virtualbox", "specify which provider to use")
}

//CreateCluster creates a new Kubernetes cluster with a provider name and options.
func CreateCluster(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Fatal("name needs to be provided")
	}

	provider, err := drivers.Factory(providerName)
	if err != nil {
		log.Panic("error loading driver %s. %v", providerName, err)
	}
	fmt.Println("Using", providerName)

	cluster, err := kubeUp(provider)
	if err != nil {
		log.Fatal("could not install kubernetes in %v", providerName)
	}

	if !cluster.IsValid() {
		log.Fatal("could not install kubernetes in %v. invalid cluster outcome.", providerName)
	}

	fmt.Println("cluster info")
	fmt.Println("\t", cluster.Info())
}

func kubeUp(provider drivers.Driver) (*cluster.Cluster, error) {
	c := new(cluster.Cluster)
	fmt.Println("kube up - start")
	defer fmt.Println("kube up - done")

	fmt.Println("\tgen kube basicauth")
	fmt.Println("\tget tokens")

	// c.Master := provider.ProvisionMaster()
	provider.ProvisionMaster()
	// c.Nodes := provider.ProvisionNode()
	provider.ProvisionNode()

	fmt.Println("\tscp kubernetes")
	fmt.Println("\tscp certificates")
	fmt.Println("\tcreate kubeconfig")

	fmt.Println("\tverify cluster")

	return c, nil
}
