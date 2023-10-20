package actions

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type CreateCollective struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	Name        string
	Description string
	Policy      Policy
}

func (c *CreateCollective) Reasoning() string {
	return c.Reasons
}

func (c *CreateCollective) Hashed() crypto.Hash {
	return crypto.Hasher([]byte(c.Name))
}

// Afeta um hash raiz (?)
func (c *CreateCollective) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.ZeroHash}
}

func (c *CreateCollective) Authored() crypto.Token {
	return c.Author
}

func (c *CreateCollective) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACreateCollective, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Name, &bytes)
	util.PutString(c.Description, &bytes)
	PutPolicy(c.Policy, &bytes)
	return bytes
}

func ParseCreateCollective(create []byte) *CreateCollective {
	action := CreateCollective{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACreateCollective {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Name, position = util.ParseString(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Policy, position = ParsePolicy(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type UpdateCollective struct {
	Epoch         uint64
	Author        crypto.Token
	Reasons       string
	OnBehalfOf    string
	Description   *string
	Majority      *byte
	SuperMajority *byte
}

func (c *UpdateCollective) Reasoning() string {
	return c.Reasons
}

func (c *UpdateCollective) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// Afeta o coletivo que vai ser atualizado
func (c *UpdateCollective) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.OnBehalfOf))}
}

func (c *UpdateCollective) Authored() crypto.Token {
	return c.Author
}

func (c *UpdateCollective) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AUpdateCollective, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	if c.Description != nil {
		util.PutByte(1, &bytes)
		util.PutString(*c.Description, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Majority != nil {
		util.PutByte(1, &bytes) // there is a policy
		util.PutByte(byte(*c.Majority), &bytes)
	} else {
		util.PutByte(0, &bytes) // there is no policy
	}
	if c.SuperMajority != nil {
		util.PutByte(1, &bytes) // there is a policy
		util.PutByte(byte(*c.SuperMajority), &bytes)
	} else {
		util.PutByte(0, &bytes) // there is no policy
	}
	return bytes
}

func ParseUpdateCollective(update []byte) *UpdateCollective {
	action := UpdateCollective{}
	position := 0
	action.Epoch, position = util.ParseUint64(update, position)
	action.Author, position = util.ParseToken(update, position)
	if update[position] != AUpdateCollective {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(update, position)
	action.OnBehalfOf, position = util.ParseString(update, position)

	if update[position] == 1 {
		position += 1
		var des string
		des, position = util.ParseString(update, position)
		action.Description = &des
	} else if update[position] == 0 {
		position += 1
	} else {
		return nil
	}
	if update[position] == 1 {
		var majority byte
		position += 1
		majority, position = util.ParseByte(update, position)
		action.Majority = &majority
	} else if update[position] != 0 {
		return nil
	} else {
		position += 1
	}
	if update[position] == 1 {
		var superMajority byte
		position += 1
		superMajority, position = util.ParseByte(update, position)
		action.SuperMajority = &superMajority
	} else if update[position] != 0 {
		return nil
	} else {
		position += 1
	}
	if position != len(update) {
		return nil
	}
	return &action
}

type RequestMembership struct {
	Epoch      uint64
	Author     crypto.Token
	Reasons    string
	Collective string
	Include    bool
}

func (c *RequestMembership) Reasoning() string {
	return c.Reasons
}

func (c *RequestMembership) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// Afeta o coletivo
func (c *RequestMembership) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.Collective))}
}

func (c *RequestMembership) Authored() crypto.Token {
	return c.Author
}

func (c *RequestMembership) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ARequestMembership, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Collective, &bytes)
	util.PutBool(c.Include, &bytes)
	return bytes
}

func ParseRequestMembership(update []byte) *RequestMembership {
	action := RequestMembership{}
	position := 0
	action.Epoch, position = util.ParseUint64(update, position)
	action.Author, position = util.ParseToken(update, position)
	if update[position] != ARequestMembership {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(update, position)
	action.Collective, position = util.ParseString(update, position)
	action.Include, position = util.ParseBool(update, position)
	if position != len(update) {
		return nil
	}
	return &action
}

type RemoveMember struct {
	Epoch      uint64
	Author     crypto.Token
	OnBehalfOf string
	Reasons    string
	Member     crypto.Token
}

func (c *RemoveMember) Reasoning() string {
	return c.Reasons
}

func (c *RemoveMember) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// Afeta o coletivo
func (c *RemoveMember) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.OnBehalfOf))}
}

func (c *RemoveMember) Authored() crypto.Token {
	return c.Author
}

func (c *RemoveMember) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ARemoveMember, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutToken(c.Member, &bytes)
	return bytes
}

func ParseRemoveMember(update []byte) *RemoveMember {
	action := RemoveMember{}
	position := 0
	action.Epoch, position = util.ParseUint64(update, position)
	action.Author, position = util.ParseToken(update, position)
	if update[position] != ARemoveMember {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(update, position)
	action.OnBehalfOf, position = util.ParseString(update, position)
	action.Member, position = util.ParseToken(update, position)
	if position != len(update) {
		return nil
	}
	return &action
}
