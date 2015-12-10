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
	"github.com/gerred/kube-cluster/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gerred/kube-cluster/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gerred/kube-cluster/cluster"
)

var KubeRoot = ""

// KubeClusterCmd is the root command. Attach all other commands to this.
var KubeClusterCmd = &cobra.Command{
	Use:   "kube-cluster",
	Short: "kube-cluster provisions, scales, and manages kubernetes environments",
	Long:  "kube-cluster provisions, scales, and manages kubernetes environments",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var providerName string

var createEnvCmd = &cobra.Command{
	Use:   "create-env [name]",
	Short: "Create a Kubernetes cluster with the given name and provider options",
	Run:   CreateCluster,
}

var envVars = [...]string{
	"ADMISSION_CONTROL",
	"CONTAINER_SUBNET",
	"DNS_DOMAIN",
	"DNS_REPLICAS",
	"DNS_SERVER_IP",
	"ENABLE_CLUSTER_DNS",
	"ENABLE_CLUSTER_MONITORING",
	"ENABLE_CLUSTER_UI",
	"ENABLE_CPU_CFS_QUOTA",
	"ENABLE_NODE_LOGGING",
	"EXTRA_DOCKER_OPTS",
	"INSTANCE_PREFIX",
	"KUBE_PASSWORD",
	"KUBE_PROXY_TOKEN",
	"KUBE_USER",
	"KUBELET_TOKEN",
	"LOGGING_DESTINATION",
	"MASTER_CONTAINER_ADDR",
	"MASTER_CONTAINER_NETMASK",
	"MASTER_CONTAINER_SUBNET",
	"MASTER_EXTRA_SANS",
	"MASTER_IP",
	"MASTER_PASSWD",
	"MASTER_USER",
	"MINION_CONTAINER_NETMASKS",
	"MINION_CONTAINER_SUBNETS",
	"MINION_IPS",
	"MINION_NAMES",
	"RUNTIME_CONFIG",
	"SERVICE_CLUSTER_IP_RANGE",
	"VAGRANT_DEFAULT_PROVIDER",
}

func init() {
	createEnvCmd.Flags().StringVarP(&providerName, "provider", "p", "virtualbox", "specify which provider to use")
	cluster.GetOptions(createEnvCmd)
	for _, envVar := range envVars {
		viper.BindEnv(envVar)
	}
}
