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

package kubectlfwd

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

const clusterObjectName = "env"

var acceptableClusterCalls = [...]struct {
	cmd string
	obj string
}{
	{"create", clusterObjectName},
	{"describe", clusterObjectName},
	{"env", clusterObjectName},
	{"get", clusterObjectName},
}

type fwd struct {
	args          []string
	kubectlBinary string
	stdout        *os.File
	stderr        *os.File

	returnCode int
}

func New(args []string, kubectl string, stdout, stderr *os.File) *fwd {
	return &fwd{
		args:          args,
		kubectlBinary: kubectl,
		stdout:        stdout,
		stderr:        stderr,
	}
}

func (f *fwd) ForwardCall() (bool, error) {
	if f.isClusterCall() {
		return false, nil
	}

	cmd := exec.Command(f.kubectlBinary, f.args[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return true, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return true, err
	}
	if err := cmd.Start(); err != nil {
		return true, err
	}

	io.Copy(f.stderr, stderr)
	io.Copy(f.stdout, stdout)
	if err := cmd.Wait(); err != nil {
		return true, err
	}

	return true, nil
}

func (f *fwd) isClusterCall() bool {
	if len(f.args) < 3 || '-' == f.args[2][0] {
		return true
	}

	cmd := strings.ToLower(f.args[1])
	obj := strings.ToLower(f.args[2])
	for _, call := range acceptableClusterCalls {
		if cmd == call.cmd && obj == call.obj {
			return true
		}
	}

	return false
}
