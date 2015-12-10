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

package cluster

import "testing"

func TestGenerateBasicAuth(t *testing.T) {
	c := new(Cluster)
	c.GenerateBasicAuth()

	if c.username == "" || c.password == "" {
		t.Error("expected some username and password")
	}
}

func TestGenerateTokens(t *testing.T) {
	c := new(Cluster)
	c.GenerateTokens()

	if c.kubeletToken == "" || c.proxyToken == "" {
		t.Error("expected kubete and proxy tokens to not be empty")
	}
}
