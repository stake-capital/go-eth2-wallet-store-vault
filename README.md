# go-eth2-wallet-store-vault
<!-- 
[![Tag](https://img.shields.io/github/tag/wealdtech/go-eth2-wallet-store-s3.svg)](https://github.com/wealdtech/go-eth2-wallet-store-s3/releases/)
[![License](https://img.shields.io/github/license/wealdtech/go-eth2-wallet-store-s3.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/wealdtech/go-eth2-wallet-store-s3?status.svg)](https://godoc.org/github.com/wealdtech/go-eth2-wallet-store-s3)
[![Travis CI](https://img.shields.io/travis/wealdtech/go-eth2-wallet-store-s3.svg)](https://travis-ci.org/wealdtech/go-eth2-wallet-store-s3)
[![codecov.io](https://img.shields.io/codecov/c/github/wealdtech/go-eth2-wallet-store-s3.svg)](https://codecov.io/github/wealdtech/go-eth2-wallet-store-s3)
[![Go Report Card](https://goreportcard.com/badge/github.com/wealdtech/go-eth2-wallet-store-s3)](https://goreportcard.com/report/github.com/wealdtech/go-eth2-wallet-store-s3) -->

Hashicorp Vault based store for the [Ethereum 2 wallet](https://github.com/wealdtech/go-eth2-wallet).


## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-eth2-wallet-store-vault` is a standard Go module which can be installed with:

```sh
go get github.com/wealdtech/go-eth2-wallet-store-vault
```

## Usage

In normal operation this module should not be used directly.  Instead, it should be configured to be used as part of [go-eth2-wallet](https://github.com/wealdtech/go-eth2-wallet).

The Vault store has the following options:

  - `vault_addr`: the Vault address in which the wallet is to be stored. Exemple: http://localhost:8200 for local vault
  - `id`: an ID that is used to differentiate multiple stores created by the same account.  If this is not configured an empty ID is used
  - `vault_auth`: Vault authentication type. Values: `token` or `kubernetes`
  - `vault_token`: Vault token to use for requesting vault (Mandatory if `vault_auth` is `token`)
  - `vault_k8s_auth_role`: Name of the kubernetes auth role to use (Mandatory if `vault_auth` is `kubernetes`)
  - `vault_k8s_auth_sa_token_path`: Local path to access to the kubernetes service account token. Default: `/var/run/secrets/kubernetes.io/serviceaccount/token`
  - `vault_k8s_auth_mount_path`: Kubernetes auth module path. Default: `kubernetes`
  - `vault_secrets_mount_path`: KVv2 secrets module path (Mandatory)
  - `passphrase`: a key used to encrypt all data written to the store.  If this is not configured data is written to the store unencrypted (although wallet- and account-specific private information may be protected by their own passphrases)

When initiating a connection to Amazon S3 the Amazon credentials are required.  Details on how to make the credentials available to the store are available at [the Amazon S3 documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#shared-credentials-file)

### Example

```go
package main

import (
	e2wallet "github.com/wealdtech/go-eth2-wallet"
	vault "github.com/bliiitz/go-eth2-wallet-store-vault"
)

func main() {
    // Set up and use an encrypted store
    store, err := vault.New(vault.WithPassphrase([]byte("my secret")))
    if err != nil {
        panic(err)
    }
    e2wallet.UseStore(store)

    // Set up and use an encrypted store in the central Canada region
    store, err = vault.New(vault.WithPassphrase([]byte("my secret")), vault.WithRegion("ca-central-1"))
    if err != nil {
        panic(err)
    }
    e2wallet.UseStore(store)

    // Set up and use an encrypted store with a custom ID
    store, err = vault.New(vault.WithPassphrase([]byte("my secret")), vault.WithID([]byte("store 2")))
    if err != nil {
        panic(err)
    }
    e2wallet.UseStore(store)
}
```

## Maintainers

Bliiitz: [@mcdee](https://github.com/bliiitz).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/wealdtech/go-eth2-wallet-store-vault/issues).

## License

[Apache-2.0](LICENSE) © 2019 Bliiitz 
