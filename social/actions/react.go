package actions

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/util"
)

type React struct {
	Epoch      uint64
	Author     crypto.Token
	Reasons    string
	OnBehalfOf string
	Hash       crypto.Hash
	Reaction   byte
}

func (c *React) Reasoning() string {
	return c.Reasons
}

func (c *React) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// Afeta o objeto a que se aplica o react
func (c *React) Affected() []crypto.Hash {
	return []crypto.Hash{c.Hash}
}

func (c *React) Authored() crypto.Token {
	return c.Author
}

func (c *React) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AReact, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutHash(c.Hash, &bytes)
	util.PutByte(c.Reaction, &bytes)
	return bytes
}

func ParseReact(create []byte) *React {
	action := React{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AReact {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.Hash, position = util.ParseHash(create, position)
	action.Reaction, position = util.ParseByte(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
