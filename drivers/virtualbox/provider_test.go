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
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// some code here to check arguments perhaps?
	fmt.Fprintln(os.Stdout, os.Args)
	os.Exit(0)
}

func TestIsMinimumVirtualBoxVersion(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	v := new(Virtualbox)
	v.mgmtbin = "mgmtbin"
	err := v.isMinimumVirtualBoxVersion()
	if err != ErrParsingVirtualBoxVersion {
		t.Error("expect %v. got %v", ErrParsingVirtualBoxVersion, err)
	}
}

func TestIsISOdownloaded(t *testing.T) {
	const fn = "foo"
	v, err := New("testEnv", VagrantBox(fn))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if v.isISOdownloaded() {
		t.Error("unexpected file found: %v", fn)
	}
}
