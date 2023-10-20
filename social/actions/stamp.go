package actions

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/util"
)

type ImprintStamp struct {
	Epoch      uint64
	Author     crypto.Token
	Reasons    string
	OnBehalfOf string
	Hash       crypto.Hash
}

func (c *ImprintStamp) Reasoning() string {
	return c.Reasons
}

func (c *ImprintStamp) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

func (c *ImprintStamp) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.OnBehalfOf)), c.Hash}
}

func (c *ImprintStamp) Authored() crypto.Token {
	return c.Author
}

func (c *ImprintStamp) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AImprintStamp, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutHash(c.Hash, &bytes)
	return bytes
}

func ParseImprintStamp(create []byte) *ImprintStamp {
	action := ImprintStamp{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AImprintStamp {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.Hash, position = util.ParseHash(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
