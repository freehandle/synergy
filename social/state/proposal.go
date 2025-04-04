package state

import (
	"errors"
	"sync"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
)

// contagem como eh feito no caso das acoes, lista das proposals
const (
	UpdateCollectiveProposal byte = iota
	RequestMembershipProposal
	RemoveMemberProposal
	DraftProposal
	EditProposal
	CreateBoardProposal
	UpdateBoardProposal
	PinProposal
	BoardEditorProposal
	ReleaseDraftProposal
	ImprintStampProposal
	ReactProposal
	CreateEventProposal
	CancelEventProposal
	UpdateEventProposal
	EventCheckinGreetProposal
	UnkownProposal
)

var proposalNames = []string{
	"Update Collective",
	"Request Membership",
	"Remove Member",
	"Draft",
	"Edit",
	"Create Board",
	"Update Board",
	"Pin",
	"Board Editor",
	"Release Draft",
	"Imprint Stamp",
	"React",
	"Create Event",
	"Cancel Event",
	"Update Event",
	"Unkown",
}

var ErrProposalNotFound = errors.New("proposal not found")

type Proposal interface {
	IncorporateVote(vote actions.Vote, state *State) error
}

func NewProposals(i Indexer) *Proposals {
	return &Proposals{
		mu:         &sync.Mutex{},
		all:        make(map[crypto.Hash]byte),
		stateIndex: i,
		//index:             make(map[crypto.Token]*SetOfHashes),
		UpdateCollective:  make(map[crypto.Hash]*PendingUpdate),
		RequestMembership: make(map[crypto.Hash]*PendingRequestMembership),
		RemoveMember:      make(map[crypto.Hash]*PendingRemoveMember),
		Draft:             make(map[crypto.Hash]*Draft),
		Edit:              make(map[crypto.Hash]*Edit),
		CreateBoard:       make(map[crypto.Hash]*PendingBoard),
		UpdateBoard:       make(map[crypto.Hash]*PendingUpdateBoard),
		Pin:               make(map[crypto.Hash]*Pin),
		BoardEditor:       make(map[crypto.Hash]*BoardEditor),
		ReleaseDraft:      make(map[crypto.Hash]*Release),
		ImprintStamp:      make(map[crypto.Hash]*Stamp),
		//react map[crypto.Hash]*
		CreateEvent:  make(map[crypto.Hash]*Event),
		CancelEvent:  make(map[crypto.Hash]*CancelEvent),
		UpdateEvent:  make(map[crypto.Hash]*EventUpdate),
		GreetCheckin: make(map[crypto.Hash]*EventCheckinGreet),
	}
}

/*type SetOfHashes struct {
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
*/

type Proposals struct {
	mu         *sync.Mutex
	all        map[crypto.Hash]byte // hash do que ta pendente pro tipo de proposal
	stateIndex Indexer
	//index             map[crypto.Token]*SetOfHashes // token do membro pra um conjunto de hashs dos votos que ele precisa dar
	UpdateCollective  map[crypto.Hash]*PendingUpdate
	RequestMembership map[crypto.Hash]*PendingRequestMembership
	RemoveMember      map[crypto.Hash]*PendingRemoveMember
	Draft             map[crypto.Hash]*Draft
	Edit              map[crypto.Hash]*Edit
	CreateBoard       map[crypto.Hash]*PendingBoard
	UpdateBoard       map[crypto.Hash]*PendingUpdateBoard
	Pin               map[crypto.Hash]*Pin
	BoardEditor       map[crypto.Hash]*BoardEditor
	ReleaseDraft      map[crypto.Hash]*Release
	ImprintStamp      map[crypto.Hash]*Stamp
	//react map[crypto.Hash]*
	CreateEvent  map[crypto.Hash]*Event
	CancelEvent  map[crypto.Hash]*CancelEvent
	UpdateEvent  map[crypto.Hash]*EventUpdate
	GreetCheckin map[crypto.Hash]*EventCheckinGreet
}

func (p *Proposals) GetEvent(hash crypto.Hash) *Event {
	e, _ := p.CreateEvent[hash]
	return e
}

func (p *Proposals) Delete(hash crypto.Hash) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stateIndex != nil {
		p.stateIndex.RemoveVoteHash(hash)
	}
	delete(p.all, hash)
	delete(p.UpdateCollective, hash)
	delete(p.RequestMembership, hash)
	delete(p.RemoveMember, hash)
	delete(p.Draft, hash)
	delete(p.Edit, hash)
	delete(p.CreateBoard, hash)
	delete(p.UpdateBoard, hash)
	delete(p.Pin, hash)
	delete(p.BoardEditor, hash)
	delete(p.ReleaseDraft, hash)
	delete(p.ImprintStamp, hash)
	//react map[crypto.Hash]*
	delete(p.CreateEvent, hash)
	delete(p.CancelEvent, hash)
	delete(p.UpdateEvent, hash)
	/*for _, hashes := range p.index {
		hashes.Remove(hash)
	}*/
	delete(p.GreetCheckin, hash)
}

