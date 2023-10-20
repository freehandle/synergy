package actions

import (
	"reflect"
	"testing"

	"github.com/freehandle/breeze/crypto"
)

var (
	edit = &Edit{
		Epoch:         20,
		Author:        crypto.Token{},
		Reasons:       "edit test",
		OnBehalfOf:    "first_collective",
		CoAuthors:     []crypto.Token{},
		EditedDraft:   crypto.Hash{},
		ContentType:   "txt",
		ContentHash:   crypto.Hash{},
		NumberOfParts: 1,
		Content:       content,
	}
)

func TestEdit(t *testing.T) {
	e := ParseEdit(edit.Serialize())
	if e == nil {
		t.Error("Could not parse actions Edit")
		return
	}
	if !reflect.DeepEqual(e, edit) {
		t.Error("Parse and Serialize not working for actions Edit")
	}
}
