package actions

import (
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type CreateEvent struct {
	Epoch           uint64
	Author          crypto.Token
	Reasons         string
	OnBehalfOf      string // non-optional
	StartAt         time.Time
	EstimatedEnd    time.Time
	Description     string
	Venue           string
	Open            bool
	Public          bool
	ManagerMajority byte
	Managers        []crypto.Token // default é qualquer um do coletivo
}

func (c *CreateEvent) Reasoning() string {
	return c.Reasons
}

func (c *CreateEvent) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// afeta o coletivo em nome do qual esta sendo criado o evento
func (c *CreateEvent) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.Hasher([]byte(c.OnBehalfOf))}
}

func (c *CreateEvent) Authored() crypto.Token {
	return c.Author
}

func (c *CreateEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACreateEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutString(c.OnBehalfOf, &bytes)
	util.PutTime(c.StartAt, &bytes)
	util.PutTime(c.EstimatedEnd, &bytes)
	util.PutString(c.Description, &bytes)
	util.PutString(c.Venue, &bytes)
	util.PutBool(c.Open, &bytes)
	util.PutBool(c.Public, &bytes)
	util.PutByte(c.ManagerMajority, &bytes)
	PutTokenArray(c.Managers, &bytes)
	return bytes
}

func ParseCreateEvent(create []byte) *CreateEvent {
	action := CreateEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACreateEvent {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.OnBehalfOf, position = util.ParseString(create, position)
	action.StartAt, position = util.ParseTime(create, position)
	action.EstimatedEnd, position = util.ParseTime(create, position)
	action.Description, position = util.ParseString(create, position)
	action.Venue, position = util.ParseString(create, position)
	action.Open, position = util.ParseBool(create, position)
	action.Public, position = util.ParseBool(create, position)
	action.ManagerMajority, position = util.ParseByte(create, position)
	action.Managers, position = ParseTokenArray(create, position)

	if position != len(create) {
		return nil
	}
	return &action
}

type CancelEvent struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Hash    crypto.Hash
}

func (c *CancelEvent) Reasoning() string {
	return c.Reasons
}

func (c *CancelEvent) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// afeta apenas o proprio evento
func (c *CancelEvent) Affected() []crypto.Hash {
	return []crypto.Hash{c.Hashed()}
}

func (c *CancelEvent) Authored() crypto.Token {
	return c.Author
}

func (c *CancelEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACancelEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.Hash, &bytes)
	return bytes
}

