package actions

import (
	"reflect"
	"testing"

	"github.com/freehandle/breeze/crypto"
)

var (
	policy = &Policy{
		Majority:      10,
		SuperMajority: 20,
	}

	collective = &CreateCollective{
		Epoch:       14,
		Author:      crypto.Token{},
		Reasons:     "create collective test",
		Name:        "first_collective",
		Description: "create collective test",
		Policy:      *policy,
	}

	uCollective = &UpdateCollective{
		Epoch:         15,
		Author:        crypto.Token{},
		OnBehalfOf:    "first_collective",
		Description:   nil,
		Majority:      nil,
		SuperMajority: nil,
	}

	request = &RequestMembership{
		Epoch:      16,
		Author:     crypto.Token{},
		Reasons:    "request membership test",
		Collective: "first_collective",
		Include:    true,
	}

	remove = &RemoveMember{
		Epoch:      17,
		Author:     crypto.Token{},
		OnBehalfOf: "first_collective",
		Reasons:    "remove member test",
		Member:     crypto.Token{},
	}
)

func TestCreateCollective(t *testing.T) {
	c := ParseCreateCollective(collective.Serialize())
	if c == nil {
		t.Error("Could not parse actions CreateCollective")
		return
	}
	if !reflect.DeepEqual(c, collective) {
		t.Error("Parse and Serialize not working for actions CreateCollective")
	}
}

func TestUpdateCollective(t *testing.T) {
	u := ParseUpdateCollective(uCollective.Serialize())
	if u == nil {
		t.Error("Could not parse actions UpdateCollective")
		return
	}
	if !reflect.DeepEqual(u, uCollective) {
		t.Error("Parse and Serialize not working for actions UpdateCollective")
	}
}

func TestRequestMembership(t *testing.T) {
	r := ParseRequestMembership(request.Serialize())
	if r == nil {
		t.Error("Could not parse actions RequestMembership")
		return
	}
	if !reflect.DeepEqual(r, request) {
		t.Error("Parse and Serialize not working for actions RequestMemebership")
	}
}

func TestRemoveMember(t *testing.T) {
	r := ParseRemoveMember(remove.Serialize())
	if r == nil {
		t.Error("Could not parse actions Remove Member")
		return
	}
	if !reflect.DeepEqual(r, remove) {
		t.Error("Parse and Serialize not working for actions RemoveMember")
	}
}
