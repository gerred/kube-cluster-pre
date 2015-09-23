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
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
)

const (
	VagrantBox = "https://cloud-images.ubuntu.com/vagrant/vivid/current/vivid-server-cloudimg-%s-vagrant-disk1.box"
)

var ErrDeployedEnvironment error = errors.New("environment already deployed.")
var ErrNonSupportedArchitecture error = errors.New("non supported architecture. must be either i386 or amd64.")

func (v *Virtualbox) Setup() error {
	if v.isDeployedEnv() {
		return ErrDeployedEnvironment
	}
	defer v.cleanUp()

	steps := []func() error{
		v.downloadISO,
		v.untarBox,
		v.importVM,
		v.startVM,
	}

	for _, step := range steps {
		fmt.Print(".")
		if err := step(); err != nil {
			return err
		}
	}

	return nil
}

func (v *Virtualbox) isDeployedEnv() bool {
	if _, err := exec.Command(v.mgmtbin, "showvminfo", v.envName, "--machinereadable").Output(); err != nil {
		return false
	}

	return true
}

func (v *Virtualbox) boxURL() (string, error) {
	switch runtime.GOARCH {
	case "386":
		return fmt.Sprintf(VagrantBox, "i386"), nil
	case "amd64":
		return fmt.Sprintf(VagrantBox, "amd64"), nil
	}
	return "", ErrNonSupportedArchitecture
}

func (v *Virtualbox) downloadISO() error {
	url, err := v.boxURL()
	if err != nil {
		return err
	}

	output, err := os.Create(path.Base(url))
	if err != nil {
		return err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if _, err := io.Copy(output, response.Body); err != nil {
		return err
	}

	return nil
}

func (v *Virtualbox) untarBox() error {
	url, err := v.boxURL()
	if err != nil {
		return err
	}

	boxFileReader, err := os.Open(path.Base(url))
	if err != nil {
		return err
	}
	defer boxFileReader.Close()

	tarReader := tar.NewReader(boxFileReader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fn := hdr.Name
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fn, os.FileMode(hdr.Mode)); err != nil {
				return err
			}

		case tar.TypeReg:
			writer, err := os.Create(fn)
			if err != nil {
				return err
			}
			io.Copy(writer, tarReader)
			if err = os.Chmod(fn, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
			writer.Close()

		default:
			log.Panicf("unable to untar type: %c in file %s", hdr.Typeflag, fn)
		}
	}

	return nil
}

func (v *Virtualbox) importVM() error {
	out, err := exec.Command(v.mgmtbin, "import", "box.ovf", "--vsys", "0", "--vmname", v.envName).CombinedOutput()
	if err != nil {
		log.Printf("%s\n", out)
		return err
	}
	return nil
}

func (v *Virtualbox) startVM() error {
	if err := exec.Command(v.headlessbin, "-s", v.envName).Start(); err != nil {
		return err
	}
	return nil
}

func (v *Virtualbox) cleanUp() {
	files := [...]string{
		"Vagrantfile",
		"box-disk1.vmdk",
		"box.ovf",
		"vivid-server-cloudimg-amd64-vagrant-disk1.box",
	}

	for _, fn := range files {
		if err := os.Remove(fn); err != nil {
			log.Println(err)
		}
	}
}
