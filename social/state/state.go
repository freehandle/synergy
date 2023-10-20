package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

const (
	ReactionsCount   = 5
	ProposalDeadline = 30 * 24 * 60 * 60
)

type State struct {
	Epoch        uint64
	MembersIndex map[string]crypto.Token       // mapa do handle to token
	Members      map[crypto.Hash]string        // mapa do hash do token para o handle
	PendingMedia map[crypto.Hash]*PendingMedia // multi-part media file
	Media        map[crypto.Hash][]byte        // quando termina de receber todas as partes vira o media
	Drafts       map[crypto.Hash]*Draft        // o hash do draft eh o hash da media dele
	Edits        map[crypto.Hash]*Edit         // o hash do edit eh o hash da media dele
	Releases     map[crypto.Hash]*Release      // hash do draft para instancia do release
	Events       map[crypto.Hash]*Event        // hash do evento eh hash da acao do evento
	Collectives  map[crypto.Hash]*Collective   // hash do coletivo eh o hash do nome
	Boards       map[crypto.Hash]*Board        // hash do board eh o hash do nome
	Proposals    *Proposals                    // map[crypto.Hash]Proposal // proposals pending vote actions
	Deadline     map[uint64][]crypto.Hash      // map do epoch que morre para o array de hash dos elementos que vao morrer naquele epoch
	Reactions    [ReactionsCount]map[crypto.Hash]uint
	GenesisTime  time.Time
	index        Indexer
	action       Notifier // pra ser usado pra notificacao real time

}

func (s *State) TimeOfEpoch(epoch uint64) time.Time {
	return s.GenesisTime.Add(time.Duration(epoch) * time.Second)
}

func (s *State) IndexConsensus(hash crypto.Hash, approve bool) {
	if s.index == nil {
		return
	}
	s.index.IndexConsensus(hash, approve)
}

// printa o que ta rolando no terminal
func logAction(a actions.Action) {
	if a == nil {
		fmt.Println("nil action")
	}
	var des string
	switch a.(type) {
	case *actions.Vote:
		des = "vote"
	case *actions.CreateCollective:
		des = "Create Collective"
	case *actions.UpdateCollective:
		des = "Update Collective"
	case *actions.RequestMembership:
		des = "Request Membership"
	case *actions.RemoveMember:
		des = "Remove Member"
	case *actions.Draft:
		return
	case *actions.Edit:
		return
	case *actions.MultipartMedia:
		return
	case *actions.CreateBoard:
		des = "Create Board"
	case *actions.UpdateBoard:
		des = "Update Board"
	case *actions.Pin:
		des = "Pin"
	case *actions.BoardEditor:
		des = "Board Editor"
	case *actions.ReleaseDraft:
		des = "Release Draft"
	case *actions.ImprintStamp:
		des = "Imprint Stamp"
	case *actions.React:
		des = "React"
	case *actions.CreateEvent:
		des = "Create Event"
	case *actions.CancelEvent:
		des = "Cancel Event"
	case *actions.UpdateEvent:
		des = "Update Event"
	case *actions.CheckinEvent:
		des = "Checkin Event"
	case *actions.GreetCheckinEvent:
		des = "Greet Checkin Event"
	}
	text, _ := json.Marshal(a)
	fmt.Printf("%v: %v\n\n", des, string(text))
}

