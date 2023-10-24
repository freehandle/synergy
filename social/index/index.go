package index

import (
	"fmt"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
	"github.com/freehandle/synergy/social/state"
)

const MaxRecentAction = 100

type LastAction struct {
	Author      crypto.Token
	Description string
	Epoch       uint64
}

type IndexedAction struct {
	Action   actions.Action
	Hash     crypto.Hash
	Approved byte
}

const ActionsCacheCount = 10

type ActionDetails struct {
	Description string
	ObjectHash  string
	Author      crypto.Token
	Votes       []actions.Vote
	VoteStatus  bool
	Epoch       uint64
	IsReaction  bool
	Reaction    string
}

type RecentActions struct {
	actions []actions.Action
}

func NewRecentActions(action actions.Action) *RecentActions {
	return &RecentActions{actions: []actions.Action{action}}
}

func (r *RecentActions) Append(action actions.Action) {
	if len(r.actions) == ActionsCacheCount {
		r.actions = append(r.actions[1:], action)
	} else {
		r.actions = append(r.actions, action)
	}
}

func (r *RecentActions) Last() actions.Action {
	if len(r.actions) == 0 {
		return nil
	}
	return r.actions[len(r.actions)-1]
}

type Index struct {
	allPendingactions map[crypto.Hash]actions.Action

	allUsers map[crypto.Token]*Person

	indexedMembers      map[crypto.Token]string           // token to handle
	memberToAction      map[crypto.Token][]*IndexedAction // ação e se foi aprovada ou se está pendente
	pendingIndexActions map[crypto.Hash]crypto.Token      // action

	indexVotes          map[crypto.Token]*SetOfHashes
	indexCompletedVotes map[crypto.Hash][]actions.Vote

	// central connections member connections
	memberToCollective map[crypto.Token][]string
	memberToBoard      map[crypto.Token][]string
	memberToEvent      map[crypto.Token][]crypto.Hash
	MemberToDraft      map[crypto.Token][]*state.Draft
	MemberToEdit       map[crypto.Token][]*state.Edit

	MemberToCheckin map[crypto.Token][]*state.Event
	//memberToEdit
	//memberToDraft

	objectHashToActionHash map[crypto.Hash]*RecentActions // object to recent actions

	// central connections collectives card
	collectiveToBoards map[*state.Collective][]*state.Board
	collectiveToStamps map[*state.Collective][]*state.Stamp
	collectiveToEvents map[*state.Collective][]*state.Event

	RecentActions []*IndexedAction

	// collectiveLastAction map[*state.Collective][]lastaction

	// central connections edit card
	//editToDrafts map[*state.Edit][]*state.Draft

	state *state.State

	stateProposals *state.Proposals
}

func (i *Index) Reason(hash crypto.Hash) string {
	actionWithHash := i.allPendingactions[hash]
	if actionWithHash == nil {
		return ""
	}
	return actionWithHash.Reasoning()
}

func (i *Index) PendingAction(hash crypto.Hash) actions.Action {
	return i.allPendingactions[hash]
}

type Person struct {
	Collectives []string
	Boards      []string
	Events      []crypto.Hash
	Drafts      []crypto.Hash
	Edits       []crypto.Hash
}

func (p *Person) AddCollective(collective string) {
	for _, b := range p.Collectives {
		if b == collective {
			return
		}
	}
	p.Collectives = append(p.Collectives, collective)
}

func (p *Person) RemoveCollective(collective string) {
	for n, b := range p.Collectives {
		if b == collective {
			p.Collectives = append(p.Collectives[:n], p.Collectives[n+1:]...)
			return
		}
	}
}

func (p *Person) AddBoard(board string) {
	for _, b := range p.Boards {
		if b == board {
			return
		}
	}
	p.Boards = append(p.Boards, board)
}

func (p *Person) RemoveBoard(board string) {
	for n, b := range p.Boards {
		if b == board {
			p.Boards = append(p.Boards[:n], p.Boards[n+1:]...)
			return
		}
	}
}

func (p *Person) AddEvent(event crypto.Hash) {
	for _, e := range p.Events {
		if e.Equal(event) {
			return
		}
	}
	p.Events = append(p.Events, event)
}

func (p *Person) AddDraft(draft crypto.Hash) {
	for _, e := range p.Drafts {
		if e.Equal(draft) {
			return
		}
	}
	p.Drafts = append(p.Drafts, draft)
}

