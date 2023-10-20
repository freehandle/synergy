package actions

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type Draft struct {
	Epoch         uint64
	Author        crypto.Token
	Reasons       string
	OnBehalfOf    string
	CoAuthors     []crypto.Token
	Policy        *Policy
	Title         string
	Keywords      []string
	Description   string
	ContentType   string
	ContentHash   crypto.Hash // hash of the entire content, not of the part
	NumberOfParts byte
	Content       []byte // entire content of the first part
	PreviousDraft crypto.Hash
	References    []crypto.Hash
}

func (c *Draft) Reasoning() string {
	return c.Reasons
}

func (c *Draft) Hashed() crypto.Hash {
	return c.ContentHash
}

// Se for em nome de um coletivo, afeta o coletivo
func (c *Draft) Affected() []crypto.Hash {
	if c.OnBehalfOf != "" {
		return []crypto.Hash{crypto.Hasher([]byte(c.OnBehalfOf))}
	}
	return nil
}

func (c *Draft) Authored() crypto.Token {
	return c.Author
}

func (c *Draft) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ADraft, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	PutTokenArray(c.CoAuthors, &bytes)
	if c.Policy != nil {
		util.PutByte(1, &bytes) // there is a policy
		PutPolicy(*c.Policy, &bytes)
	} else {
		util.PutByte(0, &bytes) // there is no policy
	}
	util.PutString(c.Title, &bytes)
	PutKeywords(c.Keywords, &bytes)
	util.PutString(c.Description, &bytes)
	util.PutString(c.ContentType, &bytes)
	util.PutHash(c.ContentHash, &bytes)
	util.PutByte(c.NumberOfParts, &bytes)
	util.PutByteArray(c.Content, &bytes)
	util.PutHash(c.PreviousDraft, &bytes)
	PutHashArray(c.References, &bytes)
	return bytes
}

func ParseDraft(create []byte) *Draft {
	action := Draft{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ADraft {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.CoAuthors, position = ParseTokenArray(create, position)
	if create[position] == 1 {
		var policy Policy
		position += 1
		policy, position = ParsePolicy(create, position)
		action.Policy = &policy
	} else if create[position] != 0 {
		return nil
	} else {
		position += 1
	}
	action.Title, position = util.ParseString(create, position)
	action.Keywords, position = ParseKeywords(create, position)
	action.Description, position = util.ParseString(create, position)
	action.ContentType, position = util.ParseString(create, position)
	action.ContentHash, position = util.ParseHash(create, position)
	action.NumberOfParts, position = util.ParseByte(create, position)
	action.Content, position = util.ParseByteArray(create, position)
	action.PreviousDraft, position = util.ParseHash(create, position)
	action.References, position = ParseHashArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type ReleaseDraft struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	ContentHash crypto.Hash
}

func (c *ReleaseDraft) Reasoning() string {
	return c.Reasons
}

func (c *ReleaseDraft) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// Release afeta apenas o proprio draft
func (c *ReleaseDraft) Affected() []crypto.Hash {
	return []crypto.Hash{c.ContentHash}
}

func (c *ReleaseDraft) Authored() crypto.Token {
	return c.Author
}

func (c *ReleaseDraft) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AReleaseDraft, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.ContentHash, &bytes)
	return bytes
}

func ParseReleaseDraft(create []byte) *ReleaseDraft {
	action := ReleaseDraft{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AReleaseDraft {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.ContentHash, position = util.ParseHash(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
