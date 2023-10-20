package actions

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type MultipartMedia struct {
	Epoch  uint64
	Author crypto.Token
	Hash   crypto.Hash
	Part   byte
	Of     byte
	Data   []byte
}

func (c *MultipartMedia) Reasoning() string {
	return ""
}

func (c *MultipartMedia) Hashed() crypto.Hash {
	return c.Hash
}

func (c *MultipartMedia) Authored() crypto.Token {
	return c.Author
}

func (c *MultipartMedia) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AMultipartMedia, &bytes)
	util.PutHash(c.Hash, &bytes)
	util.PutByte(c.Part, &bytes)
	util.PutByte(c.Of, &bytes)
	util.PutByteArray(c.Data, &bytes)
	return bytes
}

func ParseMultipartMedia(create []byte) *MultipartMedia {
	action := MultipartMedia{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AMultipartMedia {
		return nil
	}
	position += 1
	action.Hash, position = util.ParseHash(create, position)
	action.Part, position = util.ParseByte(create, position)
	action.Of, position = util.ParseByte(create, position)
	action.Data, position = util.ParseByteArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
