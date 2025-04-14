package state

import (
	"errors"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
)

// Trata do objeto board (que vai existir no state) e todas as suas fases
// e funcionalidades

// objeto board
type Board struct {
	Name        string
	Keyword     []string
	Description string
	Collective  *Collective
	Editors     *UnamedCollective
	Pinned      []*Draft
	Hash        crypto.Hash
}

func (b *Board) Pin(d *Draft) error {
	for _, pinned := range b.Pinned {
		if pinned == d {
			return errors.New("already pinned")
		}
	}
	b.Pinned = append(b.Pinned, d)
	return nil
}

func (b *Board) Remove(d *Draft) error {
	for n, pinned := range b.Pinned {
		if pinned == d {
			b.Pinned = append(b.Pinned[0:n], b.Pinned[n+1:]...)
			return nil
		}
	}
	return errors.New("not pinned")
}

func (b *Board) Last(n int) []*Draft {
	if len(b.Pinned) <= n {
		return b.Pinned
	}
	return b.Pinned[len(b.Pinned)-n:]
}

func (b *Board) First(n int) []*Draft {
	if len(b.Pinned) <= n {
		return b.Pinned
	}
	return b.Pinned[0:n]
}

type PendingUpdateBoard struct {
	Origin      *actions.UpdateBoard
	Keywords    *[]string
	PinMajority *byte
	Description *string
	Board       *Board
	Hash        crypto.Hash
	Votes       []actions.Vote
}

func (b *PendingUpdateBoard) IncorporateVote(vote actions.Vote, state *State) error {
	IsNewValidVote(vote, b.Votes, b.Hash)
	b.Votes = append(b.Votes, vote)
	consensus := b.Board.Collective.Consensus(vote.Hash, b.Votes)
	state.index.IndexActionStatus(b.Hash, consensus)
	if consensus == Undecided {
		return nil
	}
	state.IndexConsensus(vote.Hash, consensus)
	state.Proposals.Delete(b.Hash)
	if consensus == Against {
		return nil
	}
	if b.PinMajority != nil {
		b.Board.Editors.ChangeMajority(int(*b.PinMajority))
	}
	if b.Description != nil {
		b.Board.Description = *b.Description
	}
	if b.Keywords != nil {
		b.Board.Keyword = *b.Keywords
	}
	return nil
}

type PendingBoard struct {
	Origin *actions.CreateBoard
	Board  *Board
	Hash   crypto.Hash
	Votes  []actions.Vote
}

func (b *PendingBoard) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, b.Votes, b.Hash); err != nil {
	}
	b.Votes = append(b.Votes, vote)
	consensus := b.Board.Collective.Consensus(vote.Hash, b.Votes)
	state.index.IndexActionStatus(b.Hash, consensus)
	if consensus == Undecided {
		return nil
	}
	state.IndexConsensus(vote.Hash, consensus)
	state.Proposals.Delete(b.Hash)
	if consensus == Against {
		return nil
	}
	if state.index != nil {
		state.index.AddBoardToCollective(b.Board, b.Board.Collective)
	}
	hash := crypto.Hasher([]byte(b.Board.Name))
	if _, ok := state.Boards[hash]; ok {
		return errors.New("board already exists")
	}
	state.Boards[hash] = b.Board
	return nil
}

type Pin struct {
	Hash  crypto.Hash // hash da [epoch, board, draft, pin/unpin]
	Epoch uint64
	Board *Board
	Draft *Draft
	Pin   bool
	Votes []actions.Vote
}

func (p *Pin) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	consensus := p.Board.Editors.Consensus(vote.Hash, p.Votes)
	state.index.IndexActionStatus(p.Hash, consensus)
	if consensus == Undecided {
		return nil
	}
	state.IndexConsensus(vote.Hash, consensus)
	state.Proposals.Delete(p.Hash)
	if consensus == Against {
		return nil
	}
	if p.Pin {
		// coloca o pin no draft
		p.Draft.Pinned = append(p.Draft.Pinned, p.Board)
		err := p.Board.Pin(p.Draft)
		return err
	}
	// aqui eh um unpin
	if len(p.Draft.Pinned) > 0 {
		for n, pin := range p.Draft.Pinned {
			// se estiver na lista
			if pin == p.Board {
				// tira da lista
				p.Draft.Pinned = append(p.Draft.Pinned[:n], p.Draft.Pinned[n+1:]...)
				break
			}
		}
	}
	return p.Board.Remove(p.Draft)
}

type BoardEditor struct {
	Hash   crypto.Hash // hash of combined draft hash + board name + epoch + pin/remove?
	Epoch  uint64
	Board  *Board
	Editor crypto.Token
	Insert bool
	Votes  []actions.Vote
}

func (e *BoardEditor) IncorporateVote(vote actions.Vote, state *State) error {
	IsNewValidVote(vote, e.Votes, e.Hash)
	e.Votes = append(e.Votes, vote)
	consensus := e.Board.Collective.Consensus(vote.Hash, e.Votes)
	state.index.IndexActionStatus(e.Hash, consensus)
	if consensus == Undecided {
		return nil
	}
	state.IndexConsensus(vote.Hash, consensus)
	state.Proposals.Delete(vote.Hash)
	if consensus == Against {
		return nil
	}
	if e.Insert {
		e.Board.Editors.IncludeMember(e.Editor)
	} else {
		e.Board.Editors.RemoveMember(e.Editor)
	}
	return nil
}
