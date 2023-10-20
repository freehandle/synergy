package actions

import (
	"reflect"
	"testing"

	"github.com/lienkolabs/breeze/crypto"
)

var (
	reaction = &React{
		Epoch:      25,
		Author:     crypto.Token{},
		Reasons:    "react test",
		OnBehalfOf: "first_collective",
		Hash:       crypto.ZeroValueHash,
		Reaction:   1,
	}
)

func TestReact(t *testing.T) {
	r := ParseReact(reaction.Serialize())
	if r == nil {
		t.Error("Could not parse actions React")
		return
	}
	if !reflect.DeepEqual(r, reaction) {
		t.Error("Parse and Serialize not working for actions React")
	}
}
