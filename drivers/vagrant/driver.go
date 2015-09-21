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
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Based on
// https://github.com/kubernetes/kubernetes/blob/master/cluster/vagrant/util.sh

var providers = [...]struct {
	executable, name, plugin string
}{
	{"", "vmware_fusion", "vagrant-vmware-fusion"},
	{"", "vmware_workstation", "vagrant-vmware-workstation"},
	{"prlctl", "parallels", "vagrant-parallels"},
	{"VBoxManage", "virtualbox", ""},
	{"virsh", "libvirt", "vagrant-libvirt"},
}

type Vagrant struct {
	bin              string
	detectedProvider string
}

func New() (*Vagrant, error) {
	vagrantBin, err := exec.LookPath("vagrant")
	if err != nil {
		return nil, fmt.Errorf("could not find vagrant: %v", err)
	}

	v := &Vagrant{
		bin: vagrantBin,
	}

	if err := v.detectProvider(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Vagrant) extractPlugins() (map[string]struct{}, error) {
	plugins := make(map[string]struct{})

	out, err := exec.Command(v.bin, "plugin", "list").Output()
	if err != nil {
		return plugins, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		name := strings.Fields(line)[0]
		plugins[name] = struct{}{}
	}

	return plugins, nil
}

func (v *Vagrant) detectProvider() error {
	plugins, err := v.extractPlugins()
	if err != nil {
		return fmt.Errorf("error detecting vagrant providers: %v", err)
	}

	// TODO(ccf): read $VAGRANT_DEFAULT_PROVIDER (env var)

	providerFound := ""
	for _, p := range providers {
		_, err := exec.LookPath(p.executable)
		_, ok := plugins[p.plugin]
		if (p.executable != "" && err == nil) || (p.plugin != "" && ok) {
			providerFound = p.name
			break
		}
	}

	if providerFound == "" {
		return errors.New("no vagrant provider found")
	}

	v.detectedProvider = providerFound
	return nil
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