func (p *Person) AddEdit(edit crypto.Hash) {
	for _, e := range p.Edits {
		if e.Equal(edit) {
			return
		}
	}
	p.Edits = append(p.Edits, edit)
}

func (i *Index) Personal(token crypto.Token) *Person {
	person, ok := i.allUsers[token]
	if !ok {
		person = &Person{
			Collectives: make([]string, 0),
			Boards:      make([]string, 0),
			Events:      make([]crypto.Hash, 0),
			Drafts:      make([]crypto.Hash, 0),
			Edits:       make([]crypto.Hash, 0),
		}
		i.allUsers[token] = person
	}
	return person
}

func (i *Index) IndexActionToPerson(hash crypto.Hash) {
	action, _ := i.allPendingactions[hash]
	if action == nil {
		return
	}
	switch v := action.(type) {
	case *actions.CreateCollective:
		person := i.Personal(v.Author)
		person.AddCollective(v.Name)
	case *actions.RequestMembership:
		person := i.Personal(v.Author)
		person.AddCollective(v.Collective)
	case *actions.Draft:
		if len(v.OnBehalfOf) == 0 {
			person := i.Personal(v.Author)
			person.AddDraft(v.ContentHash)
			for _, author := range v.CoAuthors {
				person := i.Personal(author)
				person.AddDraft(v.ContentHash)
			}
		}
	case *actions.Edit:
		if len(v.OnBehalfOf) == 0 {
			person := i.Personal(v.Author)
			person.AddDraft(v.ContentHash)
			for _, author := range v.CoAuthors {
				person := i.Personal(author)
				person.AddDraft(v.ContentHash)
			}
		}
	case *actions.CreateBoard:
		person := i.Personal(v.Author)
		person.AddBoard(v.Name)
	case *actions.BoardEditor:
		person := i.Personal(v.Editor)
		if v.Insert {
			person.AddBoard(v.Board)
		} else {
			person.RemoveBoard(v.Board)
		}
	case *actions.RemoveMember:
		person := i.Personal(v.Member)
		person.RemoveBoard(v.OnBehalfOf)
	case *actions.CreateEvent:
		person := i.Personal(v.Author)
		person.AddEvent(v.Hashed())
	case *actions.CheckinEvent:
		person := i.Personal(v.Author)
		person.AddEvent(v.EventHash)
	}
}

func NewIndex() *Index {
	return &Index{

		allPendingactions: make(map[crypto.Hash]actions.Action),

		allUsers: make(map[crypto.Token]*Person),
		// central connections
		memberToCollective: make(map[crypto.Token][]string),
		memberToBoard:      make(map[crypto.Token][]string),
		memberToEvent:      make(map[crypto.Token][]crypto.Hash),
		MemberToCheckin:    make(map[crypto.Token][]*state.Event),

		MemberToDraft: make(map[crypto.Token][]*state.Draft),
		MemberToEdit:  make(map[crypto.Token][]*state.Edit),

		//memberToEdit:       make(map[string][]*state.Edit),
		collectiveToBoards: make(map[*state.Collective][]*state.Board),
		collectiveToStamps: make(map[*state.Collective][]*state.Stamp),
		collectiveToEvents: make(map[*state.Collective][]*state.Event),
		// collectiveLastAction: make(map[*state.Collective][]lastaction),
		//editToDrafts: make(map[*state.Edit][]*state.Draft),

		indexedMembers:      make(map[crypto.Token]string),
		memberToAction:      make(map[crypto.Token][]*IndexedAction),
		pendingIndexActions: make(map[crypto.Hash]crypto.Token),

		indexVotes:          make(map[crypto.Token]*SetOfHashes),
		indexCompletedVotes: make(map[crypto.Hash][]actions.Vote),

		objectHashToActionHash: make(map[crypto.Hash]*RecentActions),

		RecentActions: make([]*IndexedAction, 0),
	}
}

func (i *Index) SetState(s *state.State) {
	i.state = s
	i.stateProposals = s.Proposals
}

func (i *Index) ActionStatus(action actions.Action) ([]actions.Vote, bool) {
	hash := action.Hashed()
	if votes, ok := i.indexCompletedVotes[hash]; ok {
		return votes, true
	}
	return i.stateProposals.Votes(hash), false
}

// Objects related to a given collective

