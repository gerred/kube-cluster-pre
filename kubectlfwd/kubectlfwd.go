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

// Package kubectlfwd analyses incoming CLI call, and detects whether it should
// be hijacked and forwarded, as is, to kubectl. "create", "describe", "env" and
// "get" commands, if accompanied by an environment name, they are not trapped,
// and kube-cluster main execution course takes place.
//
// Whenever there is ambiguity when parsing the CLI commands and parameters, it
// will prefer kube-cluster over kubectl.
package kubectlfwd

import (
	"io"
	"os/exec"
	"strings"
)

const clusterObjectName = "env"

var acceptableClusterCalls = [...]struct {
	cmd string
	obj string
}{
	{"create-env", ""},
	{"describe", clusterObjectName},
	{"env", clusterObjectName},
	{"get", clusterObjectName},
}

// Fwd holds the CLI environment state which is used to make a decision about
// forwarding the call to kubectl.
type Fwd struct {
	args          []string // args expects pristine os.Args
	kubectlBinary string   // kubectlBinary is the full path of the detected kubectl

	// stdio for kubectl execution
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New instantiates a call forwarder (*Fwd). Feed it with os.Args, and os.Stdin,
// os.Stdout and os.Stderr.
func New(args []string, kubectl string, stdin io.Reader, stdout, stderr io.Writer) *Fwd {
	return &Fwd{
		args:          args,
		kubectlBinary: kubectl,

		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// Hijack effectively forwards the CLI call to kubectl, if the combination of
// command and objects are not targeting manipulation of Kubernetes environments.
func (f *Fwd) Hijack() (bool, error) {
	if f.isClusterCall() {
		return false, nil
	}

	cmd := exec.Command(f.kubectlBinary, f.args[1:]...)
	cmd.Stdin = f.stdin
	cmd.Stdout = f.stdout
	cmd.Stderr = f.stderr

	if err := cmd.Start(); err != nil {
		return true, err
	}

	if err := cmd.Wait(); err != nil {
		return true, err
	}

	return true, nil
}

// isClusterCall iterates through command and object pairs looking for
// non-hijackable calls.
func (f *Fwd) isClusterCall() bool {
	if len(f.args) == 1 || len(f.args) == 2 && '-' == f.args[1][0] {
		return true
	}
	if len(f.args) < 3 {
		return false
	}

	cmd := strings.ToLower(f.args[1])
	obj := strings.ToLower(f.args[2])
	for _, call := range acceptableClusterCalls {
		if cmd == call.cmd && (obj == call.obj || "" == call.obj) {
			return true
		}
	}

	return false
}
