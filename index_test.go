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
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	vault "github.com/stake-capital/go-eth2-wallet-store-vault"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/go-indexer"
)

func TestStoreRetrieveIndex(t *testing.T) {
	rand.Seed(time.Now().Unix())
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
	walletName := "test wallet"
	walletData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID.String()))
	accountID := uuid.New()
	accountName := "test account"
	accountData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))

	index := indexer.New()
	index.Add(accountID, accountName)

	err = store.StoreWallet(walletID, walletName, walletData)
	require.Nil(t, err)
	err = store.StoreAccount(walletID, accountID, accountData)
	require.Nil(t, err)

	serializedIndex, err := index.Serialize()
	require.Nil(t, err)
	err = store.StoreAccountsIndex(walletID, serializedIndex)
	require.Nil(t, err)

	fetchedIndex, err := store.RetrieveAccountsIndex(walletID)
	require.Nil(t, err)

	reIndex, err := indexer.Deserialize(fetchedIndex)
	require.Nil(t, err)

	fetchedAccountName, exists := reIndex.Name(accountID)
	require.Equal(t, true, exists)
	require.Equal(t, accountName, fetchedAccountName)

	fetchedAccountID, exists := reIndex.ID(accountName)
	require.Equal(t, true, exists)
	require.Equal(t, accountID, fetchedAccountID)
}
