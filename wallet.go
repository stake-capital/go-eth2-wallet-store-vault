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
	"encoding/json"
	"log"

	b64 "encoding/base64"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type WalletListSecret struct {
	keys []string
}

type WalletSecret struct {
	data []byte
}

// StoreWallet stores wallet-level data.  It will fail if it cannot store the data.
// Note that this will overwrite any existing data; it is up to higher-level functions to check for the presence of a wallet with
// the wallet name and handle clashes accordingly.
func (s *Store) StoreWallet(id uuid.UUID, name string, data []byte) error {
	path := s.walletHeaderPath(id)
	var err error
	data, err = s.encryptIfRequired(data)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt wallet")
	}

	sEnc := b64.URLEncoding.EncodeToString(data)
	_, err = s.client.KVv2(s.vault_secrets_mount_path).Put(context.Background(), path, map[string]interface{}{
		"data": sEnc,
	})
	if err != nil {
		return errors.Wrap(err, "failed to store wallet")
	}
	return nil
}

// RetrieveWallet retrieves wallet-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveWallet(walletName string) ([]byte, error) {
	for data := range s.RetrieveWallets() {
		info := &struct {
			Name string `json:"name"`
		}{}
		err := json.Unmarshal(data, info)
		if err == nil && info.Name == walletName {
			return data, nil
		}
	}
	return nil, errors.New("wallet not found")
}

// RetrieveWalletByID retrieves wallet-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
	for data := range s.RetrieveWallets() {
		info := &struct {
			ID uuid.UUID `json:"uuid"`
		}{}
		err := json.Unmarshal(data, info)
		if err == nil && info.ID == walletID {
			return data, nil
		}
	}
	return nil, errors.New("wallet not found")
}

// RetrieveWallets retrieves wallet-level data for all wallets.
func (s *Store) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		endpoint := "/" + s.vault_secrets_mount_path + "/metadata/wallets"
		walletList, err := s.client.Logical().List(endpoint)

		if err != nil {
			log.Fatalln(err)
		}

		if err == nil && walletList != nil && walletList.Data != nil {
			k := walletList.Data["keys"].([]interface{})
			if k != nil {
				for _, item := range k {
					walletIdWithSuffix := item.(string)
					lastCharPos := len(walletIdWithSuffix) - 1
					walletId := walletIdWithSuffix[:lastCharPos]

					uuidId, _ := uuid.Parse(walletId)
					path := s.walletHeaderPath(uuidId)

					secret, err := s.client.KVv2(s.vault_secrets_mount_path).Get(context.Background(), path)
					if err != nil {
						continue
					}

					returnedData, _ := secret.Data["data"].(string)

					sDec, _ := b64.URLEncoding.DecodeString(returnedData)
					data, err := s.decryptIfRequired([]byte(sDec))
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
