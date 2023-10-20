package state

import (
	"errors"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type Consensual interface {
	Consensus(hash crypto.Hash, votes []actions.Vote) ConsensusState
	ConsensusEpoch(votes []actions.Vote) uint64
	IsMember(token crypto.Token) bool
	IncludeMember(token crypto.Token)
	RemoveMember(token crypto.Token)
	ChangeMajority(majority int)
	ListOfMembers() map[crypto.Token]struct{}
	ListOfTokens() map[crypto.Token]struct{}
	CollectiveName() string
	GetPolicy() (majority int, supermajority int)
	Unanimous(hash crypto.Hash, votes []actions.Vote) ConsensusState
}

func consensus(members map[crypto.Token]struct{}, votesRequired int, votesCount int, hash crypto.Hash, votes []actions.Vote) ConsensusState {
	favor := 0
	against := 0
	for _, vote := range votes {
		_, isMember := members[vote.Author]
		if isMember && hash == vote.Hash {
			if vote.Approve {
				favor += 1
				if favor >= votesRequired {
					return Favorable
				}
			} else {
				against += 1
				if against > (votesCount - votesRequired) {
					return Against
				}
			}
		}
	}
	return Undecided
}

func consensusEpoch(members map[crypto.Token]struct{}, votesRequired int, votes []actions.Vote) uint64 {
	count := 0
	for _, vote := range votes {
		_, isMember := members[vote.Author]
		if isMember && vote.Approve {
			count += 1
			if count >= votesRequired {
				return vote.Epoch
			}
		}
	}
	return 0
}

func isValidVote(hash crypto.Hash, vote actions.Vote, signatures []actions.Vote) error {
	if vote.Hash != hash {
		return errors.New("invalid hash")
	}
	for _, cast := range signatures {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	return nil
}

func Authors(majority int, tokens ...crypto.Token) Consensual {
	collective := UnamedCollective{
		Members:  make(map[crypto.Token]struct{}),
		Majority: majority,
	}
	for _, token := range tokens {
		collective.Members[token] = struct{}{}
	}
	return &collective
}
