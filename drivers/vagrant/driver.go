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

package vagrant

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Based on
// https://github.com/kubernetes/kubernetes/blob/master/cluster/vagrant/util.sh

var providers = [...]struct {
	executable, name, pluginRe string
}{
	{"", "vmware_fusion", "vagrant-vmware-fusion"},
	{"", "vmware_workstation", "vagrant-vmware-workstation"},
	{"prlctl", "parallels", "vagrant-parallels"},
	{"VBoxManage", "virtualbox", ""},
	{"virsh", "libvirt", "vagrant-libvirt"},
}

type Vagrant struct {
	bin     string
	plugins []string
}

func New() (*Vagrant, error) {
	vagrantBin, err := exec.LookPath("vagrant")
	if err != nil {
		return nil, fmt.Errorf("could not find vagrant: %v", err)
	}

	v := &Vagrant{
		bin: vagrantBin,
	}

	plugins, err := v.extractPlugins()
	if err != nil {
		return nil, err
	}

	v.plugins = plugins

	return v, nil
}

func (v *Vagrant) extractPlugins() ([]string, error) {
	var plugins []string

	out, err := exec.Command(v.bin, "plugin", "list").Output()
	if err != nil {
		return plugins, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		plugins = append(plugins, strings.Fields(line)[0])
	}

	return plugins, nil
}

func (v *Vagrant) GenerateCerts() {
}

func (v *Vagrant) GetTokens() {
}

func (v *Vagrant) ProvisionMaster() {
}

func (v *Vagrant) ConfigureMaster() {
}

func (v *Vagrant) ProvisionNode() {
}

func (v *Vagrant) ConfigureNode() {
}
