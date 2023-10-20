package actions

import (
	"reflect"
	"testing"

	"github.com/lienkolabs/breeze/crypto"
)

var (
	content = make([]byte, 10)

	draft = &Draft{
		Epoch:         18,
		Author:        crypto.Token{},
		Reasons:       "draft test",
		OnBehalfOf:    "first_collective",
		CoAuthors:     []crypto.Token{},
		Policy:        policy,
		Title:         "first_draft",
		Keywords:      exampleArray[:],
		Description:   "draft test",
		ContentType:   "txt",
		ContentHash:   crypto.Hash{},
		NumberOfParts: 1,
		Content:       content,
		PreviousDraft: crypto.Hash{},
		References:    []crypto.Hash{},
	}

	release = &ReleaseDraft{
		Epoch:       19,
		Author:      crypto.Token{},
		Reasons:     "release draft test",
		ContentHash: crypto.Hash{},
	}
)

func TestDraft(t *testing.T) {
	d := ParseDraft(draft.Serialize())
	if d == nil {
		t.Error("Could not parse actions Draft")
		return
	}
	if !reflect.DeepEqual(d, draft) {
		t.Error("Parse and Serialize not working for actions Draft")
	}
}

func TestReleaseDraft(t *testing.T) {
	r := ParseReleaseDraft(release.Serialize())
	if r == nil {
		t.Error("Could nor parse actions ReleaseDraft")
		return
	}
	if !reflect.DeepEqual(r, release) {
		t.Error("Parse and Serialize not working for actions ReleaseDraft")
	}
}