func (p *Proposals) Kind(hash crypto.Hash) byte {
	kind, ok := p.all[hash]
	if !ok {
		return UnkownProposal
	}
	return kind
}

func (p *Proposals) KindText(hash crypto.Hash) string {
	return proposalNames[p.Kind(hash)]
}

// colocando os hashs pendentes na lista de cada token que precisa votar
func (p *Proposals) indexHash(c Consensual, hash crypto.Hash) {
	if p.stateIndex != nil {
		p.stateIndex.IndexVoteHash(c, hash)
	}
	/*
		p.mu.Lock()
		defer p.mu.Unlock()
		members := c.ListOfTokens()
		for token := range members {
			if _, ok := p.index[token]; !ok {
				p.index[token] = NewSetOfHashes()
			}
			p.index[token].Add(hash)
		}
	*/
}

func casted(votes []actions.Vote, token crypto.Token) bool {
	for _, vote := range votes {
		if vote.Author.Equal(token) {
			return true
		}
	}
	return false
}

/*func (p *Proposals) GetVotes(token crypto.Token) map[crypto.Hash]struct{} {
	hashes := p.index[token]
	if hashes == nil {
		return nil
	}
	open := make(map[crypto.Hash]struct{})
	for hash := range hashes.All() {
		if !casted(p.Votes(hash), token) {
			open[hash] = struct{}{}
		}
	}
	return open
}
*/

func (p *Proposals) AddUpdateCollective(update *PendingUpdate, reason actions.Action) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = UpdateCollectiveProposal
	p.UpdateCollective[update.Hash] = update
}

func (p *Proposals) AddRequestMembership(update *PendingRequestMembership, reason actions.Action) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = RequestMembershipProposal
	p.RequestMembership[update.Hash] = update
}

func (p *Proposals) AddPendingRemoveMember(update *PendingRemoveMember, reason actions.Action) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = RemoveMemberProposal
	p.RemoveMember[update.Hash] = update
}

func (p *Proposals) AddDraft(update *Draft, reason actions.Action) {
	p.indexHash(update.Authors, update.DraftHash)
	if update.PreviousVersion != nil {
		p.indexHash(update.PreviousVersion.Authors, update.DraftHash)
	}
	p.all[update.DraftHash] = DraftProposal
	p.Draft[update.DraftHash] = update
}

func (p *Proposals) AddEdit(update *Edit, reason actions.Action) {
	p.indexHash(update.Draft.Authors, update.Edit)
	p.indexHash(update.Authors, update.Edit)
	p.all[update.Edit] = EditProposal
	p.Edit[update.Edit] = update
}

func (p *Proposals) AddPendingBoard(update *PendingBoard, reason actions.Action) {
	p.indexHash(update.Board.Collective, update.Hash)
	p.all[update.Hash] = CreateBoardProposal
	p.CreateBoard[update.Hash] = update
}

func (p *Proposals) AddPendingUpdateBoard(update *PendingUpdateBoard, reason actions.Action) {
	p.indexHash(update.Board.Editors, update.Hash)
	p.all[update.Hash] = UpdateBoardProposal
	p.UpdateBoard[update.Hash] = update
}

// adicionando aos proposals o que chegou pra ser votado
func (p *Proposals) AddPin(update *Pin, reason actions.Action) {
	// quem vai receber o pedido de voto
	p.indexHash(update.Board.Editors, update.Hash)
	p.all[update.Hash] = PinProposal // adiciona a
	p.Pin[update.Hash] = update
}

func (p *Proposals) AddBoardEditor(update *BoardEditor, reason actions.Action) {
	p.indexHash(update.Board.Collective, update.Hash)
	p.all[update.Hash] = BoardEditorProposal
	p.BoardEditor[update.Hash] = update
}

func (p *Proposals) AddRelease(update *Release, reason actions.Action) {
	p.indexHash(update.Draft.Authors, update.Hash)
	p.all[update.Hash] = ReleaseDraftProposal
	p.ReleaseDraft[update.Hash] = update
}

