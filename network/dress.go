package network

import (
	"github.com/freehandle/breeze/crypto"
)

const (
	breezeTailSize = crypto.SignatureSize + 8 + crypto.TokenSize
	axeTailsize    = breezeTailSize + crypto.TokenSize + crypto.SignatureSize
)

// Breeze Void + Axé Void Specification
// version for breeze           (byte)           0
// void breeze instruction      (byte)           1
// Epoch                        (8 bytes)        2
// Protocol Code                (4 bytes)        10
// Axé Void instruction code    (byte)           14
// Author                       (32 bytes)       15
// Data ....                    (Variable)
// Signer                       (32 bytes)
// Signature                    (64 bytes)
// Wallet                       (32 bytes)
// Fee                          (8 bytes)
// Signature                    (64 bytes)

// Translate breeze byte array into synergy byte array
func BreezeToSynergy(action []byte) []byte {
	if len(action) < 15+axeTailsize {
		return nil
	}
	// strip first 2 bytes, the 4 bytes of protocol, the byte for the axe void and
	// the tail (signer ... wallet signayture)
	return append(action[2:10], action[15:len(action)-axeTailsize]...)

}
