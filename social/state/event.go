package state

import (
	"errors"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
)

type Greeting struct {
	Action       *actions.GreetCheckinEvent
	EphemeralKey crypto.Token
}

type Event struct {
	Collective     *Collective
	StartAt        time.Time
	EstimatedEnd   time.Time
	Description    string
	Venue          string
	Open           bool
	Public         bool
	Hash           crypto.Hash
	Managers       *UnamedCollective // default é qualquer um do coletivo
	Votes          []actions.Vote
	Checkin        map[crypto.Token]*Greeting
	CheckinReasons map[crypto.Token]string
	Live           bool
	EventReasons   string
}

func (p *Event) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if p.Live {
		return nil
	}
	consensus := p.Collective.Consensus(p.Hash, p.Votes)
	if consensus == Undecided {
		return nil
	}
	// new consensus
	state.IndexConsensus(p.Hash, consensus == Favorable)
	state.Proposals.Delete(p.Hash)
	if consensus == Favorable {
		p.Live = true
		if state.index != nil {
			state.index.AddEventToCollective(p, p.Collective)
		}
		if _, ok := state.Events[p.Hash]; !ok {
			state.Events[p.Hash] = p
			return nil
		} else {
			return errors.New("already live")
		}
	}
	return nil
}

type EventUpdate struct {
	Event        *Event
	StartAt      *time.Time
	EstimatedEnd *time.Time
	Description  *string
	Venue        *string
	Open         *bool
	Public       *bool
	Managers     *UnamedCollective
	Votes        []actions.Vote
	Hash         crypto.Hash
	Updated      bool
	Reasons      string
}

func (p *EventUpdate) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if p.Updated {
		return nil
	}
	consensus := p.Event.Managers.Consensus(p.Hash, p.Votes)
	if consensus == Undecided {
		return nil
	}
	// new consensus, update event details
	state.IndexConsensus(vote.Hash, consensus == Favorable)
	state.Proposals.Delete(p.Hash)
	if consensus == Against {
		return nil
	}
	p.Updated = true
	if event := p.Event; event != nil {
		if p.StartAt != nil {
			event.StartAt = *p.StartAt
		}
		if p.EstimatedEnd != nil {
			event.EstimatedEnd = *p.EstimatedEnd
		}
		if p.Description != nil {
			event.Description = *p.Description
		}
		if p.Venue != nil {
			event.Venue = *p.Venue
		}
		if p.Open != nil {
			event.Open = *p.Open
		}
		if p.Public != nil {
			event.Public = *p.Public
		}
		if p.Managers != nil {
			event.Managers = p.Managers
		}
		return nil
	}
	return errors.New("event not found")
}

type CancelEvent struct {
	Event   *Event
	Hash    crypto.Hash
	Votes   []actions.Vote
	Reasons string
}

func (p *CancelEvent) IncorporateVote(vote actions.Vote, state *State) error {
	if err := IsNewValidVote(vote, p.Votes, p.Hash); err != nil {
		return err
	}
	p.Votes = append(p.Votes, vote)
	if !p.Event.Live {
		return nil
	}
	consensus := p.Event.Managers.Consensus(p.Hash, p.Votes)
	if consensus == Undecided {
		return nil
	}
	// new consensus, update event details
	state.IndexConsensus(vote.Hash, consensus == Favorable)
	if consensus == Favorable {
		p.Event.Live = false
		if state.index != nil {
			state.index.RemoveEventFromCollective(p.Event, p.Event.Collective)
		}
	}
	state.Proposals.Delete(p.Hash)
	return nil
}

type EventCheckinGreet struct {
	Event  *Event
	Hash   crypto.Hash
	Greets []actions.GreetCheckinEvent
}

func (p *EventCheckinGreet) IncorporateGreet(greet actions.GreetCheckinEvent, state *State) error {

	if greet.EventHash != p.Hash {
		return errors.New("invalid hash")
	}
	for _, cast := range p.Greets {
		if cast.CheckedIn == greet.CheckedIn {
			return errors.New("checkin already greeted")
		}
	}
	p.Greets = append(p.Greets, greet)

	if !p.Event.Live {
		return nil
	}

	state.Proposals.Delete(p.Hash)
	return nil
}
