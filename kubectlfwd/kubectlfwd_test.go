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

import "testing"

func TestIsClusterCall(t *testing.T) {
	argsSet := [...]struct {
		args   []string
		result bool
	}{
		{[]string{"kube-cluster"}, true},
		{[]string{"kube-cluster", "--help"}, true},
		{[]string{"kube-cluster", "create-env"}, false},
		{[]string{"kube-cluster", "create-env", "dev"}, true},
		{[]string{"kube-cluster", "create-env", "--help"}, true},
		{[]string{"kube-cluster", "get", "env"}, true},
		{[]string{"kube-cluster", "get", "pod"}, false},
		{[]string{"kube-cluster", "get", "pod", "some-pod"}, false},
	}

	for _, set := range argsSet {
		fwd := New(
			set.args,
			"kubectl",

			nil,
			nil,
			nil,
		)

		if fwd.isClusterCall() != set.result {
			t.Errorf("unexpected behavior found. %#v", set.args)
		}
	}
}
