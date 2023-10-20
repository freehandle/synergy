package actions

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type CreateBoard struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	OnBehalfOf  string
	Name        string
	Description string
	Keywords    []string
	PinMajority byte
}

func (c *CreateBoard) Reasoning() string {
	return c.Reasons
}

func (c *CreateBoard) Hashed() crypto.Hash {
	return crypto.Hasher([]byte(c.Name))
}

// afeta apenas o coletivo em nome do qual ta sendo criado o board
func (c *CreateBoard) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.OnBehalfOf))}
}

func (c *CreateBoard) Authored() crypto.Token {
	return c.Author
}

func (c *CreateBoard) Serialize() []byte {
	//bytes := []byte{0, breeze.IVoid}
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	//util.PutByteArray([]byte{1, 1, 0, 0}, &bytes)
	//util.PutByte(attorney.VoidType, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACreateBoard, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutString(c.Name, &bytes)
	util.PutString(c.Description, &bytes)
	PutKeywords(c.Keywords, &bytes)
	util.PutByte(c.PinMajority, &bytes)
	return bytes
}

func ParseCreateBoard(create []byte) *CreateBoard {
	action := CreateBoard{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACreateBoard {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.Name, position = util.ParseString(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Keywords, position = ParseKeywords(create, position)
	action.PinMajority, position = util.ParseByte(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type UpdateBoard struct {
	Epoch       uint64
	Author      crypto.Token
	Reasons     string
	Board       string
	Description *string
	Keywords    *[]string
	PinMajority *byte
}

func (c *UpdateBoard) Reasoning() string {
	return c.Reasons
}

func (c *UpdateBoard) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// afeta apenas o board
func (c *UpdateBoard) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.Board))}
}

func (c *UpdateBoard) Authored() crypto.Token {
	return c.Author
}

func (c *UpdateBoard) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AUpdateBoard, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	if c.Description != nil {
		util.PutByte(1, &bytes)
		util.PutString(*c.Description, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Keywords != nil {
		util.PutByte(1, &bytes)
		PutKeywords(*c.Keywords, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.PinMajority != nil {
		util.PutByte(1, &bytes)
		util.PutByte(*c.PinMajority, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	return bytes
}

func ParseUpdateBoard(create []byte) *UpdateBoard {
	action := UpdateBoard{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AUpdateBoard {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Board, position = util.ParseString(create, position)
	if create[position] == 0 {
		position += 1
	} else {
		var des string
		position += 1
		des, position = util.ParseString(create, position)
		action.Description = &des
	}
	if create[position] == 0 {
		position += 1
	} else {
		var key []string
		position += 1
		key, position = ParseKeywords(create, position)
		action.Keywords = &key
	}
	if create[position] == 0 {
		position += 1
	} else {
		var pin byte
		position += 1
		pin, position = util.ParseByte(create, position)
		action.PinMajority = &pin
	}
	if position != len(create) {
		return nil
	}
	return &action
}

type Pin struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Board   string
	Draft   crypto.Hash
	Pin     bool
}

func (c *Pin) Reasoning() string {
	return c.Reasons
}

func (c *Pin) Hashed() crypto.Hash {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutHash(c.Draft, &bytes)
	util.PutString(c.Board, &bytes)
	// checa se eh um pin ou um unpin
	if c.Pin {
		util.PutByte(1, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	hash := crypto.Hasher(bytes)
	return hash
}

// afeta o board e o draft
func (c *Pin) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.Board)), c.Draft}
}

func (c *Pin) Authored() crypto.Token {
	return c.Author
}

func (c *Pin) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes) // 8 bytes
	util.PutToken(c.Author, &bytes) // 32 bytes do token
	util.PutByte(APin, &bytes)      // 41o byte eh o byte da acao
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	util.PutHash(c.Draft, &bytes)
	util.PutBool(c.Pin, &bytes)
	return bytes
}

func ParsePin(create []byte) *Pin {
	action := Pin{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != APin {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Board, position = util.ParseString(create, position)
	action.Draft, position = util.ParseHash(create, position)
	action.Pin, position = util.ParseBool(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type BoardEditor struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Board   string
	Editor  crypto.Token
	Insert  bool
}

func (c *BoardEditor) Reasoning() string {
	return c.Reasons
}

func (c *BoardEditor) Hashed() crypto.Hash {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Editor, &bytes)
	util.PutString(c.Board, &bytes)
	if c.Insert {
		util.PutByte(1, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	hash := crypto.Hasher(bytes)
	return hash
}

// afeta o board
func (c *BoardEditor) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.Board))}
}

func (c *BoardEditor) Authored() crypto.Token {
	return c.Author
}

func (c *BoardEditor) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ABoardEditor, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.Board, &bytes)
	util.PutToken(c.Editor, &bytes)
	util.PutBool(c.Insert, &bytes)
	return bytes
}

func ParseBoardEditor(create []byte) *BoardEditor {
	action := BoardEditor{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ABoardEditor {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Board, position = util.ParseString(create, position)
	action.Editor, position = util.ParseToken(create, position)
	action.Insert, position = util.ParseBool(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
