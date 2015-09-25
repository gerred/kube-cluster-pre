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

import (
	"github.com/gerred/kube-cluster/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gerred/kube-cluster/Godeps/_workspace/src/github.com/spf13/viper"
)

type Option func(*Cluster)

var Flags = [...]struct {
	Name, defaultValue, description string

	Action func(string) Option
}{
	{"username", "admin", "specify username for Kubernetes cluster", Username},
	{"password", "", "specify password for Kubernetes cluster", Password},
}

func GetOptions(cmd *cobra.Command) {
	for _, v := range Flags {
		cmd.Flags().String(v.Name, v.defaultValue, v.description)
		viper.BindPFlag(v.Name, cmd.Flags().Lookup(v.Name))
	}
}

func Username(u string) Option {
	return func(c *Cluster) {
		c.username = u
	}
}

func Password(p string) Option {
	return func(c *Cluster) {
		c.password = p
	}
}
