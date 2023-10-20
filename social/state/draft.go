package state

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
)

type Draft struct {
	Title           string
	Date            uint64
	Description     string
	Authors         Consensual
	DraftType       string
	DraftHash       crypto.Hash // this must be a valid Media in the state
	PreviousVersion *Draft
	Keywords        []string
	References      []crypto.Hash
	Votes           []actions.Vote
	Pinned          []*Board
	Edits           []*Edit
	Aproved         bool
}

// RULES:
// ======
// If there is no previous version:
// If on behalf of collective, collective must consent
// If with co-authors every co-author must consent
//
// If there is a previous version
// If new version is with co-authors, every new co-author must consent
// Previous authors must collectively consent according to policy
// Current authors must collectively consent according to policy
func (d *Draft) Consensus() ConsensusState {
	if d.Aproved {
		return Favorable
	}

	currentAuthorsConsensus := d.Authors.Consensus(d.DraftHash, d.Votes)
	if currentAuthorsConsensus != Favorable {
		return currentAuthorsConsensus
	}
	previous := d.PreviousVersion
	if previous == nil {
		if d.Authors.CollectiveName() == "" {
			// every co-author must vote
			return d.Authors.Unanimous(d.DraftHash, d.Votes)
		}
		// it is a collective with consensus formed
		return Favorable
	}
	previousAuthorsConsensus := previous.Authors.Consensus(d.DraftHash, d.Votes)
	if previousAuthorsConsensus != Favorable {
		return previousAuthorsConsensus
	}
	collective := d.Authors.CollectiveName()
	if collective != "" {
		// current version is collective and previous version consent, its ok
		return Favorable
	}
	previousCollective := previous.Authors.CollectiveName()
	if previousCollective != "" {
		// if previous is a collective and current is not... every co-author must
		// sign
		return d.Authors.Unanimous(d.DraftHash, d.Votes)
	}
	// here we know that neither current nor previous version is a collective
	currentMembers := d.Authors.ListOfMembers()
	previousMembers := previous.Authors.ListOfMembers()
	newMembers := make(map[crypto.Token]struct{})
	for token, _ := range currentMembers {
		if _, ok := previousMembers[token]; !ok {
			newMembers[token] = struct{}{}
		}
	}
	for _, vote := range d.Votes {
		if _, ok := newMembers[vote.Author]; ok && vote.Hash == d.DraftHash {
			delete(newMembers, vote.Author)
		}
	}
	if len(newMembers) == 0 {
		return Favorable
	}
	return Undecided
}

// IncorpoateVote checks if vote scope is correct (hash) if vote was not alrerady
// cast. If new valid vote returns if the new vote is sufficient for consensus
func (d *Draft) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, d.Votes, d.DraftHash); err != nil {
		return err
	}
	d.Votes = append(d.Votes, vote)
	if d.Aproved {
		return nil
	}
	consensus := d.Consensus()
	if consensus == Undecided {
		return nil
	}
	if consensus == Favorable {
		d.Aproved = true
		state.Drafts[d.DraftHash] = d
	}
	state.Proposals.Delete(d.DraftHash)
	state.IndexConsensus(d.DraftHash, consensus == Favorable)
	return nil
}
