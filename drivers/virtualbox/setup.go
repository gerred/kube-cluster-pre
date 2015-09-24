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
	"compress/gzip"
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

const (
	VagrantBox     = "http://opscode-vm-bento.s3.amazonaws.com/vagrant/virtualbox/opscode_fedora-21_chef-provisionerless.box"
	VagrantBoxSHA1 = "f520b4ca37ce1721fb60ff20636eeaf12bc633ca"
)

var ErrDeployedEnvironment error = errors.New("environment already deployed.")
var ErrNonSupportedArchitecture error = errors.New("non supported architecture. must be either i386 or amd64.")

func (v *Virtualbox) Setup() error {
	log.Println("launching VM")
	defer log.Println("done")
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

func (v *Virtualbox) isISOdownloaded() bool {
	fn := path.Base(VagrantBox)
	if _, err := os.Stat(fn); err == nil {
		h := sha1.New()
		r, rerr := os.Open(fn)
		if rerr != nil {
			return false
		}
		if _, err := io.Copy(h, r); err != nil {
			return false
		}
		if err := r.Close(); err != nil {
			return false
		}

		isoSHA1 := h.Sum(nil)
		if fmt.Sprintf("%x", isoSHA1) == VagrantBoxSHA1 {
			return true
		}

		v.cleanUp()
	}
	return false
}

func (v *Virtualbox) downloadISO() error {
	if v.isISOdownloaded() {
		return nil
	}

	fn := path.Base(VagrantBox)
	output, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer func(c io.Closer) {
		if err := c.Close(); err != nil {
			log.Panic(err)
		}
	}(output)

	response, err := http.Get(VagrantBox)
	if err != nil {
		return err
	}
	defer func(c io.Closer) {
		if err := c.Close(); err != nil {
			log.Panic(err)
		}
	}(response.Body)

	if _, err := io.Copy(output, response.Body); err != nil {
		return err
	}

	return nil
}

func (v *Virtualbox) openBoxTarFile() (*tar.Reader, error) {
	boxFileReader, err := os.Open(path.Base(VagrantBox))
	if err != nil {
		return nil, err
	}

	gzReader, err := gzip.NewReader(boxFileReader)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(gzReader)
	return tarReader, nil
}

func untarFile(hdr *tar.Header, tarReader *tar.Reader) error {
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
		if _, err = io.Copy(writer, tarReader); err != nil {
			return err
		}
		if err = os.Chmod(fn, os.FileMode(hdr.Mode)); err != nil {
			return err
		}
		if err = writer.Close(); err != nil {
			return err
		}

	default:
		log.Panicf("unable to untar type: %c in file %s", hdr.Typeflag, fn)
	}
	return nil
}
func (v *Virtualbox) untarBox() error {
	tarReader, err := v.openBoxTarFile()
	if err != nil {
		return err
	}
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err = untarFile(hdr, tarReader); err != nil {
			return err
		}
	}

	return nil
}

func (v *Virtualbox) importVM() error {
	steps := [...][]string{
		{"import", "box.ovf", "--vsys", "0", "--vmname", v.envName},
		{"modifyvm", v.envName, "--natpf1", "guestssh,tcp,,2222,,22"},

		// From k8s Vagrantfile
		// Use faster paravirtualized networking
		{"modifyvm", v.envName, "--nictype1", "virtio"},
		{"modifyvm", v.envName, "--nictype2", "virtio"},
	}

	for _, step := range steps {
		out, err := exec.Command(v.mgmtbin, step...).CombinedOutput()
		if err != nil {
			log.Printf("%s\n", out)
			return err
		}
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
		"fedora-21-x86_64-disk1.vmdk",
		"box.ovf",
		"opscode_fedora-21_chef-provisionerless.box",
		"metadata.json",
	}

	for _, fn := range files {
		if err := os.Remove(fn); err != nil {
			log.Println(err)
		}
	}
}
