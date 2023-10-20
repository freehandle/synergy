package state

import (
	"errors"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

// verifica se eh um voto novo e se o hash bate
func IsNewValidVote(vote actions.Vote, voted []actions.Vote, hash crypto.Hash) error {
	if vote.Hash != hash {
		return errors.New("invalid hash")
	}
	for _, cast := range voted {
		if cast.Author == vote.Author {
			return errors.New("vote already cast")
		}
	}
	return nil
}
