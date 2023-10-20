package state

import "github.com/lienkolabs/breeze/crypto"

/*
Notify implements an interface to send notification messages through a
channel about modifications on objects within the state of the synergy
protocol. This is useful for example, for the implementation of user
interfaces related to the protocol.
*/

// Action is action that is trigering the notification. For example a
// instruction to pin a draft on a board will notify about a PinAction
type Action byte

const (
	ReactAction Action = iota
	PinAction
	PublishAction
	DraftAction
	EditAction
	BoardAction
	JournalAction
	CollectiveAction
	MediaAction
	VoteAction
	ExpireProposal
	AcceptProposal
	SigninAction
	MediaUpload
)

type Object byte

const (
	NoObject Object = iota
	AuthorObject
	DraftObject
	EditObject
	BoardObject
	JournalObject
	EventObject
	CollectiveObject
	MemberObject
	MediaObject
)

type Updated struct {
	Action Action
	Object Object
	Hash   crypto.Hash
}

type Notifier chan Updated

func (n Notifier) Notify(origin Action, affects Object, id crypto.Hash) {
	n <- Updated{Action: origin, Object: affects, Hash: id}
}
