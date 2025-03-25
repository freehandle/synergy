package config

import (
	"log"
	"os"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type SecretsVault struct {
	Secrets map[crypto.Token]crypto.PrivateKey
	PK      crypto.PrivateKey
}

func OpenVaultFromPassword(password []byte, fileName string) (*SecretsVault, error) {
	vault, err := util.OpenVaultFromPassword(password, fileName)
	if err != nil {
		if _, fileErr := os.Stat(fileName); os.IsNotExist(fileErr) {
			vault, err = util.NewSecureVault(password, fileName)
			if err != nil {
				log.Fatal("could not create secure vault")
				return nil, err
			}
		} else {
			log.Fatal("could not open secure vault")
			return nil, err
		}
	}
	safe := SecretsVault{
		Secrets: make(map[crypto.Token]crypto.PrivateKey),
	}
	safe.PK = vault.SecretKey
	safe.Secrets[vault.SecretKey.PublicKey()] = vault.SecretKey
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
