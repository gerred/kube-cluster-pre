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

package virtualbox

import (
	"fmt"
	"os"
	"strings"
)

//https://github.com/kubernetes/kubernetes/blob/master/cluster/vagrant/util.sh#L281

//gen-kube-basicauth
//get-tokens
//create-provision-scripts

func (v *Virtualbox) createProvisionScripts() error {
	// ensure-temp-dir

	masterStart, err := os.Create(fmt.Sprintf("%v/master-start.sh", v.tempDir))
	if err != nil {
		return err
	}
	defer masterStart.Close()

	fmt.Fprintf(masterStart, "#! /bin/bash\n")
	// TODO(carlos):shared directory with hosting machine, might need some assets scp into machine
	fmt.Fprintf(masterStart, "KUBE_ROOT=/vagrant\n")

	fmt.Fprintf(masterStart, "INSTANCE_PREFIX='%v'\n", v.instancePrefix)
	fmt.Fprintf(masterStart, "MASTER_NAME='%v-master'\n", v.instancePrefix)
	fmt.Fprintf(masterStart, "MASTER_IP='%v'\n", v.masterIP)
	fmt.Fprintf(masterStart, "MINION_NAMES=(%v)\n", strings.Join(v.minionNames, " "))
	fmt.Fprintf(masterStart, "MINION_IPS=(%v)\n", strings.Join(v.minionIPs, " "))
	fmt.Fprintf(masterStart, "NODE_IP='%v'\n", v.masterIP)
	fmt.Fprintf(masterStart, "CONTAINER_SUBNET='%v'\n", v.containerSubnet)
	fmt.Fprintf(masterStart, "CONTAINER_NETMASK='%v'\n", v.masterContainerNetmask)
	fmt.Fprintf(masterStart, "MASTER_CONTAINER_SUBNET='%v'\n", v.masterContainerSubnet)
	fmt.Fprintf(masterStart, "CONTAINER_ADDR='%v'\n", v.masterContainerAddr)
	fmt.Fprintf(masterStart, "MINION_CONTAINER_NETMASKS='%v'\n", strings.Join(v.minionContainerNetmasks, " "))
	fmt.Fprintf(masterStart, "MINION_CONTAINER_SUBNETS=(%v)\n", strings.Join(v.minionContainerSubnets, " "))
	fmt.Fprintf(masterStart, "SERVICE_CLUSTER_IP_RANGE='%v'\n", v.serviceClusterIpRange)
	fmt.Fprintf(masterStart, "MASTER_USER='%v'\n", v.masterUser)
	fmt.Fprintf(masterStart, "MASTER_PASSWD='%v'\n", v.masterPasswd)
	fmt.Fprintf(masterStart, "KUBE_USER='%v'\n", v.kubeUser)
	fmt.Fprintf(masterStart, "KUBE_PASSWORD='%v'\n", v.kubePassword)
	fmt.Fprintf(masterStart, "ENABLE_CLUSTER_MONITORING='%v'\n", v.enableClusterMonitoring)
	fmt.Fprintf(masterStart, "ENABLE_NODE_LOGGING='%v'\n", v.enableNodeLogging)
	fmt.Fprintf(masterStart, "ENABLE_CLUSTER_UI='%v'\n", v.enableClusterUI)
	fmt.Fprintf(masterStart, "LOGGING_DESTINATION='%v'\n", v.loggingDestination)
	fmt.Fprintf(masterStart, "ENABLE_CLUSTER_DNS='%v'\n", v.enableClusterDns)
	fmt.Fprintf(masterStart, "DNS_SERVER_IP='%v'\n", v.dnsServerIp)
	fmt.Fprintf(masterStart, "DNS_DOMAIN='%v'\n", v.dnsDomain)
	fmt.Fprintf(masterStart, "DNS_REPLICAS='%v'\n", v.dnsReplicas)
	fmt.Fprintf(masterStart, "RUNTIME_CONFIG='%v'\n", v.runtimeConfig)
	fmt.Fprintf(masterStart, "ADMISSION_CONTROL='%v'\n", v.admissionControl)
	fmt.Fprintf(masterStart, "DOCKER_OPTS='%v'\n", v.extraDockerOpts)
	// fmt.Fprintf(masterStart,"VAGRANT_DEFAULT_PROVIDER='${VAGRANT_DEFAULT_PROVIDER:-}'\n")
	fmt.Fprintf(masterStart, "KUBELET_TOKEN='%v'\n", v.kubeletToken)
	fmt.Fprintf(masterStart, "KUBE_PROXY_TOKEN='%v'\n", v.kubeProxyToken)
	fmt.Fprintf(masterStart, "MASTER_EXTRA_SANS='%v'\n", v.masterExtraSans)
	fmt.Fprintf(masterStart, "ENABLE_CPU_CFS_QUOTA='%v'\n", v.enableCpuCfsQuota)

	// TODO(carlos): copy non commented lines of provision-network-master.sh and provision-master.sh
	// awk '!/^#/' "${KUBE_ROOT}/cluster/vagrant/provision-network-master.sh"
	// awk '!/^#/' "${KUBE_ROOT}/cluster/vagrant/provision-master.sh"

	for k, m := range v.minionNames {
		minionStart, err := os.Create(fmt.Sprintf("%v/minion-start-%v.sh", v.tempDir, k))
		if err != nil {
			return err
		}
		fmt.Fprintf(minionStart, "#! /bin/bash\n")
		fmt.Fprintf(minionStart, "MASTER_NAME='%v'\n", v.masterName)
		fmt.Fprintf(minionStart, "MASTER_IP='%v'\n", v.masterIP)
		fmt.Fprintf(minionStart, "MINION_NAMES=(%v)\n", strings.Join(v.minionNames, " "))
		fmt.Fprintf(minionStart, "MINION_NAME=(%v)\n", m)
		fmt.Fprintf(minionStart, "MINION_IPS=(%v)\n", strings.Join(v.minionIPs, " "))
		fmt.Fprintf(minionStart, "MINION_IP='%v'\n", v.minionIPs[k])
		fmt.Fprintf(minionStart, "MINION_ID='%v'\n", k)
		fmt.Fprintf(minionStart, "NODE_IP='%v'\n", v.minionIPs[k])
		fmt.Fprintf(minionStart, "MASTER_CONTAINER_SUBNET='%v'\n", v.masterContainerSubnet)
		fmt.Fprintf(minionStart, "CONTAINER_ADDR='%v'\n", v.minionContainerAddrs[k])
		fmt.Fprintf(minionStart, "CONTAINER_NETMASK='%v'\n", v.minionContainerNetmasks[k])
		fmt.Fprintf(minionStart, "MINION_CONTAINER_SUBNETS=(%v)\n", strings.Join(v.minionContainerSubnets, " "))
		fmt.Fprintf(minionStart, "CONTAINER_SUBNET='%v'\n", v.containerSubnet)
		fmt.Fprintf(minionStart, "DOCKER_OPTS='%v'\n", v.extraDockerOpts)
		fmt.Fprintf(minionStart, "KUBELET_TOKEN='%v'\n", v.kubeletToken)
		fmt.Fprintf(minionStart, "KUBE_PROXY_TOKEN='%v'\n", v.kubeProxyToken)
		fmt.Fprintf(minionStart, "MASTER_EXTRA_SANS='%v'\n", v.masterExtraSans)

		// TODO(carlos): copy non commented lines of provision-network-master.sh and provision-minion.sh
		// awk '!/^#/' "${KUBE_ROOT}/cluster/vagrant/provision-network-minion.sh"
		// awk '!/^#/' "${KUBE_ROOT}/cluster/vagrant/provision-minion.sh"
		minionStart.Close()
	}

	return nil
}