func (i *Index) BoardsOnCollective(collective *state.Collective) []*state.Board {
	return i.collectiveToBoards[collective]
}

func (i *Index) StampsOnCollective(collective *state.Collective) []*state.Stamp {
	return i.collectiveToStamps[collective]
}

func (i *Index) EventsOnCollective(collective *state.Collective) []*state.Event {
	return i.collectiveToEvents[collective]
}

// Objects related to a given member

func (i *Index) CollectivesOnMember(member crypto.Token) []string {
	return i.memberToCollective[member]
}

func (i *Index) BoardsOnMember(member crypto.Token) []string {
	return i.memberToBoard[member]
}

func (i *Index) EventsOnMember(member crypto.Token) []crypto.Hash {
	return i.memberToEvent[member]
}

func (i *Index) AddMemberToIndex(token crypto.Token, handle string) {
	i.indexedMembers[token] = handle
}

func (i Index) GetLastAction(objectHash crypto.Hash) *ActionDetails {
	recent := i.objectHashToActionHash[objectHash]
	if recent == nil || len(recent.actions) == 0 {
		return nil
	}
	// TODO: check consensus status
	des, hash, author, epoch, _ := i.ActionToString(recent.actions[len(recent.actions)-1], true)
	return &ActionDetails{
		Description: des,
		Author:      author,
		ObjectHash:  hash,
		Epoch:       epoch,
	}
}

func (i Index) GetRecentActions(objectHash crypto.Hash) []ActionDetails {
	recent := i.objectHashToActionHash[objectHash]
	if recent == nil {
		return nil
	}
	details := make([]ActionDetails, len(recent.actions))
	for n, r := range recent.actions {
		// TODO: check consensus status
		status := true
		des, _, _, epoch, _ := i.ActionToString(r, status)
		votes, status := i.ActionStatus(r)
		details[n] = ActionDetails{
			Description: des,
			Votes:       votes,
			VoteStatus:  status,
			Epoch:       epoch,
		}
	}
	return details
}

func (i Index) GetRecentActionsWithLinks(objectHash crypto.Hash) []ActionDetails {
	recent := i.objectHashToActionHash[objectHash]
	if recent == nil {
		return nil
	}
	actionDetails := make([]ActionDetails, len(recent.actions))
	for n, r := range recent.actions {
		// TODO: check consensus status
		votes, status := i.ActionStatus(r)
		des, epoch, reasons := i.ActionToStringWithLinks(r, status)
		details := ActionDetails{
			Description: des,
			Votes:       votes,
			VoteStatus:  status,
			Epoch:       epoch,
		}
		if _, ok := r.(*actions.React); ok {
			details.IsReaction = true
			details.Reaction = reasons
		}
		actionDetails[n] = details
	}
	return actionDetails
}

func (i *Index) IndexAction(action actions.Action) {
	hash := action.Hashed()
	i.allPendingactions[hash] = action
	author := action.Authored()
	objects := i.ActionToObjects(action)
	for _, object := range objects {
		if recent, ok := i.objectHashToActionHash[object]; ok {
			recent.Append(action)
		} else {
			i.objectHashToActionHash[object] = NewRecentActions(action)
		}
	}
	//hash := action.Hashed()
	newAction := IndexedAction{
		Action:   action,
		Hash:     action.Hashed(),
		Approved: 0,
	}
	if _, ok := i.indexedMembers[author]; ok {
		switch v := action.(type) {
		case *actions.GreetCheckinEvent:
			newAction.Approved = 1
		case *actions.CheckinEvent:
			newAction.Approved = 1
		case *actions.React:
			newAction.Approved = 1
		case *actions.Signin:
			newAction.Approved = 1
		case *actions.Vote:
			newAction.Approved = 1
		case *actions.RequestMembership:
			if !v.Include {
				i.IndexConsensusAction(action)
				newAction.Approved = 1
			}
		case *actions.CreateCollective:
			i.IndexConsensusAction(action)
			newAction.Approved = 1
			if person := i.Personal(v.Author); person != nil {
				person.AddCollective(v.Name)
			}
		}
		if indexedActions, ok := i.memberToAction[author]; ok {
			i.memberToAction[author] = append(indexedActions, &newAction)
		} else {
			i.memberToAction[author] = []*IndexedAction{&newAction}
		}
		if newAction.Approved == 0 {
			hash := action.Hashed()
			i.pendingIndexActions[hash] = author
		}
	}
	if len(i.RecentActions) > MaxRecentAction {
		i.RecentActions = append([]*IndexedAction{&newAction}, i.RecentActions[0:MaxRecentAction-1]...)
	} else {
		i.RecentActions = append([]*IndexedAction{&newAction}, i.RecentActions[0:]...)
	}
}

