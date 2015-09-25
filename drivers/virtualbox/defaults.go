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

const DefaultNumMinions = "1"
const DefaultMasterIp = "10.245.1.2"
const DefaultKubeMasterIp = "10.245.1.2"
const DefaultInstancePrefix = "kubernetes"

// const DefaultMasterName="${INSTANCE_PREFIX}-master"
const DefaultRegisterMasterKubelet = false
const DefaultMinionIpBase = "10.245.1."
const DefaultMinionContainerSubnetBase = "10.246"
const DefaultMasterContainerNetmask = "255.255.255.0"
const DefaultMasterContainerAddr = "10.246.0.1"
const DefaultMasterContainerSubnet = "10.246.0.1/24"
const DefaultContainerSubnet = "10.246.0.0/16"

// for ((i=0; i < NUM_MINIONS; i++)) do
//   MINION_IPS[$i]="${MINION_IP_BASE}$((i+3))"
//   MINION_NAMES[$i]="${INSTANCE_PREFIX}-minion-$((i+1))"
//   MINION_CONTAINER_SUBNETS[$i]="${MINION_CONTAINER_SUBNET_BASE}.$((i+1)).1/24"
//   MINION_CONTAINER_ADDRS[$i]="${MINION_CONTAINER_SUBNET_BASE}.$((i+1)).1"
//   MINION_CONTAINER_NETMASKS[$i]="255.255.255.0"
//   VAGRANT_MINION_NAMES[$i]="minion-$((i+1))"
// done

const DefaultMasterUser = "vagrant"
const DefaultMasterPasswd = "vagrant"
const DefaultAdmissionControl = "NamespaceLifecycle,NamespaceExists,LimitRanger,SecurityContextDeny,ServiceAccount,ResourceQuota"
const DefaultEnableNodeLogging = "false"
const DefaultLoggingDestination = "elasticsearch"
const DefaultEnableClusterLogging = "false"
const DefaultElasticsearchLoggingReplicas = "1"
const DefaultEnableClusterMonitoring = "influxdb"
const DefaultExtraDockerOpts = "-b=cbr0 --insecure-registry 10.0.0.0/8"
const DefaultEnableClusterDns = "true"
const DefaultDnsServerIp = "10.247.0.10"
const DefaultDnsDomain = "cluster.local"
const DefaultDnsReplicas = "1"
const DefaultRuntimeConfig = "api/v1"

// octets=($(echo "$SERVICE_CLUSTER_IP_RANGE" | sed -e 's|/.*||' -e 's/\./ /g'))
// ((octets[3]+=1))
// service_ip=$(echo "${octets[*]}" | sed 's/ /./g')
// MASTER_EXTRA_SANS="IP:${service_ip},DNS:kubernetes,DNS:kubernetes.default,DNS:kubernetes.default.svc,DNS:kubernetes.default.svc.${DNS_DOMAIN},DNS:${MASTER_NAME}"