func (p *Proposals) AddStamp(update *Stamp, reason actions.Action) {
	p.indexHash(update.Reputation, update.Hash) // reputation aqui é = um membro ou coletivo que vai dar o stamp ??
	p.all[update.Hash] = ImprintStampProposal
	p.ImprintStamp[update.Hash] = update
}

func (p *Proposals) AddEvent(update *Event, reason actions.Action) {
	p.indexHash(update.Collective, update.Hash)
	p.all[update.Hash] = CreateEventProposal
	p.CreateEvent[update.Hash] = update
}

func (p *Proposals) AddCancelEvent(update *CancelEvent, reason actions.Action) {
	p.indexHash(update.Event.Collective, update.Hash)
	p.all[update.Hash] = CancelEventProposal
	p.CancelEvent[update.Hash] = update
}

func (p *Proposals) AddEventUpdate(update *EventUpdate, reason actions.Action) {
	p.indexHash(update.Event.Managers, update.Hash)
	p.all[update.Hash] = UpdateEventProposal
	p.UpdateEvent[update.Hash] = update
}

func (p *Proposals) AddEventCheckinGreet(update *EventCheckinGreet, reason actions.Action) {
	// p.indexHash(update.Event.Greets, update.Hash)
	p.all[update.Hash] = EventCheckinGreetProposal
	p.GreetCheckin[update.Hash] = update
}

func (p *Proposals) Has(hash crypto.Hash) bool {
	_, ok := p.all[hash]
	return ok
}

func (p *Proposals) IncorporateVote(vote actions.Vote, state *State) error {
	hash := vote.Hash
	var proposal Proposal
	// qual tipo de proposta ta associado ao hash
	kind, ok := p.all[hash]
	if !ok {
		return ErrProposalNotFound
	}
	switch kind {
	case UpdateCollectiveProposal:
		proposal = p.UpdateCollective[hash]
	case RequestMembershipProposal:
		proposal = p.RequestMembership[hash]
	case RemoveMemberProposal:
		proposal = p.RemoveMember[hash]
	case DraftProposal:
		proposal = p.Draft[hash]
	case EditProposal:
		proposal = p.Edit[hash]
	case CreateBoardProposal:
		proposal = p.CreateBoard[hash]
	case UpdateBoardProposal:
		proposal = p.UpdateBoard[hash]
	case PinProposal:
		proposal = p.Pin[hash]
	case BoardEditorProposal:
		proposal = p.BoardEditor[hash]
	case ReleaseDraftProposal:
		proposal = p.ReleaseDraft[hash]
	case ImprintStampProposal:
		proposal = p.ImprintStamp[hash]
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal = p.CreateEvent[hash]
	case CancelEventProposal:
		proposal = p.CancelEvent[hash]
	case UpdateEventProposal:
		proposal = p.UpdateEvent[hash]
	}
	if proposal == nil {
		return ErrProposalNotFound
	}
	return proposal.IncorporateVote(vote, state)
}

type Pool struct {
	Voters   map[crypto.Token]struct{}
	Majority int
	Votes    []actions.Vote
}

func DeepCopyMembers(m map[crypto.Token]struct{}) map[crypto.Token]struct{} {
	copiedmap := make(map[crypto.Token]struct{})
	for keyvoter, valuevoter := range m {
		copiedmap[keyvoter] = valuevoter
	}
	return copiedmap
}

func (p *Proposals) Pooling(hash crypto.Hash) *Pool {
	kind, ok := p.all[hash]
	if !ok {
		return nil
	}
	switch kind {
	case RequestMembershipProposal:
		proposal := p.RequestMembership[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Collective.ListOfMembers()),
			Majority: proposal.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case UpdateCollectiveProposal:
		proposal := p.UpdateCollective[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Collective.ListOfMembers()),
			Majority: proposal.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case RemoveMemberProposal:
		proposal := p.RemoveMember[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Collective.ListOfMembers()),
			Majority: proposal.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case DraftProposal:
		proposal := p.Draft[hash]
		majority, _ := proposal.Authors.GetPolicy()
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Authors.ListOfMembers()),
			Majority: majority,
			Votes:    proposal.Votes,
		}
	case EditProposal:
		proposal := p.Edit[hash]
		majority, _ := proposal.Authors.GetPolicy()
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Authors.ListOfMembers()),
			Majority: majority,
			Votes:    proposal.Votes,
		}
	case CreateBoardProposal:
		proposal := p.CreateBoard[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Board.Collective.ListOfMembers()),
			Majority: proposal.Board.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case UpdateBoardProposal:
		proposal := p.UpdateBoard[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Board.Collective.ListOfMembers()),
			Majority: proposal.Board.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case PinProposal:
		proposal := p.Pin[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Board.Editors.ListOfMembers()),
			Majority: proposal.Board.Editors.Majority,
			Votes:    proposal.Votes,
		}
	case BoardEditorProposal:
		proposal := p.BoardEditor[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Board.Collective.ListOfMembers()),
			Majority: proposal.Board.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case ReleaseDraftProposal:
		proposal := p.ReleaseDraft[hash]
		majority, _ := proposal.Draft.Authors.GetPolicy()
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Draft.Authors.ListOfMembers()),
			Majority: majority,
			Votes:    proposal.Votes,
		}
	case ImprintStampProposal:
		proposal := p.ImprintStamp[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Reputation.ListOfMembers()),
			Majority: proposal.Reputation.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal := p.CreateEvent[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Collective.ListOfMembers()),
			Majority: proposal.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case CancelEventProposal:
		proposal := p.CancelEvent[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Event.Collective.ListOfMembers()),
			Majority: proposal.Event.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	case UpdateEventProposal:
		proposal := p.UpdateEvent[hash]
		return &Pool{
			Voters:   DeepCopyMembers(proposal.Event.Collective.ListOfMembers()),
			Majority: proposal.Event.Collective.Policy.Majority,
			Votes:    proposal.Votes,
		}
	}
	return nil
}