func (i *Index) isIndexedMember(token crypto.Token) bool {
	_, ok := i.indexedMembers[token]
	return ok
}

func appendOrCreate[T any](values []T, value T) []T {
	if values == nil {
		return []T{value}
	}
	return append(values, value)
}

func removeItem[T comparable](values []T, value T) []T {
	for n, item := range values {
		if item == value {
			if n == len(values)-1 {
				return values[0:n]
			}
			return append(values[0:n], values[n+1:]...)
		}
	}
	return values
}

func (i *Index) IndexConsensusAction(action actions.Action) {
	switch v := action.(type) {
	case *actions.CreateCollective:
		if i.isIndexedMember(v.Author) {
			i.memberToCollective[v.Author] = appendOrCreate[string](i.memberToCollective[v.Author], v.Name)
		}
	case *actions.RequestMembership:
		if i.isIndexedMember(v.Author) {
			if v.Include {
				i.memberToCollective[v.Author] = appendOrCreate[string](i.memberToCollective[v.Author], v.Collective)
			} else {
				i.memberToCollective[v.Author] = removeItem[string](i.memberToCollective[v.Author], v.Collective)
			}
		}
	case *actions.RemoveMember:
		if i.isIndexedMember(v.Member) {
			i.memberToCollective[v.Member] = removeItem[string](i.memberToCollective[v.Member], v.OnBehalfOf)
		}
	case *actions.CreateBoard:
		if i.isIndexedMember(v.Author) {
			i.memberToBoard[v.Author] = appendOrCreate[string](i.memberToBoard[v.Author], v.Name)
		}
	case *actions.CreateEvent:
		if i.isIndexedMember(v.Author) {
			hash := crypto.Hasher(v.Serialize())
			i.memberToEvent[v.Author] = appendOrCreate[crypto.Hash](i.memberToEvent[v.Author], hash)
		}
	case *actions.UpdateEvent:
		if v.Managers != nil {
			if event, ok := i.state.Events[v.EventHash]; ok {
				oldmanagers := event.Managers.ListOfMembers()
				for _, manager := range *v.Managers {
					if _, ok := oldmanagers[manager]; !ok {
						if i.isIndexedMember(manager) {
							i.memberToEvent[manager] = appendOrCreate[crypto.Hash](i.memberToEvent[manager], v.EventHash)
						}
					} else {
						delete(oldmanagers, manager)
					}
				}
				for manager := range oldmanagers {
					if i.isIndexedMember(manager) {
						i.memberToEvent[manager] = removeItem[crypto.Hash](i.memberToEvent[manager], v.EventHash)
					}
				}
			}
		}
	}
}

func (i *Index) IndexConsensus(hash crypto.Hash, approved bool) {
	if approved {
		i.IndexActionToPerson(hash)
	}
	author, ok := i.pendingIndexActions[hash]
	if !ok {
		return
	}
	delete(i.pendingIndexActions, hash)
	indexActions, ok := i.memberToAction[author]
	if !ok {
		return
	}
	for _, action := range indexActions {
		if action.Hash.Equal(hash) {
			if approved {
				i.IndexConsensusAction(action.Action)
				action.Approved = 1
			} else {
				action.Approved = 2
			}
		}
	}
}

func (i *Index) LastManagerPinOnBoard(manager crypto.Token, board string) *LastAction {
	allActions, ok := i.memberToAction[manager]
	if !ok {
		return nil
	}
	for n := len(allActions) - 1; n >= 0; n-- {
		switch v := allActions[n].Action.(type) {
		case *actions.Pin:
			if v.Board == board {
				actions := &LastAction{
					Author: manager,
					Epoch:  v.Epoch,
				}
				if draft, ok := i.state.Drafts[v.Draft]; ok {
					pin := "pin"
					if !v.Pin {
						pin = "unpin"
					}
					actions.Description = fmt.Sprintf("%s %s", pin, draft.Title)
					return actions
				}
			}
		}
	}
	return nil
}

