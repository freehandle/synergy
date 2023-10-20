package actions

import (
	"reflect"
	"testing"

	"github.com/lienkolabs/breeze/crypto"
)

var (
	stamp = &ImprintStamp{
		Epoch:      27,
		Author:     crypto.Token{},
		Reasons:    "imprint stamp test",
		OnBehalfOf: "first_collective",
		Hash:       crypto.Hash{},
	}
)

func TestImprintStamp(t *testing.T) {
	s := ParseImprintStamp(stamp.Serialize())
	if s == nil {
		t.Error("Could not parse actions ImprintStamp")
		return
	}
	if !reflect.DeepEqual(s, stamp) {
		t.Error("Parse and Serialize not working for actions ImprintStamp")
	}
}
