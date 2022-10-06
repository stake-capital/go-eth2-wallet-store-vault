// Copyright 2019, 2020 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vaultstorage_test

import (
	"testing"

	vault "github.com/bliiitz/go-eth2-wallet-store-vault"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
	// wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestNew(t *testing.T) {
	store, err := vault.New(
		vault.WithPassphrase([]byte("test")),
		vault.WithVaultAddr("http://localhost:8200"),
		vault.WithVaultSecretMountPath("secret"),
		vault.WithVaultToken("golang-test"),
		vault.WithVaultAuth("token"),
	)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "vault", store.Name())

}
