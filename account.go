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

package vaultstorage

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	// Ensure the wallet exists
	_, err := s.RetrieveWalletByID(walletID)
	if err != nil {
		return errors.New("unknown wallet")
	}

	// See if an account with this name already exists
	existingAccount, err := s.RetrieveAccount(walletID, accountID)
	if err == nil {
		// It does; they need to have the same ID for us to overwrite it
		info := &struct {
			ID string `json:"uuid"`
		}{}
		err := json.Unmarshal(existingAccount, info)
		if err != nil {
			return err
		}
		if info.ID != accountID.String() {
			return errors.New("account already exists")
		}
	}

	data, err = s.encryptIfRequired(data)
	if err != nil {
		return err
	}

	path := s.accountPath(walletID, accountID)
	sEnc := b64.URLEncoding.EncodeToString(data)
	s.client.KVv2(s.vault_secrets_mount_path).Put(context.Background(), path, map[string]interface{}{
		"data": sEnc,
	})

	if err != nil {
		return errors.Wrap(err, "failed to store key")
	}
	return nil
}

// RetrieveAccount retrieves account-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
	path := s.accountPath(walletID, accountID)

	secret, err := s.client.KVv2(s.vault_secrets_mount_path).Get(context.Background(), path)
	if err != nil {
		return nil, err
	}
	returnedData, _ := secret.Data["data"].(string)
	sDec, _ := b64.URLEncoding.DecodeString(returnedData)
	data, err := s.decryptIfRequired(sDec)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	path := s.walletPath(walletID)
	ch := make(chan []byte, 1024)
	go func() {
		endpoint := "/" + s.vault_secrets_mount_path + "/metadata/wallets/" + path
		accountList, err := s.client.Logical().List(endpoint)
		if err == nil && accountList != nil && accountList.Data != nil {
			k := accountList.Data["keys"].([]interface{})
			if k != nil {
				for _, item := range k {
					account := item.(string)
					if strings.HasSuffix(account, "/") {
						// Directory
						continue
					}
					if strings.HasSuffix(account, walletID.String()) {
						// Wallet
						continue
					}

					uuidId, _ := uuid.Parse(account)
					secret, err := s.client.KVv2(s.vault_secrets_mount_path).Get(context.Background(), s.accountPath(walletID, uuidId))
					if err != nil {
						continue
					}

					returnedData, _ := secret.Data["data"].(string)

					sDec, _ := b64.URLEncoding.DecodeString(returnedData)
					data, err := s.decryptIfRequired(sDec)
					if err != nil {
						continue
					}
					ch <- data
				}
			}
		}
		close(ch)
	}()
	return ch
}
