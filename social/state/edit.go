package state

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
)

type Edit struct {
	Authors  Consensual
	Date     uint64
	Reasons  string
	Draft    *Draft
	EditType string
	Edit     crypto.Hash
	Votes    []actions.Vote
}

func (e *Edit) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, e.Votes, e.Edit); err != nil {
		return err
	}
	e.Votes = append(e.Votes, vote)
	consensus := e.Authors.Consensus(e.Edit, e.Votes)
	if consensus == Undecided {
		return nil
	}
	// new consensus
	if consensus == Favorable {
		e.Draft.Edits = append(e.Draft.Edits, e)
	}
	state.Edits[e.Edit] = e
	state.Proposals.Delete(e.Edit)
	state.IndexConsensus(e.Edit, consensus == Favorable)
	// to do where to put edits?
	return nil
}
