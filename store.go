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

package vault

import (
	"log"
	"context"
	"errors"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"

	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// options are the options for the S3 store
type options struct {
	id          					[]byte
	vault_addr  					string
	vault_auth  					string
	vault_token  					string
	vault_k8s_auth_role  			string
	vault_k8s_auth_sa_token_path  	string
	vault_k8s_auth_mount_path  		string
	vault_secrets_mount_path    	string
	passphrase  					[]byte
}

// Option gives options to New
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithPassphrase sets the passphrase for the store.
func WithPassphrase(passphrase []byte) Option {
	return optionFunc(func(o *options) {
		o.passphrase = passphrase
	})
}

// WithID sets the ID for the store
func WithID(t []byte) Option {
	return optionFunc(func(o *options) {
		o.id = t
	})
}

// Store is the store for the wallet held encrypted on Amazon S3.
type Store struct {
	client      					*vault.Client
	id          					[]byte
	vault_addr  					string
	vault_auth  					string
	vault_token  					string
	vault_k8s_auth_role  			string
	vault_k8s_auth_sa_token_path  	string
	vault_k8s_auth_mount_path  		string
	vault_secrets_mount_path  		string
	passphrase  					[]byte
}

// New creates a new Amazon S3 store.
// This takes the following options:
//  - region: a string specifying the Amazon S3 region, defaults to "us-east-1", set with WithRegion()
//  - id: a byte array specifying an identifying key for the store, defaults to nil, set with WithID()
// This expects the access credentials to be in a standard place, e.g. ~/.aws/credentials
func New(opts ...Option) (wtypes.Store, error) {
	options := options{
		vault_addr: 					"",
		vault_auth: 					"",
		vault_token: 					"",
		vault_k8s_auth_role: 			"",
		vault_k8s_auth_sa_token_path: 	"/var/run/secrets/kubernetes.io/serviceaccount/token",
		vault_k8s_auth_mount_path:		"kubernetes",
		vault_secrets_mount_path: 		"",
	}
	for _, o := range opts {
		o.apply(&options)
	}

	if options.vault_addr == "" {
		return nil, errors.New("vault_addr option missing")
	}

	if options.vault_auth == "" {
		return nil, errors.New("vault_auth option missing")
	}

	if options.vault_secrets_mount_path == "" {
		return nil, errors.New("vault_secrets_mount_path option missing")
	}

	if options.vault_auth == "token" && options.vault_token == "" {
		return nil, errors.New("vault_token option missing")
	}

	if options.vault_auth == "kubernetes" && options.vault_k8s_auth_role == "" {
		return nil, errors.New("vault_k8s_auth_role option missing")
	}

	// If set, the VAULT_ADDR environment variable will be the address that
	// your pod uses to communicate with Vault.
	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = options.vault_addr

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	if options.vault_auth == "token" {
		client.SetToken(options.vault_token)
	}

	if options.vault_auth == "kubernetes" {
		// The service-account token will be read from the path where the token's
		// Kubernetes Secret is mounted. By default, Kubernetes will mount it to
		// /var/run/secrets/kubernetes.io/serviceaccount/token, but an administrator
		// may have configured it to be mounted elsewhere.
		// In that case, we'll use the option WithServiceAccountTokenPath to look
		// for the token there.
		k8sAuth, err := auth.NewKubernetesAuth(
			options.vault_k8s_auth_role,
			auth.WithMountPath(options.vault_k8s_auth_mount_path),
			auth.WithServiceAccountTokenPath(options.vault_k8s_auth_sa_token_path),
		)
		if err != nil {
			return nil, err
		}

		authInfo, err := client.Auth().Login(context.Background(), k8sAuth)
		if err != nil {
			return nil, err
		}
		if authInfo == nil {
			log.Fatal("no auth info was returned after login")
		}
		client.SetToken(authInfo.Auth.ClientToken)
	}
	

	return &Store{
		client:     					client,
		id:         					options.id,
		vault_addr:  					options.vault_addr,
		vault_auth:  					options.vault_auth,
		vault_token:  					options.vault_token,
		vault_k8s_auth_role:  			options.vault_k8s_auth_role,
		vault_k8s_auth_sa_token_path:  	options.vault_k8s_auth_sa_token_path,
		vault_k8s_auth_mount_path:  	options.vault_k8s_auth_mount_path,
		vault_secrets_mount_path:  		options.vault_secrets_mount_path,
		passphrase: 					options.passphrase,
	}, nil
}

// Name returns the name of this store.
func (s *Store) Name() string {
	return "vault"
}

// Location returns the location of this store.
func (s *Store) Location() string {
	return s.vault_addr
}