func (i *Index) LastPinOnBoard(token crypto.Token, board string) *LastAction {
	recent, ok := i.objectHashToActionHash[crypto.Hasher([]byte(board))]
	if !ok || recent == nil || len(recent.actions) == 0 {
		return nil
	}
	for n := len(recent.actions) - 1; n >= 0; n-- {
		switch v := recent.actions[n].(type) {
		case *actions.Pin:
			if v.Board == board && (!v.Author.Equal(token)) {
				actions := &LastAction{
					Author: v.Author,
					Epoch:  v.Epoch,
				}
				if draft, ok := i.state.Drafts[v.Draft]; ok {
					pin := "pin"
					if !v.Pin {
						pin = "unpin"
					}
					actions.Description = fmt.Sprintf("%s %s", pin, draft.Title)
					return actions
				}
			}
		}
	}
	return nil
}

type PendingAction struct {
	Description string
	Epoch       uint64
	Votes       []actions.Vote
}

type PendingActionDetailed struct {
	Description string
	Epoch       uint64
	Pool        *state.Pool
}

func (i *Index) GetPendingActions(token crypto.Token) []PendingAction {
	pendingActions := make([]PendingAction, 0)
	actions := i.memberToAction[token]
	if actions == nil {
		return pendingActions
	}
	for _, action := range actions {
		if action.Approved == 0 {
			description, _, _, epoch, _ := i.ActionToString(action.Action, false)
			votes := i.state.Proposals.Votes(action.Hash)
			pending := PendingAction{
				Description: description,
				Epoch:       epoch,
				Votes:       votes,
			}
			pendingActions = append(pendingActions, pending)
		}
	}
	return pendingActions
}

func (i *Index) GetPendingActionsDetailed(token crypto.Token) []PendingActionDetailed {
	pendingActions := make([]PendingActionDetailed, 0)
	actions := i.memberToAction[token]
	if actions == nil {
		return pendingActions
	}
	for _, action := range actions {
		if action.Approved == 0 {
			//description, _, _, epoch, _ := i.ActionToString(action.Action, false)
			description, epoch, _ := i.ActionToStringWithLinks(action.Action, false)
			if pool := i.state.Proposals.Pooling(action.Hash); pool != nil {
				pending := PendingActionDetailed{
					Description: description,
					Epoch:       epoch,
					Pool:        pool,
				}
				pendingActions = append(pendingActions, pending)
			}
		}
	}
	return pendingActions
}

func (i *Index) LastMemberActionOnCollective(member crypto.Token, collective string) *LastAction {
	allActions, ok := i.memberToAction[member]
	if !ok {
		return nil
	}
	for n := len(allActions) - 1; n >= 0; n-- {
		switch v := allActions[n].Action.(type) {
		case *actions.CreateBoard:
			if v.OnBehalfOf == collective {
				return &LastAction{
					Author:      member,
					Description: "create board",
					Epoch:       v.Epoch,
				}
			}
		case *actions.Draft:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "submit draft",
					Epoch:       v.Epoch,
				}
			}
		case *actions.Edit:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "submit edit",
					Epoch:       v.Epoch,
				}
			}

		case *actions.CreateEvent:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "create event",
					Epoch:       v.Epoch,
				}
			}
		case *actions.RemoveMember:
			if v.OnBehalfOf == collective && v.Authored().Equal(member) {
				return &LastAction{
					Author:      member,
					Description: "remove member",
					Epoch:       v.Epoch,
				}
			}
		}
	}
	return nil
}

func (i *Index) AddCheckin(token crypto.Token, event *state.Event) {
	if person := i.Personal(token); person != nil {
		person.AddEvent(event.Hash)
	}
	if _, ok := i.indexedMembers[token]; ok {
		if events, ok := i.MemberToCheckin[token]; ok {
			i.MemberToCheckin[token] = append(events, event)
		} else {
			i.MemberToCheckin[token] = []*state.Event{event}
		}
	}
}

func (i *Index) AddBoardToCollective(board *state.Board, collective *state.Collective) {
	if boards, ok := i.collectiveToBoards[collective]; ok {
		i.collectiveToBoards[collective] = append(boards, board)
	} else {
		i.collectiveToBoards[collective] = []*state.Board{board}
	}
}

