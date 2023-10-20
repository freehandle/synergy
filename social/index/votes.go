package index

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/state"
)

type SetOfHashes struct {
	set map[crypto.Hash]struct{}
}

func NewSetOfHashes() *SetOfHashes {
	return &SetOfHashes{
		set: make(map[crypto.Hash]struct{}),
	}
}

func (s *SetOfHashes) Add(hash crypto.Hash) {
	s.set[hash] = struct{}{}
}

func (s *SetOfHashes) Remove(hash crypto.Hash) {
	delete(s.set, hash)
}

func (s *SetOfHashes) All() map[crypto.Hash]struct{} {
	all := make(map[crypto.Hash]struct{})
	if s.set == nil || len(s.set) == 0 {
		return all
	}
	for hash := range s.set {
		all[hash] = struct{}{}
	}
	return all
}

func (i *Index) IndexVoteHash(c state.Consensual, hash crypto.Hash) {
	//p.mu.Lock()
	//defer p.mu.Unlock()
	members := c.ListOfTokens()
	for token := range members {
		if _, ok := i.indexedMembers[token]; ok {
			if tokenIndex, ok := i.indexVotes[token]; ok {
				tokenIndex.Add(hash)
			} else {
				newSet := NewSetOfHashes()
				newSet.Add(hash)
				i.indexVotes[token] = newSet
			}
		}
	}
}

func (i *Index) RemoveVoteHash(hash crypto.Hash) {
	i.indexCompletedVotes[hash] = i.stateProposals.Votes(hash)
	for _, hashes := range i.indexVotes {
		hashes.Remove(hash)
	}
}

func casted(votes []actions.Vote, token crypto.Token) bool {
	for _, vote := range votes {
		if vote.Author.Equal(token) {
			return true
		}
	}
	return false
}

func (i *Index) GetVotes(token crypto.Token) map[crypto.Hash]struct{} {
	hashes := i.indexVotes[token]
	if hashes == nil {
		return nil
	}
	open := make(map[crypto.Hash]struct{})
	for hash := range hashes.All() {
		if !casted(i.stateProposals.Votes(hash), token) {
			open[hash] = struct{}{}
		}
	}
	return open
}
