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
	Approved bool
}

// Se escrito em coautoria precisa necessariamente de aprovacao de todos

func (e *Edit) Consensus() ConsensusState {
	if e.Approved {
		return Favorable
	}
	currentAuthorsConsensus := e.Authors.Consensus(e.Edit, e.Votes)
	// testa se é coletivo (consenso undecided) se for manda a funcao de consenso do coletivo
	if currentAuthorsConsensus != Favorable {
		return currentAuthorsConsensus
	}
	// testa se é um coletivo sem nome (coautoria), se for manda funcao de unanimidade
	if e.Authors.CollectiveName() == "" {
		return e.Authors.Unanimous(e.Edit, e.Votes)
	}
	// caso contrario é aprovado, autoria unica
	return Favorable
}

func (e *Edit) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, e.Votes, e.Edit); err != nil {
		return err
	}
	e.Votes = append(e.Votes, vote)
	if e.Approved {
		return nil
	}
	// consensus := e.Authors.Consensus(e.Edit, e.Votes)
	consensus := e.Consensus()
	state.index.IndexActionStatus(e.Edit, consensus)
	if consensus == Undecided {
		return nil
	}
	// new consensus
	if consensus == Favorable {
		e.Approved = true
		e.Draft.Edits = append(e.Draft.Edits, e)
		state.Edits[e.Edit] = e
	}
	state.Proposals.Delete(e.Edit)
	state.IndexConsensus(e.Edit, consensus)
	// to do where to put edits?
	return nil
}
