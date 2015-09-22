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
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

// todo(carlos): detect CPU architecture and download appropriate ISO
// todo(carlos): detect OS and adjust VBoxManageBin properly (.exe)
// todo(carlos): test for minimum VBox version

const (
	VBoxManageBin    = "VBoxManage"
	DefaultLinuxISO  = "http://archive.ubuntu.com/ubuntu/dists/vivid/main/installer-amd64/current/images/netboot/mini.iso"
	DefaultLinuxSHA1 = "97282a3b066de4ee4c9409979737f3911f95ceab"
)

var ErrInvalidISO error = errors.New("Invalid Linux ISO")

type Virtualbox struct {
	bin     string
	envName string
}

func New(envName string) (*Virtualbox, error) {
	bin, err := exec.LookPath(VBoxManageBin)
	if err != nil {
		return nil, fmt.Errorf("could not find %s: %v", VBoxManageBin, err)
	}

	v := &Virtualbox{
		bin:     bin,
		envName: envName,
	}

	return v, nil
}

func (v *Virtualbox) Setup() error {
	if err := v.downloadISO(); err != nil {
		return err
	}

	// todo(carlos): detect CPU architecture and adjust ostype accordingly
	steps := [...][]string{
		{"createmedium", "disk", "--filename", v.envName + ".vdi", "--size", "32768"},
		{"createvm", "--name", v.envName, "--ostype", "Ubuntu_64", "--register"},

		{"storagectl", v.envName, "--name", "\"SATA Controller\"", "--add", "sata", "--controller", "IntelAHCI"},
		{"storageattach", v.envName, "--storagectl", "\"SATA Controller\"", "--port", "0", "--device", "0", "--type", "hdd", "--medium", v.envName + ".vdi"},

		{"storagectl", v.envName, "--name", "\"IDE Controller\"", "--add", "ide"},
		{"storageattach", v.envName, "--storagectl", "\"IDE Controller\"", "--port", "0", "--device", "0", "--type", "dvddrive", "--medium", "mini.iso"},
	}

	for _, s := range steps {
		out, err := exec.Command(v.bin, s...).CombinedOutput()
		if err != nil {
			log.Printf("%s\n", out)
			return err
		}
	}
	return nil
}

func (v *Virtualbox) downloadISO() error {
	// todo(carlos): if the ISO is available local, and it is consistent, then avoid the second download.

	fn := path.Base(DefaultLinuxISO)
	output, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer output.Close()

	response, err := http.Get(DefaultLinuxISO)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	h := sha1.New()
	mw := io.MultiWriter(output, h)

	if _, err := io.Copy(mw, response.Body); err != nil {
		return err
	}

	if fmt.Sprintf("%x", h.Sum(nil)) != DefaultLinuxSHA1 {
		return ErrInvalidISO
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
