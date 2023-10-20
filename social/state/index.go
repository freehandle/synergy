package state

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Indexer interface {
	AddBoardToCollective(*Board, *Collective)
	RemoveBoardFromCollective(*Board, *Collective)
	AddStampToCollective(*Stamp, *Collective)
	AddEventToCollective(*Event, *Collective)
	RemoveEventFromCollective(*Event, *Collective)
	IndexConsensus(crypto.Hash, bool)
	IndexAction(action actions.Action)
	IndexVoteHash(Consensual, crypto.Hash)
	RemoveVoteHash(crypto.Hash)
	AddDraftToIndex(*Draft)
	AddEditToIndex(*Edit)
	AddCheckin(crypto.Token, *Event)
}
