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

package cli

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a Kubernetes cluster with the given name and provider options",
	Run:   CreateCluster,
}

//CreateCluster creates a new Kubernetes cluster with a provider name and options.
func CreateCluster(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Error("name needs to be provided")
	}
}
