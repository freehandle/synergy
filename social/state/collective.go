package state

import (
	"errors"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
)

type ConsensusState byte

const (
	Undecided ConsensusState = iota
	Favorable
	Against
)

type Collective struct {
	Name        string
	Members     map[crypto.Token]struct{}
	Description string
	Policy      actions.Policy
}

func (c *Collective) GetPolicy() (majority int, supermajority int) {
	majority = c.Policy.Majority
	supermajority = c.Policy.SuperMajority
	return
}

func (c *Collective) ListOfMembers() map[crypto.Token]struct{} {
	return c.Members
}

func (c *Collective) ListOfTokens() map[crypto.Token]struct{} {
	return c.Members
}

func (c *Collective) CollectiveName() string {
	return c.Name
}

func (c *Collective) Photo() *Collective {
	cloned := Collective{
		Name:    c.Name,
		Members: make(map[crypto.Token]struct{}),
		Policy: actions.Policy{
			Majority:      c.Policy.Majority,
			SuperMajority: c.Policy.SuperMajority,
		},
		Description: c.Description,
	}
	for member, _ := range c.Members {
		cloned.Members[member] = struct{}{}
	}
	return &cloned
}

func (c *Collective) IncludeMember(token crypto.Token) {
	c.Members[token] = struct{}{}
}

func (c *Collective) RemoveMember(token crypto.Token) {
	delete(c.Members, token)
}

func (c *Collective) ChangeMajority(majority int) {
	c.Policy.Majority = majority
}

func (c *Collective) Consensus(hash crypto.Hash, votes []actions.Vote) ConsensusState {
	required := len(c.Members)*c.Policy.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, len(c.Members), hash, votes)
}

func (c *Collective) ConsensusEpoch(votes []actions.Vote) uint64 {
	required := len(c.Members)*c.Policy.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensusEpoch(c.Members, required, votes)
}

func (c *Collective) Unanimous(hash crypto.Hash, votes []actions.Vote) ConsensusState {
	required := len(c.Members)
	return consensus(c.Members, required, len(c.Members), hash, votes)
}

func (c *Collective) SuperConsensus(hash crypto.Hash, votes []actions.Vote) ConsensusState {
	required := len(c.Members)*c.Policy.SuperMajority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, len(c.Members), hash, votes)
}

func (c *Collective) IsMember(token crypto.Token) bool {
	_, ok := c.Members[token]
	return ok
}

type UnamedCollective struct {
	Members  map[crypto.Token]struct{}
	Majority int
}

func (c *UnamedCollective) CollectiveName() string {
	return ""
}

func (c *UnamedCollective) ListOfMembers() map[crypto.Token]struct{} {
	return c.Members
}

func (c *UnamedCollective) ListOfTokens() map[crypto.Token]struct{} {
	return c.Members
}

func (c *UnamedCollective) Consensus(hash crypto.Hash, votes []actions.Vote) ConsensusState {
	required := len(c.Members)*c.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensus(c.Members, required, len(c.Members), hash, votes)
}

func (c *UnamedCollective) ConsensusEpoch(votes []actions.Vote) uint64 {
	required := len(c.Members)*c.Majority/100 + 1
	if required > len(c.Members) {
		required = len(c.Members)
	}
	return consensusEpoch(c.Members, required, votes)
}

func (c *UnamedCollective) Unanimous(hash crypto.Hash, votes []actions.Vote) ConsensusState {
	required := len(c.Members)
	return consensus(c.Members, required, len(c.Members), hash, votes)
}

func (c *UnamedCollective) IsMember(token crypto.Token) bool {
	_, ok := c.Members[token]
	return ok
}

func (c *UnamedCollective) IncludeMember(token crypto.Token) {
	c.Members[token] = struct{}{}
}

func (c *UnamedCollective) RemoveMember(token crypto.Token) {
	delete(c.Members, token)
}

func (c *UnamedCollective) ChangeMajority(majority int) {
	c.Majority = majority
}

func (c *UnamedCollective) GetPolicy() (majority int, supermajority int) {
	majority = c.Majority
	supermajority = c.Majority
	return
}

type PendingUpdate struct {
	Update       *actions.UpdateCollective
	Collective   *Collective
	Hash         crypto.Hash // hash of update action
	ChangePolicy bool
	Votes        []actions.Vote
}

func (p *PendingUpdate) IncorporateVote(vote actions.Vote, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	consensus := Undecided
	if p.ChangePolicy {
		consensus = p.Collective.SuperConsensus(p.Hash, p.Votes)
	} else {
		consensus = p.Collective.Consensus(p.Hash, p.Votes)
	}
	state.index.IndexActionStatus(p.Hash, consensus)
	if consensus == Undecided {
		return nil
	}
	state.IndexConsensus(vote.Hash, consensus)
	state.Proposals.Delete(p.Hash)
	if consensus == Against {
		return nil
	}
	// consensus is favorable, update collective

	// p.Collective is a photo, we need the original to update
	collective, ok := state.Collective(p.Collective.Name)
	if !ok {
		return nil
	}

	if p.Update.Description != nil {
		collective.Description = *p.Update.Description
	}
	if p.ChangePolicy {
		newPolicy := actions.Policy{
			Majority:      p.Collective.Policy.Majority,
			SuperMajority: p.Collective.Policy.SuperMajority,
		}

		if p.Update.Majority != nil {
			newPolicy.Majority = int(*p.Update.Majority)
		}
		if p.Update.SuperMajority != nil {
			newPolicy.SuperMajority = int(*p.Update.SuperMajority)
		}
		collective.Policy = newPolicy
	}
	return nil
}

type PendingRequestMembership struct {
	Request    *actions.RequestMembership
	Collective *Collective
	Hash       crypto.Hash
	Votes      []actions.Vote
}

func (p *PendingRequestMembership) IncorporateVote(vote actions.Vote, state *State) error {
	// if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
	// 	return err
	// }
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	consensus := p.Collective.Consensus(vote.Hash, p.Votes)
	state.index.IndexActionStatus(p.Hash, consensus)
	if consensus == Undecided {
		return nil
	}
	state.IndexConsensus(vote.Hash, consensus)
	state.Proposals.Delete(p.Hash)
	if consensus == Against {
		return nil
	}
	collective, ok := state.Collective(p.Collective.Name)
	if !ok {
		return errors.New("collective not found")
	}
	collective.Members[p.Request.Author] = struct{}{}
	return nil
}

type PendingRemoveMember struct {
	Remove     *actions.RemoveMember
	Collective *Collective
	Hash       crypto.Hash
	Votes      []actions.Vote
}

func (p *PendingRemoveMember) IncorporateVote(vote actions.Vote, state *State) error {
	if err := isValidVote(p.Hash, vote, p.Votes); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	consensus := p.Collective.Consensus(p.Hash, p.Votes)
	state.index.IndexActionStatus(p.Hash, consensus)
	if consensus == Undecided {
		return nil
	}
	state.IndexConsensus(vote.Hash, consensus)
	state.Proposals.Delete(p.Hash)
	if consensus == Against {
		return nil
	}
	collective, ok := state.Collective(p.Collective.Name)
	if !ok {
		return errors.New("collective not found")
	}
	delete(collective.Members, p.Remove.Member)
	return nil
}
