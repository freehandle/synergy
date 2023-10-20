package actions

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/util"
)

// % + 1 Vote...or 100%
// Supermahority to change policy rule
// Majority for anything else
type Policy struct {
	Majority      int
	SuperMajority int
}

func PutPolicy(policy Policy, bytes *[]byte) {
	*bytes = append(*bytes, byte(policy.Majority), byte(policy.SuperMajority))
}

func ParsePolicy(data []byte, position int) (Policy, int) {
	policy := Policy{
		int(data[position]),
		int(data[position+1]),
	}
	return policy, position + 2
}

type Vote struct {
	Epoch   uint64
	Author  crypto.Token
	Reasons string
	Hash    crypto.Hash
	Approve bool
}

func (c *Vote) Reasoning() string {
	return c.Reasons
}

func (c *Vote) Hashed() crypto.Hash {
	return crypto.Hasher(c.Serialize())
}

func (c *Vote) Authored() crypto.Token {
	return c.Author
}

func (v *Vote) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(v.Epoch, &bytes)
	util.PutToken(v.Author, &bytes)
	util.PutByte(AVote, &bytes)
	util.PutString(v.Reasons, &bytes)
	util.PutHash(v.Hash, &bytes)
	util.PutBool(v.Approve, &bytes)
	return bytes
}

func PutOptionalPolicy(policy *Policy, bytes *[]byte) {
	if policy == nil {
		*bytes = append(*bytes, 0)
		return
	}
	*bytes = append(*bytes, 1, byte(policy.Majority), byte(policy.SuperMajority))
}

func ParseOptionalPolicy(data []byte, position int) (*Policy, int) {
	if data[position] == 0 {
		return nil, position + 1
	}
	if data[position] != 1 || data[position+1] > 100 || data[position+2] > 100 {
		return nil, position + 3
	}
	return &Policy{
		Majority:      int(data[position+1]),
		SuperMajority: int(data[position+2]),
	}, position + 3
}

func ParseVote(vote []byte) *Vote {
	action := Vote{}
	position := 0
	action.Epoch, position = util.ParseUint64(vote, position)
	action.Author, position = util.ParseToken(vote, position)
	if vote[position] != AVote {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(vote, position)
	action.Hash, position = util.ParseHash(vote, position)
	action.Approve, position = util.ParseBool(vote, position)
	if position != len(vote) {
		return nil
	}
	return &action
}
