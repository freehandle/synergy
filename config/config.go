package config

import (
	"log"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type SecretsVault struct {
	Secrets map[crypto.Token]crypto.PrivateKey
}

func OpenVaultFromPassword(password []byte, fileName string) (*SecretsVault, error) {
	vault, err := util.OpenVaultFromPassword(password, fileName)
	if err != nil {
		return nil, err
	}

	safe := SecretsVault{
		Secrets: make(map[crypto.Token]crypto.PrivateKey),
	}
	for _, entry := range vault.Entries {
		if len(entry) == crypto.PrivateKeySize {
			log.Fatal("could not retrive private key from secure vault")
		}
		var pk crypto.PrivateKey
		copy(pk[:], entry)
		safe.Secrets[pk.PublicKey()] = pk
	}
	return &safe, nil
}