func ParseCancelEvent(create []byte) *CancelEvent {
	action := CancelEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACancelEvent {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.Hash, position = util.ParseHash(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type UpdateEvent struct {
	Epoch           uint64
	Author          crypto.Token
	Reasons         string
	EventHash       crypto.Hash
	StartAt         *time.Time
	EstimatedEnd    *time.Time
	Description     *string
	Venue           *string
	Open            *bool
	Public          *bool
	ManagerMajority *byte
	Managers        *[]crypto.Token
}

func (c *UpdateEvent) Reasoning() string {
	return c.Reasons
}

func (c *UpdateEvent) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// afeta apenas o proprio evento
func (c *UpdateEvent) Affected() []crypto.Hash {
	return []crypto.Hash{c.EventHash}
}

func (c *UpdateEvent) Authored() crypto.Token {
	return c.Author
}

func (c *UpdateEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AUpdateEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	if c.StartAt != nil {
		util.PutByte(1, &bytes)
		util.PutTime(*c.StartAt, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.EstimatedEnd != nil {
		util.PutByte(1, &bytes)
		util.PutTime(*c.EstimatedEnd, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Description != nil {
		util.PutByte(1, &bytes)
		util.PutString(*c.Description, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Venue != nil {
		util.PutByte(1, &bytes)
		util.PutString(*c.Venue, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Open != nil {
		util.PutByte(1, &bytes)
		util.PutBool(*c.Open, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Public != nil {
		util.PutByte(1, &bytes)
		util.PutBool(*c.Public, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.ManagerMajority != nil {
		util.PutByte(1, &bytes)
		util.PutByte(*c.ManagerMajority, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	if c.Managers != nil {
		util.PutByte(1, &bytes)
		PutTokenArray(*c.Managers, &bytes)
	} else {
		util.PutByte(0, &bytes)
	}
	return bytes
}

func ParseUpdateEvent(create []byte) *UpdateEvent {
	action := UpdateEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AUpdateEvent {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.EventHash, position = util.ParseHash(create, position)
	if create[position] == 0 {
		position += 1
	} else {
		var t time.Time
		position += 1
		t, position = util.ParseTime(create, position)
		action.StartAt = &t
	}
	if create[position] == 0 {
		position += 1
	} else {
		var t time.Time
		position += 1
		t, position = util.ParseTime(create, position)
		action.EstimatedEnd = &t
	}
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
		var venue string
		position += 1
		venue, position = util.ParseString(create, position)
		action.Venue = &venue
	}
	if create[position] == 0 {
		position += 1
	} else {
		var b bool
		position += 1
		b, position = util.ParseBool(create, position)
		action.Open = &b
	}
	if create[position] == 0 {
		position += 1
	} else {
		var b bool
		position += 1
		b, position = util.ParseBool(create, position)
		action.Public = &b
	}
	if create[position] == 0 {
		position += 1
	} else {
		var b byte
		position += 1
		b, position = util.ParseByte(create, position)
		action.ManagerMajority = &b
	}
	if create[position] == 0 {
		position += 1
	} else {
		var tokens []crypto.Token
		position += 1
		tokens, position = ParseTokenArray(create, position)
		action.Managers = &tokens
	}
	if position != len(create) {
		return nil
	}
	return &action
}

type CheckinEvent struct {
	Epoch          uint64
	Author         crypto.Token
	EphemeralToken crypto.Token
	Reasons        string
	EventHash      crypto.Hash
}

func (c *CheckinEvent) Reasoning() string {
	return c.Reasons
}

func (c *CheckinEvent) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

// afeta apenas o proprio evento
func (c *CheckinEvent) Affected() []crypto.Hash {
	return []crypto.Hash{c.EventHash}
}

func (c *CheckinEvent) Authored() crypto.Token {
	return c.Author
}

func (c *CheckinEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(ACheckinEvent, &bytes)
	util.PutToken(c.EphemeralToken, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	return bytes
}

func ParseCheckinEvent(create []byte) *CheckinEvent {
	action := CheckinEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != ACheckinEvent {
		return nil
	}
	position += 1
	action.EphemeralToken, position = util.ParseToken(create, position)
	action.Reasons, position = util.ParseString(create, position)
	action.EventHash, position = util.ParseHash(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}

type GreetCheckinEvent struct {
	Epoch          uint64
	Author         crypto.Token
	Reasons        string
	EventHash      crypto.Hash
	CheckedIn      crypto.Token
	EphemeralToken crypto.Token
	SecretKey      []byte // diffie-hellman
	PrivateContent []byte
}

func (c *GreetCheckinEvent) Reasoning() string {
	return c.Reasons
}

func (g *GreetCheckinEvent) Hashed() crypto.Hash {
	return crypto.Hasher(g.Serialize())
}

// afeta apenas o proprio evento
func (c *GreetCheckinEvent) Affected() []crypto.Hash {
	return []crypto.Hash{c.EventHash}
}

func (c *GreetCheckinEvent) Authored() crypto.Token {
	return c.Author
}

func (c *GreetCheckinEvent) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoch, &bytes)
	util.PutToken(c.Author, &bytes)
	util.PutByte(AGreetCheckinEvent, &bytes)
	util.PutString(c.Reasons, &bytes)
	util.PutHash(c.EventHash, &bytes)
	util.PutToken(c.CheckedIn, &bytes)
	util.PutToken(c.EphemeralToken, &bytes)
	util.PutByteArray(c.SecretKey, &bytes)
	util.PutByteArray(c.PrivateContent, &bytes)
	return bytes
}

func ParseGreetCheckinEvent(create []byte) *GreetCheckinEvent {
	action := GreetCheckinEvent{}
	position := 0
	action.Epoch, position = util.ParseUint64(create, position)
	action.Author, position = util.ParseToken(create, position)
	if create[position] != AGreetCheckinEvent {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	action.EventHash, position = util.ParseHash(create, position)
	action.CheckedIn, position = util.ParseToken(create, position)
	action.EphemeralToken, position = util.ParseToken(create, position)
	action.SecretKey, position = util.ParseByteArray(create, position)
	action.PrivateContent, position = util.ParseByteArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
