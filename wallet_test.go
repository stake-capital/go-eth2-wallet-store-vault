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
	"fmt"
	"testing"

	vault "github.com/bliiitz/go-eth2-wallet-store-vault"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreRetrieveWallet(t *testing.T) {
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

	walletID := uuid.New()
	walletName := uuid.New()
	data := []byte(fmt.Sprintf(`{"uuid":%q,"name":%q}`, walletID, walletName.String()))

	err = store.StoreWallet(walletID, walletName.String(), data)
	require.Nil(t, err)
	retData, err := store.RetrieveWallet(walletName.String())
	require.Nil(t, err)
	assert.Equal(t, data, retData)

	for range store.RetrieveWallets() {
	}
}
