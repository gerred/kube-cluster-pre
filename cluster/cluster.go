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
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Cluster struct {
	User     string
	Password string

	KubeletToken string
	ProxyToken   string
}

func (c *Cluster) GenerateBasicAuth() {
	c.User = "admin"
	c.Password = generatePassword(16)
}

func generatePassword(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const letterIdxBits = 6
	const letterIdxMask = 1<<letterIdxBits - 1
	const letterIdxMax = 63 / letterIdxBits

	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (c *Cluster) GenerateTokens() error {
	kubeletToken, err := c.readRandomToken()
	if err != nil {
		return err
	}

	proxyToken, err := c.readRandomToken()
	if err != nil {
		return err
	}

	c.KubeletToken = kubeletToken
	c.ProxyToken = proxyToken

	return nil
}

func (c *Cluster) readRandomToken() (string, error) {
	b := make([]byte, 128)

	if _, err := cryptoRand.Read(b); err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(b)
	filters := [...]string{"=", "+", "/"}
	for _, fchr := range filters {
		token = strings.Replace(token, fchr, "", -1)
	}
	return token[0:32], nil
}

func (c *Cluster) IsValid() bool {
	// todo(carlos): test for cluster correctness.

	return true
}

func (c *Cluster) Info() string {
	return fmt.Sprintf("%#v\n", c)
}