func (p *Proposals) Votes(hash crypto.Hash) []actions.Vote {
	kind, ok := p.all[hash]
	if !ok {
		return nil
	}
	switch kind {
	case RequestMembershipProposal:
		proposal := p.RequestMembership[hash]
		return proposal.Votes
	case UpdateCollectiveProposal:
		proposal := p.UpdateCollective[hash]
		return proposal.Votes
	case RemoveMemberProposal:
		proposal := p.RemoveMember[hash]
		return proposal.Votes
	case DraftProposal:
		proposal := p.Draft[hash]
		return proposal.Votes
	case EditProposal:
		proposal := p.Edit[hash]
		return proposal.Votes
	case CreateBoardProposal:
		proposal := p.CreateBoard[hash]
		return proposal.Votes
	case UpdateBoardProposal:
		proposal := p.UpdateBoard[hash]
		return proposal.Votes
	case PinProposal:
		proposal := p.Pin[hash]
		return proposal.Votes
	case BoardEditorProposal:
		proposal := p.BoardEditor[hash]
		return proposal.Votes
	case ReleaseDraftProposal:
		proposal := p.ReleaseDraft[hash]
		return proposal.Votes
	case ImprintStampProposal:
		proposal := p.ImprintStamp[hash]
		return proposal.Votes
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal := p.CreateEvent[hash]
		return proposal.Votes
	case CancelEventProposal:
		proposal := p.CancelEvent[hash]
		return proposal.Votes
	case UpdateEventProposal:
		proposal := p.UpdateEvent[hash]
		return proposal.Votes
	}
	return nil
}

func (p *Proposals) OnBehalfOf(hash crypto.Hash) string {
	kind, ok := p.all[hash]
	if !ok {
		return ""
	}
	switch kind {
	case UpdateCollectiveProposal:
		proposal := p.UpdateCollective[hash]
		return proposal.Collective.Name
	case RemoveMemberProposal:
		proposal := p.RemoveMember[hash]
		return proposal.Collective.Name
	case DraftProposal:
		proposal := p.Draft[hash]
		return proposal.Authors.CollectiveName()
	case EditProposal:
		proposal := p.Edit[hash]
		return proposal.Authors.CollectiveName()
	case CreateBoardProposal:
		proposal := p.CreateBoard[hash]
		return proposal.Board.Collective.Name
	case UpdateBoardProposal:
		proposal := p.UpdateBoard[hash]
		return proposal.Board.Collective.Name
	case PinProposal:
		proposal := p.Pin[hash]
		return proposal.Board.Name
	case BoardEditorProposal:
		proposal := p.BoardEditor[hash]
		return proposal.Board.Collective.Name
	case ReleaseDraftProposal:
		proposal := p.ReleaseDraft[hash]
		return proposal.Draft.Authors.CollectiveName()
	case ImprintStampProposal:
		proposal := p.ImprintStamp[hash]
		return proposal.Reputation.Name
	case ReactProposal:
		//
	case CreateEventProposal:
		proposal := p.CreateEvent[hash]
		return proposal.Collective.Name
	case CancelEventProposal:
		proposal := p.CancelEvent[hash]
		return proposal.Event.Collective.Name
	case UpdateEventProposal:
		proposal := p.UpdateEvent[hash]
		return proposal.Event.Collective.Name
	}
	return ""
}
