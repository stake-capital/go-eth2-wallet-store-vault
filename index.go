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

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type WalletIndexSecret struct {
	data []byte
}

// StoreAccountsIndex stores the account index.
func (s *Store) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
	var err error

	// Do not encrypt empty index.
	if len(data) != 2 {
		data, err = s.encryptIfRequired(data)
		if err != nil {
			return err
		}
	}

	path := s.walletIndexPath(walletID)

	sEnc := b64.URLEncoding.EncodeToString(data)
	s.client.KVv2(s.vault_secrets_mount_path).Put(context.Background(), path, map[string]interface{}{
		"data": sEnc,
	})

	if err != nil {
		return errors.Wrap(err, "failed to store wallet index")
	}

	return nil
}

// RetrieveAccountsIndex retrieves the account index.
func (s *Store) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
	path := s.walletIndexPath(walletID)

	secret, err := s.client.KVv2(s.vault_secrets_mount_path).Get(context.Background(), path)
	if err != nil {
		return nil, err
	}

	returnedData, _ := secret.Data["data"].(string)

	sDec, _ := b64.URLEncoding.DecodeString(returnedData)
	// Do not decrypt empty index.
	if len(sDec) == 2 {
		return sDec, nil
	}
	data, err := s.decryptIfRequired(sDec)
	if err != nil {
		return nil, err
	}
	return data, nil
}
