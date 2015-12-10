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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/gerred/kube-cluster/Godeps/_workspace/src/github.com/spf13/viper"
)

const masterSh = `#! /bin/bash
KUBE_ROOT='/vagrant'
INSTANCE_PREFIX='{{ .InstancePrefix }}'
MASTER_NAME='{{ .InstancePrefix }}-master'
MASTER_IP='{{ .MasterIp }}'
MINION_NAMES=({{ .AllMinionNames }})
MINION_IPS=({{ .AllMinionIps }})
NODE_IP='{{ .MasterIp }}'
CONTAINER_SUBNET='{{ .ContainerSubnet }}'
CONTAINER_NETMASK='{{ .MasterContainerNetmask }}'
MASTER_CONTAINER_SUBNET='{{ .MasterContainerSubnet }}'
CONTAINER_ADDR='{{ .MasterContainerAddr }}'
MINION_CONTAINER_NETMASKS='{{ .AllMinionContainerNetmasks }}'
MINION_CONTAINER_SUBNETS=({{ .AllMinionContainerSubnets }})
SERVICE_CLUSTER_IP_RANGE='{{ .ServiceClusterIpRange }}'
MASTER_USER='{{ .MasterUser }}'
MASTER_PASSWD='{{ .MasterPasswd }}'
KUBE_USER='{{ .KubeUser }}'
KUBE_PASSWORD='{{ .KubePassword }}'
ENABLE_CLUSTER_MONITORING='{{ .EnableClusterMonitoring }}'
ENABLE_NODE_LOGGING='{{ .EnableNodeLogging }}'
ENABLE_CLUSTER_UI='{{ .EnableClusterUI }}'
LOGGING_DESTINATION='{{ .LoggingDestination }}'
ENABLE_CLUSTER_DNS='{{ .EnableClusterDns }}'
DNS_SERVER_IP='{{ .DnsServerIp }}'
DNS_DOMAIN='{{ .DnsDomain }}'
DNS_REPLICAS='{{ .DnsReplicas }}'
RUNTIME_CONFIG='{{ .RuntimeConfig }}'
ADMISSION_CONTROL='{{ .AdmissionControl }}'
DOCKER_OPTS='{{ .ExtraDockerOpts }}'
VAGRANT_DEFAULT_PROVIDER='{{ .VagrantDefaultProvider }}'
KUBELET_TOKEN='{{ .KubeletToken }}'
KUBE_PROXY_TOKEN='{{ .KubeProxyToken }}'
MASTER_EXTRA_SANS='{{ .MasterExtraSans }}'
ENABLE_CPU_CFS_QUOTA='{{ .EnableCpuCfsQuota }}'
`

func (v *Virtualbox) ProvisionMaster() error {
	masterShFile, err := os.Create("master-start.sh")
	if err != nil {
		return err
	}

	tmpl, err := template.New("test").Parse(masterSh)
	if err != nil {
		log.Panic(err)
	}

	values := struct {
		AdmissionControl           string
		AllMinionContainerNetmasks string
		AllMinionContainerSubnets  string
		AllMinionIps               string
		AllMinionNames             string
		ContainerSubnet            string
		DnsDomain                  string
		DnsReplicas                string
		DnsServerIp                string
		EnableClusterDns           string
		EnableClusterMonitoring    string
		EnableClusterUI            string
		EnableCpuCfsQuota          string
		EnableNodeLogging          string
		ExtraDockerOpts            string
		InstancePrefix             string
		KubeletToken               string
		KubePassword               string
		KubeProxyToken             string
		KubeUser                   string
		LoggingDestination         string
		MasterContainerAddr        string
		MasterContainerNetmask     string
		MasterContainerSubnet      string
		MasterExtraSans            string
		MasterIp                   string
		MasterPasswd               string
		MasterUser                 string
		RuntimeConfig              string
		ServiceClusterIpRange      string
		VagrantDefaultProvider     string
	}{
		AdmissionControl:           viper.GetString("ADMISSION_CONTROL"),
		AllMinionContainerNetmasks: viper.GetString("CONTAINER_SUBNET"),
		AllMinionContainerSubnets:  viper.GetString("DNS_DOMAIN"),
		AllMinionIps:               viper.GetString("DNS_REPLICAS"),
		AllMinionNames:             viper.GetString("DNS_SERVER_IP"),
		ContainerSubnet:            viper.GetString("ENABLE_CLUSTER_DNS"),
		DnsDomain:                  viper.GetString("ENABLE_CLUSTER_MONITORING"),
		DnsReplicas:                viper.GetString("ENABLE_CLUSTER_UI"),
		DnsServerIp:                viper.GetString("ENABLE_CPU_CFS_QUOTA"),
		EnableClusterDns:           viper.GetString("ENABLE_NODE_LOGGING"),
		EnableClusterMonitoring:    viper.GetString("EXTRA_DOCKER_OPTS"),
		EnableClusterUI:            viper.GetString("INSTANCE_PREFIX"),
		EnableCpuCfsQuota:          viper.GetString("KUBE_PASSWORD"),
		EnableNodeLogging:          viper.GetString("KUBE_PROXY_TOKEN"),
		ExtraDockerOpts:            viper.GetString("KUBE_USER"),
		InstancePrefix:             viper.GetString("KUBELET_TOKEN"),
		KubeletToken:               viper.GetString("LOGGING_DESTINATION"),
		KubePassword:               viper.GetString("MASTER_CONTAINER_ADDR"),
		KubeProxyToken:             viper.GetString("MASTER_CONTAINER_NETMASK"),
		KubeUser:                   viper.GetString("MASTER_CONTAINER_SUBNET"),
		LoggingDestination:         viper.GetString("MASTER_EXTRA_SANS"),
		MasterContainerAddr:        viper.GetString("MASTER_IP"),
		MasterContainerNetmask:     viper.GetString("MASTER_PASSWD"),
		MasterContainerSubnet:      viper.GetString("MASTER_USER"),
		MasterExtraSans:            viper.GetString("MINION_CONTAINER_NETMASKS"),
		MasterIp:                   viper.GetString("MINION_CONTAINER_SUBNETS"),
		MasterPasswd:               viper.GetString("MINION_IPS"),
		MasterUser:                 viper.GetString("MINION_NAMES"),
		RuntimeConfig:              viper.GetString("RUNTIME_CONFIG"),
		ServiceClusterIpRange:      viper.GetString("SERVICE_CLUSTER_IP_RANGE"),
		VagrantDefaultProvider:     viper.GetString("VAGRANT_DEFAULT_PROVIDER"),
	}

	if values.EnableNodeLogging == "" {
		values.EnableNodeLogging = "false"
	}
	if values.EnableClusterDns == "" {
		values.EnableClusterDns = "false"
	}

	// TODO(carlos): implode (or count) for MINION_NAMES, MINION_IPS,
	// MINION_CONTAINER_NETMASKS, MINION_CONTAINER_SUBNETS

	err = tmpl.Execute(masterShFile, values)
	if err != nil {
		log.Panic(err)
	}

	appendFiles := [...]string{
		"cluster/vagrant/provision-network-master.sh",
		"cluster/vagrant/provision-master.sh",
	}

	for _, appendFile := range appendFiles {
		f, err := os.Open(fmt.Sprintf("%v/%v", v.kubeRoot, appendFile))
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(strings.TrimSpace(line), "#") {
				continue
			}
			fmt.Fprintln(masterShFile, line)
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}
