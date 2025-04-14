package state

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
)

type Indexer interface {
	AddBoardToCollective(*Board, *Collective)
	RemoveBoardFromCollective(*Board, *Collective)
	AddStampToCollective(*Stamp, *Collective)
	AddEventToCollective(*Event, *Collective)
	RemoveEventFromCollective(*Event, *Collective)
	IndexConsensus(crypto.Hash, ConsensusState)
	IndexActionStatus(crypto.Hash, ConsensusState)
	IndexAction(action actions.Action)
	IndexVoteHash(Consensual, crypto.Hash)
	RemoveVoteHash(crypto.Hash)
	AddDraftToIndex(*Draft)
	AddEditToIndex(*Edit)
	AddCheckin(crypto.Token, *Event)
	AddMemberToIndex(crypto.Token, string)
}
