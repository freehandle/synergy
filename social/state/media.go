package state

import (
	"errors"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

type PendingMedia struct {
	Hash          crypto.Hash
	NumberOfParts byte
	Parts         []*actions.MultipartMedia
}

func (p *PendingMedia) Append(m *actions.MultipartMedia) ([]byte, error) {
	if m.Of != p.NumberOfParts || m.Part > m.Of-1 {
		return nil, errors.New("incompatible number of parts")
	}
	p.Parts[m.Part] = m
	size := 0
	for _, part := range p.Parts {
		if part == nil {
			return nil, nil
		}
		size += len(m.Data)
	}
	concanate := make([]byte, 0, size)
	for _, part := range p.Parts {
		concanate = append(concanate, part.Data...)
	}
	if crypto.Hasher(concanate) != p.Hash {
		return nil, errors.New("incompatible hash")
	}
	return concanate, nil
}
