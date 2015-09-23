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
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

const (
	VBoxManageBin   = "VBoxManage"
	VBoxHeadlessBin = "VBoxHeadless"
)

type Virtualbox struct {
	mgmtbin     string
	headlessbin string
	envName     string
}

var ErrParsingVirtualBoxVersion error = errors.New("error trying to detect VirtualBox version.")
var ErrMinVirtualBoxVersion error = errors.New("upgrade Virtualbox to at least v5.")

func New(envName string) (*Virtualbox, error) {
	dotExe := ""
	if runtime.GOOS == "windows" {
		dotExe = ".exe"
	}
	mgmtbin, err := exec.LookPath(VBoxManageBin + dotExe)
	if err != nil {
		return nil, fmt.Errorf("could not find %s: %v", VBoxManageBin, err)
	}
	headlessbin, err := exec.LookPath(VBoxHeadlessBin + dotExe)
	if err != nil {
		return nil, fmt.Errorf("could not find %s: %v", VBoxHeadlessBin, err)
	}

	v := &Virtualbox{
		mgmtbin:     mgmtbin,
		headlessbin: headlessbin,
		envName:     envName,
	}

	if err := v.isMinimumVirtualBoxVersion(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Virtualbox) isMinimumVirtualBoxVersion() error {
	// TODO(carlos): consider https://github.com/mcuadros/go-version
	versionOut, err := exec.Command(v.mgmtbin, "--version").Output()
	if err != nil {
		return err
	}
	version, err := strconv.Atoi(string(versionOut[0]))
	if err != nil {
		return ErrParsingVirtualBoxVersion
	}
	if version < 5 {
		return ErrMinVirtualBoxVersion
	}
	return nil
}

func (v *Virtualbox) GenerateCerts() {
}

func (v *Virtualbox) GetTokens() {
}

func (v *Virtualbox) ProvisionMaster() {
}

func (v *Virtualbox) ConfigureMaster() {
}

func (v *Virtualbox) ProvisionNode() {
}

func (v *Virtualbox) ConfigureNode() {
}