func (i *Index) RemoveBoardFromCollective(board *state.Board, collective *state.Collective) {
	if boards, ok := i.collectiveToBoards[collective]; ok {
		for n, e := range boards {
			if e == board {
				removed := boards[0:n]
				if n < len(boards)-1 {
					removed = append(removed, boards[n+1:]...)
				}
				i.collectiveToBoards[collective] = removed
			}
		}
	}
}

func (i *Index) AddStampToCollective(stamp *state.Stamp, collective *state.Collective) {
	if stamps, ok := i.collectiveToStamps[collective]; ok {
		i.collectiveToStamps[collective] = append(stamps, stamp)
	} else {
		i.collectiveToStamps[collective] = []*state.Stamp{stamp}
	}
}

func (i *Index) AddEventToCollective(event *state.Event, collective *state.Collective) {
	if events, ok := i.collectiveToEvents[collective]; ok {
		i.collectiveToEvents[collective] = append(events, event)
	} else {
		i.collectiveToEvents[collective] = []*state.Event{event}
	}
}

func (i *Index) RemoveEventFromCollective(event *state.Event, collective *state.Collective) {
	if events, ok := i.collectiveToEvents[collective]; ok {
		for n, e := range events {
			if e == event {
				removed := events[0:n]
				if n < len(events)-1 {
					removed = append(removed, events[n+1:]...)
				}
				i.collectiveToEvents[collective] = removed
			}
		}
	}
}

func (i *Index) AddDraftToIndex(draft *state.Draft) {
	if draft.Authors == nil {
		return
	}
	tokens := draft.Authors.ListOfMembers()
	for token := range tokens {
		if _, ok := i.indexedMembers[token]; ok {
			if drafts, ok := i.MemberToDraft[token]; ok {
				i.MemberToDraft[token] = append(drafts, draft)
			} else {
				i.MemberToDraft[token] = []*state.Draft{draft}
			}
		}
	}
}

func (i *Index) AddEditToIndex(edit *state.Edit) {
	if edit.Authors == nil {
		return
	}
	tokens := edit.Authors.ListOfMembers()
	for token := range tokens {
		if _, ok := i.indexedMembers[token]; ok {
			if drafts, ok := i.MemberToEdit[token]; ok {
				i.MemberToEdit[token] = append(drafts, edit)
			} else {
				i.MemberToEdit[token] = []*state.Edit{edit}
			}
		}
	}
}

/*
func (i *Index) RemoveMemberFromCollective(collective *state.Collective, member crypto.Token) {
	delete(i.memberToCollective, member)
}

func (i *Index) AddMemberToCollective(collective *state.Collective, member crypto.Token) {
	if collectives, ok := i.memberToCollective[member]; ok {
		i.memberToCollective[member] = append(collectives, collective)
	} else {
		i.memberToCollective[member] = []*state.Collective{collective}
	}
}

func (i *Index) AddEditorToBoard(board *state.Board, editor crypto.Token) {
	if boards, ok := i.memberToBoard[editor]; ok {
		i.memberToBoard[editor] = append(boards, board)
	} else {
		i.memberToBoard[editor] = []*state.Board{board}
	}
}

func (i *Index) RemoveEditorFromBoard(board *state.Board, editor crypto.Token) {
	delete(i.memberToBoard, editor)
}
*/

// Collective's boards

/*func (i *Index) AddBoardToCollective(board *state.Board, collective *state.Collective) {if boards, ok := i.collectiveToBoards[collective]; ok {
		i.collectiveToBoards[collective] = append(boards, board)
	} else {
		i.collectiveToBoards[collective] = []*state.Board{board}
	}
}

func (i *Index) RemoveBoardFromCollective(board *state.Board, collective *state.Collective) {
	if boards, ok := i.collectiveToBoards[collective]; ok {
		for n, e := range boards {
			if e == board {
				removed := boards[0:n]
				if n < len(boards)-1 {
					removed = append(removed, boards[n+1:]...)
				}
				i.collectiveToBoards[collective] = removed
			}
		}
	}
}

// Collective's stamps


// Collective's events


// Edit's drafts

func (i *Index) AddDraftToEdit(draft *state.Draft, edit *state.Edit) {
	if drafts, ok := i.editToDrafts[edit]; ok {
		i.editToDrafts[edit] = append(drafts, draft)
	} else {
		i.editToDrafts[edit] = []*state.Draft{draft}
	}
}
*/