// funcao que esta sendo chamada no SelfGateway do genesis
// valida e incorpora a acao
func (s *State) Action(data []byte) error {
	kind := actions.ActionKind(data)
	// verifica qual o tipo de acao ta sendo processado segundo o byte
	switch kind {
	case actions.AVote:
		action := actions.ParseVote(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.Vote(action)
		return err
	case actions.ACreateCollective:
		action := actions.ParseCreateCollective(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.CreateCollective(action)
		return err
	case actions.AUpdateCollective:
		action := actions.ParseUpdateCollective(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.UpdateCollective(action)
		return err
	case actions.ARequestMembership:
		action := actions.ParseRequestMembership(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.RequestMembership(action)
		return err
	case actions.ARemoveMember:
		action := actions.ParseRemoveMember(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.RemoveMember(action)
		return err
	case actions.ADraft:
		action := actions.ParseDraft(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.Draft(action)
		return err

	case actions.AEdit:
		action := actions.ParseEdit(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.Edit(action)
		return err
	case actions.AMultipartMedia:
		action := actions.ParseMultipartMedia(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		// no need to index
		return s.MultipartMedia(action)

	case actions.ACreateBoard:
		action := actions.ParseCreateBoard(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.CreateBoard(action)
		return err

	case actions.AUpdateBoard:
		action := actions.ParseUpdateBoard(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.UpdateBoard(action)
		return err

	case actions.APin:
		action := actions.ParsePin(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.Pin(action)
		return err
	case actions.ABoardEditor:
		action := actions.ParseBoardEditor(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		return s.BoardEditor(action)

	case actions.AReleaseDraft:
		action := actions.ParseReleaseDraft(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.ReleaseDraft(action)
		return err
	case actions.AImprintStamp:
		action := actions.ParseImprintStamp(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.ImprintStamp(action)
		return err
	case actions.AReact:
		action := actions.ParseReact(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.React(action)
		return err
	case actions.ASignIn:
		action := actions.ParseSignIn(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		// should index signin???
		return s.SignIn(action)

	case actions.ACreateEvent:
		action := actions.ParseCreateEvent(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.CreateEvent(action)
		return err

	case actions.ACancelEvent:
		action := actions.ParseCancelEvent(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.CancelEvent(action)
		return err
	case actions.AUpdateEvent:
		action := actions.ParseUpdateEvent(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.UpdateEvent(action)
		return err
	case actions.ACheckinEvent:
		action := actions.ParseCheckinEvent(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.CheckinEvent(action)
		return err
	case actions.AGreetCheckinEvent:
		action := actions.ParseGreetCheckinEvent(data)
		if action == nil {
			return errors.New("cound not parse action")
		}
		logAction(action)
		s.IndexAction(action)
		err := s.GreetCheckinEvent(action)
		return err
	}

	return errors.New("unrecognized action")
}

// cria o estado inicial
func GenesisState(indexer Indexer) *State {
	state := &State{
		Epoch:        0,
		MembersIndex: make(map[string]crypto.Token),
		Members:      make(map[crypto.Hash]string),
		PendingMedia: make(map[crypto.Hash]*PendingMedia),
		Media:        make(map[crypto.Hash][]byte),
		Drafts:       make(map[crypto.Hash]*Draft),
		Edits:        make(map[crypto.Hash]*Edit),
		Releases:     make(map[crypto.Hash]*Release),
		Events:       make(map[crypto.Hash]*Event),
		Collectives:  make(map[crypto.Hash]*Collective),
		Boards:       make(map[crypto.Hash]*Board),
		Proposals:    NewProposals(indexer),
		Deadline:     make(map[uint64][]crypto.Hash),
		index:        indexer,
	}
	for n := 0; n < ReactionsCount; n++ {
		state.Reactions[n] = make(map[crypto.Hash]uint)
	}
	return state
}

func (s *State) IndexAction(action actions.Action) {
	if s.index == nil {
		return
	}
	s.index.IndexAction(action)
}

func (s *State) Collective(name string) (*Collective, bool) {
	hash := crypto.Hasher([]byte(name))
	col, ok := s.Collectives[hash]
	return col, ok
}

func (s *State) Board(name string) (*Board, bool) {
	hash := crypto.Hasher([]byte(name))
	col, ok := s.Boards[hash]
	return col, ok
}

func (s *State) IsMember(token crypto.Token) bool {
	hash := crypto.HashToken(token)
	_, ok := s.Members[hash]
	return ok
}

func (s *State) Notify(origin Action, objHash crypto.Hash) {
	s.action.Notify(origin, s.hashToObjectType(objHash), objHash)
}

func (s *State) NextBlock() {
	if deadline, ok := s.Deadline[s.Epoch]; ok {
		for _, hash := range deadline {
			s.Proposals.Delete(hash)
			s.Notify(ExpireProposal, hash)
		}
	}
}

func (s *State) hashToObjectType(hash crypto.Hash) Object {
	if _, ok := s.Members[hash]; ok {
		return MemberObject
	}
	if _, ok := s.Drafts[hash]; ok {
		return DraftObject
	}
	if _, ok := s.Edits[hash]; ok {
		return EditObject
	}
	if _, ok := s.Media[hash]; ok {
		return MediaObject
	}
	return NoObject
}

func (s *State) setDeadline(epoch uint64, hash crypto.Hash) {
	if epoch <= s.Epoch {
		return
	}
	if deadlines, ok := s.Deadline[epoch]; ok {
		s.Deadline[epoch] = append(deadlines, hash)
	} else {
		s.Deadline[epoch] = []crypto.Hash{hash}
	}
}

func (s *State) GreetCheckinEvent(greet *actions.GreetCheckinEvent) error {
	event, ok := s.Events[greet.EventHash]
	if !ok {
		return errors.New("event not found")
	}
	greeting, ok := event.Checkin[greet.CheckedIn]
	if !ok {
		return errors.New("checkin not found")
	}
	greeting.Action = greet
	return nil
}

func (s *State) ImprintStamp(stamp *actions.ImprintStamp) error {
	if !s.IsMember(stamp.Author) {
		return errors.New("not a member")
	}
	release, ok := s.Releases[stamp.Hash]
	if !ok {
		return errors.New("release not found")
	}
	collective, ok := s.Collective(stamp.OnBehalfOf)
	if !ok {
		return errors.New("collective not found")
	}
	hash := stamp.Hashed()
	vote := actions.Vote{
		Epoch:   stamp.Epoch,
		Author:  stamp.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	newStamp := Stamp{
		Reputation: collective,
		Release:    release,
		Hash:       hash,
		Votes:      []actions.Vote{},
	}
	s.Proposals.AddStamp(&newStamp, stamp)
	return newStamp.IncorporateVote(vote, s)
}

func (s *State) CheckinEvent(checkin *actions.CheckinEvent) error {
	if !s.IsMember(checkin.Author) {
		return errors.New("not an author")
	}
	event, ok := s.Events[checkin.EventHash]
	if !ok {
		return errors.New("event not found")
	}
	if _, ok := event.Checkin[checkin.Author]; ok {
		return errors.New("already checkin")
	}
	event.Checkin[checkin.Author] = &Greeting{Action: nil, EphemeralKey: checkin.EphemeralToken}
	event.CheckinReasons[checkin.Author] = checkin.Reasons
	if s.index != nil {
		s.index.AddCheckin(checkin.Author, event)
	}
	return nil
}

func (s *State) UpdateEvent(update *actions.UpdateEvent) error {
	event, ok := s.Events[update.EventHash]
	if !ok {
		return errors.New("event not found")
	}
	if !event.Managers.IsMember(update.Author) {
		return errors.New("not a manager of the event")
	}
	hash := update.Hashed()
	selfVote := actions.Vote{
		Epoch:   update.Epoch,
		Author:  update.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	pending := EventUpdate{
		Event:        event,
		StartAt:      update.StartAt,
		EstimatedEnd: update.EstimatedEnd,
		Description:  update.Description,
		Venue:        update.Venue,
		Open:         update.Open,
		Public:       update.Public,
		Hash:         hash,
		Votes:        []actions.Vote{},
	}
	// if event.Managers.Consensus(hash, pending.Votes) {
	// 	if update.StartAt != nil {
	// 		event.StartAt = *update.StartAt
	// 	}
	// 	if update.EstimatedEnd != nil {
	// 		event.EstimatedEnd = *update.EstimatedEnd
	// 	}
	// 	if update.Description != nil {
	// 		event.Description = *update.Description
	// 	}
	// 	if update.Venue != nil {
	// 		event.Venue = *update.Venue
	// 	}
	// 	if update.Open != nil {
	// 		event.Open = *update.Open
	// 	}
	// 	if update.Public != nil {
	// 		event.Public = *update.Public
	// 	}
	// 	if update.ManagerMajority != nil {
	// 		event.Managers.Majority = int(*update.ManagerMajority)
	// 	}
	// } else {
	//
	// }
	s.Proposals.AddEventUpdate(&pending, update)
	return pending.IncorporateVote(selfVote, s)
}

func (s *State) CancelEvent(cancel *actions.CancelEvent) error {
	event, ok := s.Events[cancel.Hash]
	if !ok {
		return fmt.Errorf("event not found: %v", crypto.EncodeHash(cancel.Hash))
	}
	if !event.Managers.IsMember(cancel.Author) {
		return errors.New("not a manager")
	}
	hash := cancel.Hashed()
	selfVote := actions.Vote{
		Epoch:   cancel.Epoch,
		Author:  cancel.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	pending := CancelEvent{
		Event:   event,
		Hash:    hash,
		Votes:   []actions.Vote{},
		Reasons: cancel.Reasons,
	}
	s.Proposals.AddCancelEvent(&pending, cancel)
	return pending.IncorporateVote(selfVote, s)
}

func (s *State) CreateEvent(create *actions.CreateEvent) error {
	if !s.IsMember(create.Author) {
		return errors.New("not a member")
	}
	collective, ok := s.Collective(create.OnBehalfOf)
	if !ok {
		return errors.New("collective not found")
	}
	hash := create.Hashed()
	vote := actions.Vote{
		Epoch:   create.Epoch,
		Author:  create.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	event := Event{
		Collective:     collective,
		StartAt:        create.StartAt,
		EstimatedEnd:   create.EstimatedEnd,
		Description:    create.Description,
		Venue:          create.Venue,
		Open:           create.Open,
		Public:         create.Public,
		Hash:           hash,
		Votes:          []actions.Vote{},
		Checkin:        make(map[crypto.Token]*Greeting),
		CheckinReasons: make(map[crypto.Token]string),
		Live:           false,
		EventReasons:   create.Reasons,
	}
	if len(create.Managers) > 0 {
		managers := make(map[crypto.Token]struct{})
		for _, manager := range create.Managers {
			managers[manager] = struct{}{}
		}
		event.Managers = &UnamedCollective{
			Members:  managers,
			Majority: int(create.ManagerMajority),
		}
	}
	if s.Proposals.Has(hash) {
		return errors.New("event already booked")
	}
	s.Proposals.AddEvent(&event, create)
	return event.IncorporateVote(vote, s)
}

func (s *State) MultipartMedia(media *actions.MultipartMedia) error {
	pending, ok := s.PendingMedia[media.Hash]
	if !ok {
		return errors.New("referred media not found")
	}
	total, err := pending.Append(media)
	if err != nil {
		return err
	}
	if total != nil {
		delete(s.PendingMedia, media.Hash)
		s.Media[media.Hash] = total
		//s.Notify(MediaUpload, media.Hash)
	}
	return nil
}

func (s *State) ReleaseDraft(release *actions.ReleaseDraft) error {
	draft, ok := s.Drafts[release.ContentHash]
	if !ok {
		return errors.New("draft not found")
	}
	if !draft.Authors.IsMember(release.Author) {
		return errors.New("not an author")
	}
	hash := release.Hashed()
	vote := actions.Vote{
		Epoch:   release.Epoch,
		Author:  release.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	newRelease := Release{
		Epoch:    release.Epoch,
		Draft:    draft,
		Hash:     hash,
		Votes:    []actions.Vote{},
		Released: false,
		Stamps:   make([]*Stamp, 0),
	}
	// if draft.Authors.Consensus(hash, []actions.Vote{vote}) {
	// 	if _, ok := s.Releases[release.ContentHash]; ok {
	// 		return errors.New("already released")
	// 	}
	// 	newRelease.Released = true
	// 	s.Releases[release.ContentHash] = &newRelease
	// 	fmt.Println("Released")
	// 	return nil
	// }
	text, _ := json.Marshal(newRelease)
	fmt.Println(string(text))
	s.Proposals.AddRelease(&newRelease, release)
	return newRelease.IncorporateVote(vote, s)
}

func (s *State) UpdateBoard(update *actions.UpdateBoard) error {
	if !s.IsMember(update.Author) {
		return errors.New("not a member")
	}
	board, ok := s.Board(update.Board)
	if !ok {
		return errors.New("board not found")
	}
	if (!board.Collective.IsMember(update.Author)) && (!board.Editors.IsMember(update.Author)) {
		return errors.New("not a member of collective or and editor")
	}
	// hash := crypto.Hasher([]byte(update.Serialize()))
	hash := update.Hashed()
	vote := actions.Vote{
		Epoch:   update.Epoch,
		Author:  update.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	// if board.Collective.Consensus(hash, []actions.Vote{vote}) {
	// 	if update.Keywords != nil {
	// 		board.Keyword = *update.Keywords
	// 	}
	// 	if update.Description != nil {
	// 		board.Description = *update.Description
	// 	}
	// 	if update.PinMajority != nil {
	// 		board.Editors.ChangeMajority(int(*update.PinMajority))
	// 	}
	// 	return nil
	// }
	pending := PendingUpdateBoard{
		Origin:      update,
		Keywords:    update.Keywords,
		Description: update.Description,
		PinMajority: update.PinMajority,
		Board:       board,
		Hash:        hash,
		Votes:       []actions.Vote{},
	}
	s.Proposals.AddPendingUpdateBoard(&pending, update)
	return pending.IncorporateVote(vote, s)
	// TODO notify
}

func (s *State) CreateBoard(board *actions.CreateBoard) error {
	if !s.IsMember(board.Author) {
		return errors.New("not a member")
	}
	if _, ok := s.Board(board.Name); ok {
		return errors.New("board already exists")
	}
	collective, ok := s.Collective(board.OnBehalfOf)
	if !ok {
		return errors.New("collective unkown")
	}
	hash := board.Hashed()
	newBoard := Board{
		Name:        board.Name,
		Keyword:     board.Keywords,
		Description: board.Description,
		Collective:  collective,
		Editors: &UnamedCollective{
			Members:  map[crypto.Token]struct{}{board.Author: {}},
			Majority: int(board.PinMajority),
		},
		Pinned: make([]*Draft, 0),
		Hash:   hash,
	}
	vote := actions.Vote{
		Epoch:   board.Epoch,
		Author:  board.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	// if collective.Consensus(hash, []actions.Vote{vote}) {
	// 	s.Boards[hash] = &newBoard
	// 	return nil
	// }
	pendingboard := &PendingBoard{
		Origin: board,
		Board:  &newBoard,
		Hash:   hash,
		Votes:  []actions.Vote{},
	}
	s.Proposals.AddPendingBoard(pendingboard, board)
	// TODO: notify
	return pendingboard.IncorporateVote(vote, s)
}

func (s *State) SignIn(signin *actions.Signin) error {
	hash := crypto.HashToken(signin.Author)
	if _, ok := s.Members[hash]; ok {
		return errors.New("already a member of synergy")
	}
	s.Members[hash] = signin.Handle
	s.MembersIndex[signin.Handle] = signin.Author // TODO: who guarantees single names?
	//s.Notify(SigninAction, hash)
	return nil
}

func (s *State) CreateCollective(create *actions.CreateCollective) error {
	if !s.IsMember(create.Author) {
		return errors.New("not a member of synergy")
	}
	if _, ok := s.Collective(create.Name); ok {
		return errors.New("collective already exists")
	}
	if create.Policy.Majority < 0 || create.Policy.Majority > 100 || create.Policy.SuperMajority < 0 || create.Policy.SuperMajority > 100 {
		return errors.New("invalid policy")
	}
	// hash := crypto.Hasher([]byte(create.Name))
	hash := create.Hashed()
	s.Collectives[hash] = &Collective{
		Name:        create.Name,
		Members:     map[crypto.Token]struct{}{create.Author: {}},
		Description: create.Description,
		Policy: actions.Policy{
			Majority:      create.Policy.Majority,
			SuperMajority: create.Policy.SuperMajority,
		},
	}
	return nil
}

func (s *State) UpdateCollective(update *actions.UpdateCollective) error {
	collective, ok := s.Collective(update.OnBehalfOf)
	if !ok {
		return errors.New("unkown collective")
	}
	if !collective.IsMember(update.Author) {
		return errors.New("not a member of collective")
	}
	hash := crypto.Hasher(update.Serialize()) // proposal hash = hash of instruction
	vote := actions.Vote{
		Epoch:   update.Epoch,
		Author:  update.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}

	// if update.Majority != nil || update.SuperMajority != nil {
	// 	if update.Majority != nil && (*update.Majority < 0 || *update.Majority > 100) {
	// 		return errors.New("invalid policy")
	// 	}
	// 	if update.SuperMajority != nil && (*update.SuperMajority < 0 || *update.SuperMajority > 100) {
	// 		return errors.New("invalid policy")
	// 	}
	// 	if collective.SuperConsensus(hash, []actions.Vote{vote}) {
	// 		if update.Description != nil {
	// 			collective.Description = *update.Description
	// 		}
	// 		newPolicy := actions.Policy{
	// 			Majority:      collective.Policy.Majority,
	// 			SuperMajority: collective.Policy.SuperMajority,
	// 		}

	// 		if update.Majority != nil {
	// 			newPolicy.Majority = int(*update.Majority)
	// 		}
	// 		if update.SuperMajority != nil {
	// 			newPolicy.SuperMajority = int(*update.SuperMajority)
	// 		}
	// 		collective.Policy = newPolicy
	// 		return nil
	// 	}
	// } else {
	// 	if collective.Consensus(hash, []actions.Vote{vote}) {
	// 		if update.Description != nil {
	// 			collective.Description = *update.Description
	// 		}
	// 		return nil
	// 	}
	// }

	pending := PendingUpdate{
		Update: update,
		// consensus is based on the collective composition at the moment
		// of incorporation of instruction
		Collective: collective.Photo(),
		Hash:       hash,
		Votes:      []actions.Vote{},
	}
	if update.Majority != nil || update.SuperMajority != nil {
		pending.ChangePolicy = true
	}
	s.Proposals.AddUpdateCollective(&pending, update)
	s.setDeadline(update.Epoch+ProposalDeadline, hash)
	return pending.IncorporateVote(vote, s)

}

func (s *State) RequestMembership(request *actions.RequestMembership) error {
	if !s.IsMember(request.Author) {
		return errors.New("not a member of synergy")
	}
	collective, ok := s.Collective(request.Collective)
	if !ok {
		return errors.New("collective not found")
	}
	if request.Include && collective.IsMember(request.Author) {
		return errors.New("already a member")
	}
	if (!request.Include) && (!collective.IsMember(request.Author)) {
		return errors.New("not a member of collective")
	}
	if !request.Include {
		delete(collective.Members, request.Author)
		return nil
	}
	// hash := crypto.Hasher(request.Serialize())
	hash := request.Hashed()
	pending := PendingRequestMembership{
		Request:    request,
		Collective: collective.Photo(),
		Hash:       hash,
		Votes:      make([]actions.Vote, 0),
	}
	vote := actions.Vote{
		Epoch:   request.Epoch,
		Author:  request.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	s.Proposals.AddRequestMembership(&pending, request)
	s.setDeadline(request.Epoch+ProposalDeadline, hash)
	return pending.IncorporateVote(vote, s)
}

func (s *State) RemoveMember(remove *actions.RemoveMember) error {
	collective, ok := s.Collective(remove.OnBehalfOf)
	if !ok {
		return errors.New("collective not found")
	}
	if !collective.IsMember(remove.Author) {
		return errors.New("author not a member of collective")
	}
	if !collective.IsMember(remove.Member) {
		return errors.New("member to be removed not a member of collective")
	}
	if remove.Author.Equal(remove.Member) {
		delete(collective.Members, remove.Author)
		if s.index != nil {
			s.index.IndexConsensus(remove.Hashed(), true)
		}
		return nil
	}
	// hash := crypto.Hasher(remove.Serialize())
	hash := remove.Hashed()
	vote := actions.Vote{
		Epoch:   remove.Epoch,
		Author:  remove.Author,
		Reasons: "commit",
		Hash:    hash,
		Approve: true,
	}
	// if collective.Consensus(hash, []actions.Vote{vote}) {
	// 	delete(collective.Members, remove.Author)
	// 	return nil
	// }
	pending := PendingRemoveMember{
		Remove:     remove,
		Collective: collective.Photo(),
		Hash:       hash,
		Votes:      []actions.Vote{},
	}
	s.Proposals.AddPendingRemoveMember(&pending, remove)
	s.setDeadline(remove.Epoch+ProposalDeadline, hash)
	return pending.IncorporateVote(vote, s)
}

func (s *State) React(reaction *actions.React) error {
	if reaction.Reaction >= ReactionsCount {
		return errors.New("invalid reaction")
	}
	// TODO: should check if hash is known?
	if count, ok := s.Reactions[reaction.Reaction][reaction.Hash]; ok {
		s.Reactions[reaction.Reaction][reaction.Hash] = count + 1
	} else {
		s.Reactions[reaction.Reaction][reaction.Hash] = 1
	}
	return nil
}

func (s *State) Edit(edit *actions.Edit) error {
	if _, ok := s.Media[edit.ContentHash]; ok {
		return errors.New("hash already claimed")
	}
	if _, ok := s.PendingMedia[edit.ContentHash]; ok {
		return errors.New("hash already claimed")
	}
	if edit.NumberOfParts > 1 {
		pending := PendingMedia{
			Hash:          edit.ContentHash,
			NumberOfParts: edit.NumberOfParts,
			Parts:         make([]*actions.MultipartMedia, int(edit.NumberOfParts)),
		}
		s.PendingMedia[edit.ContentHash] = &pending
		pending.Parts[0] = &actions.MultipartMedia{
			Hash: edit.ContentHash,
			Part: 0,
			Of:   edit.NumberOfParts,
			Data: edit.Content,
		}
	} else {
		s.Media[edit.ContentHash] = edit.Content
	}
	draft, ok := s.Drafts[edit.EditedDraft]
	if !ok {
		return errors.New("unkown draft")
	}

	newEdit := Edit{
		Date:     edit.Epoch,
		Reasons:  edit.Reasons,
		Draft:    draft,
		Edit:     edit.ContentHash,
		EditType: edit.ContentType,
		Votes:    make([]actions.Vote, 0),
	}
	newVote := actions.Vote{
		Epoch:   edit.Epoch,
		Author:  edit.Author,
		Reasons: "submission",
		Hash:    edit.ContentHash,
		Approve: true,
	}

	if edit.OnBehalfOf != "" {
		collective, ok := s.Collective(edit.OnBehalfOf)
		if !ok {
			return errors.New("collective unkown")
		}
		if !collective.IsMember(edit.Author) {
			return errors.New("not a member of collective")
		}
		newEdit.Authors = collective
		//if collective.Consensus(edit.ContentHash, newEdit.Votes) {
		//	s.Edits[edit.ContentHash] = &newEdit
		//	draft.Edits = append(draft.Edits, &newEdit)
		//} else {
		//	s.Proposals.AddEdit(&newEdit)
		//}
	} else if len(edit.CoAuthors) > 0 {
		newEdit.Authors = Authors(1+len(edit.CoAuthors), append(edit.CoAuthors, edit.Author)...)
		//s.Proposals.AddEdit(&newEdit)
	} else {
		newEdit.Authors = Authors(1, edit.Author)
		//s.Edits[edit.ContentHash] = &newEdit
		//draft.Edits = append(draft.Edits, &newEdit)
	}
	s.Proposals.AddEdit(&newEdit, edit)
	if err := newEdit.IncorporateVote(newVote, s); err != nil {
		return err
	}
	if s.index != nil {
		s.index.AddEditToIndex(&newEdit)
	}
	return nil
	//s.action.Notify(EditAction, DraftObject, edit.EditedDraft)
	//s.action.Notify(EditAction, AuthorObject, crypto.HashToken(edit.Author))
	//return nil
}

// IncorporateDraftInstruction checks if proposed draft is valid and if so
// incorporate it as ProposedDraft if further consent is necessary of as
// Draft if the instruction author has alone authority to submit the draft.
//
// Checks:
// a) It must refer to a known media file hash (pending media will not be
//
//	accepted)
//
// b) If it has a designated previous version, the instruction auhtor must be
//
//	an accredited author of the previous version or a member of the collective
//
// c) If draft is submitted on behalf of a named collective, this name must
//
//	be recognized by the state
//
// d)
func (s *State) Draft(draft *actions.Draft) error {
	if _, ok := s.Media[draft.ContentHash]; ok {
		return errors.New("media file already drafted")
	}
	if draft.NumberOfParts > 1 {
		pending := PendingMedia{
			Hash:          draft.ContentHash,
			NumberOfParts: draft.NumberOfParts,
			Parts:         make([]*actions.MultipartMedia, draft.NumberOfParts),
		}
		pending.Parts[0] = &actions.MultipartMedia{
			Epoch:  draft.Epoch,
			Author: draft.Author,
			Hash:   draft.ContentHash,
			Part:   0,
			Of:     draft.NumberOfParts,
			Data:   draft.Content,
		}
		s.PendingMedia[draft.ContentHash] = &pending
	} else {
		if !crypto.Hasher(draft.Content).Equal(draft.ContentHash) {
			return errors.New("hash does not match")
		}
		s.Media[draft.ContentHash] = draft.Content
	}
	var previous *Draft
	if draft.PreviousDraft != crypto.ZeroHash && draft.PreviousDraft != crypto.ZeroValueHash {
		if previous, ok := s.Drafts[draft.PreviousDraft]; !ok {
			return errors.New("invalid previous version")
		} else {
			isPreviousAuthor := previous.Authors.IsMember(draft.Author)
			if !isPreviousAuthor {
				return errors.New("unauthorized version")
			}
		}
	}
	selfVote := actions.Vote{
		Epoch:   draft.Epoch,
		Author:  draft.Author,
		Reasons: "submission",
		Hash:    draft.ContentHash,
		Approve: true,
	}

	newDraft := &Draft{
		Title:           draft.Title,
		Description:     draft.Description,
		DraftType:       draft.ContentType,
		DraftHash:       draft.ContentHash,
		Keywords:        draft.Keywords,
		PreviousVersion: previous,
		References:      draft.References,
		Votes:           make([]actions.Vote, 0),
	}
	if len(draft.CoAuthors) == 0 {
		if draft.OnBehalfOf == "" {
			// create single author collective
			newDraft.Authors = Authors(1, draft.Author)
			//newDraft.Aproved = true
			//s.Drafts[newDraft.DraftHash] = newDraft
		} else {
			behalf, ok := s.Collective(draft.OnBehalfOf)
			if !ok {
				return errors.New("named collective not recognizedx")
			}
			newDraft.Authors = behalf
			//if behalf.Consensus(newDraft.DraftHash, newDraft.Votes) {
			//	newDraft.Aproved = true
			//	s.Drafts[newDraft.DraftHash] = newDraft
			//} else {
			//	s.Proposals.AddDraft(newDraft)
			//}
		}
	} else {
		coauthors := []crypto.Token{draft.Author}
		coauthors = append(coauthors, draft.CoAuthors...)
		if draft.Policy == nil {
			newDraft.Authors = Authors(100, coauthors...)
		} else {
			newDraft.Authors = Authors(draft.Policy.Majority, coauthors...)
		}
		//s.Proposals.AddDraft(newDraft)
	}
	//if newDraft.PreviousVersion != nil {
	//	s.action.Notify(DraftAction, DraftObject, draft.PreviousDraft)
	//}
	s.Proposals.AddDraft(newDraft, draft)
	if err := newDraft.IncorporateVote(selfVote, s); err != nil {
		return err
	}
	if s.index != nil {
		s.index.AddDraftToIndex(newDraft)
	}
	return nil
}

func (s *State) Vote(vote *actions.Vote) error {
	if draft, ok := s.Drafts[vote.Hash]; ok {
		return draft.IncorporateVote(*vote, s)
	}
	return s.Proposals.IncorporateVote(*vote, s)
}

func (s *State) Pin(pin *actions.Pin) error {
	// existe o board no state?
	board, ok := s.Board(pin.Board)
	if !ok {
		return errors.New("invalid board")
	}
	// existe o draft no state?
	draft, ok := s.Drafts[pin.Draft]
	if !ok {
		return errors.New("invalid draft")
	}
	// criando o byte array pra gerar o hash
	hash := pin.Hashed()

	action := Pin{
		Hash:  hash,
		Epoch: pin.Epoch,
		Board: board,
		Draft: draft,
		Pin:   pin.Pin,
		Votes: make([]actions.Vote, 0),
	}
	selfVote := actions.Vote{
		Epoch:   pin.Epoch,
		Author:  pin.Author,
		Reasons: "submission",
		Hash:    hash,
		Approve: true,
	}
	// coloca a proposta de pin criado nos proposals
	s.Proposals.AddPin(&action, pin)
	return action.IncorporateVote(selfVote, s)
}

func (s *State) BoardEditor(action *actions.BoardEditor) error {
	board, ok := s.Board(action.Board)
	if !ok {
		return errors.New("invalid board")
	}
	if s.IsMember(action.Editor); !ok { // should be
		return errors.New("invalid editor")
	}
	hash := action.Hashed()
	proposal := BoardEditor{
		Hash:   hash,
		Epoch:  action.Epoch,
		Board:  board,
		Editor: action.Editor,
		Insert: action.Insert,
		Votes:  make([]actions.Vote, 0),
	}
	selfVote := actions.Vote{
		Epoch:   action.Epoch,
		Author:  action.Author,
		Reasons: "submission",
		Hash:    hash,
		Approve: true,
	}
	s.Proposals.AddBoardEditor(&proposal, action)
	return proposal.IncorporateVote(selfVote, s)
}
