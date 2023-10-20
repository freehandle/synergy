package actions

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/util"
)

type Edit struct {
	Epoch         uint64
	Author        crypto.Token
	Reasons       string
	OnBehalfOf    string
	CoAuthors     []crypto.Token
	EditedDraft   crypto.Hash
	ContentType   string
	ContentHash   crypto.Hash // hash of the entire content, not of the part
	NumberOfParts byte
	Content       []byte // entire content of the first part
}

func (c *Edit) Reasoning() string {
	return c.Reasons
}

func (c *Edit) Hashed() crypto.Hash {
	return c.ContentHash
}

// Se for em nome de um coletivo, afeta o coletivo e o draft, senao somente o draft
func (c *Edit) Affected() []crypto.Hash {
	if c.OnBehalfOf != "" {
		return []crypto.Hash{crypto.Hasher([]byte(c.OnBehalfOf)), c.ContentHash}
	}
	return []crypto.Hash{c.ContentHash}
}

func (c *Edit) Authored() crypto.Token {
	return c.Author
}

func (c *Edit) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AEdit, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	PutTokenArray(c.CoAuthors, &bytes)
	util.PutHash(c.EditedDraft, &bytes)
	util.PutString(c.ContentType, &bytes)
	util.PutHash(c.ContentHash, &bytes)
	util.PutByte(c.NumberOfParts, &bytes)
	util.PutByteArray(c.Content, &bytes)
	return bytes
}

func ParseEdit(create []byte) *Edit {
	action := Edit{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AEdit {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.CoAuthors, position = ParseTokenArray(create, position)
	action.EditedDraft, position = util.ParseHash(create, position)
	action.ContentType, position = util.ParseString(create, position)
	action.ContentHash, position = util.ParseHash(create, position)
	action.NumberOfParts, position = util.ParseByte(create, position)
	action.Content, position = util.ParseByteArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
