package api

import (
	"net/url"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/crypto/dh"
	"github.com/freehandle/synergy/social/index"
	"github.com/freehandle/synergy/social/state"
)

type EventsView struct {
	Hash        string
	Live        bool
	Description string
	StartAt     time.Time
	Collective  NameLink
	Public      bool
}

type EventVoteAction struct {
	Kind   string // create, cancel, update
	Hash   string
	Update string
}

type EventsListView struct {
	Events []EventsView
	Head   HeaderInfo
}

type VoteUpdateEventView struct {
	Description     string
	OldDescription  string
	StartAt         string
	OldStartAt      string
	EstimatedEnd    string
	OldEstimatedEnd string
	Venue           string
	OldVenue        string
	Open            string
	OldOpen         string
	Public          string
	OldPublic       string
	Hash            string
	Reasons         string
	Collective      string
	CollectiveLink  string
	Managing        bool
	VoteHash        string
	Head            HeaderInfo
	Voting          DetailedVoteView
}

func yesorno(b *bool) string {
	if b == nil {
		return ""
	}
	if *b {
		return "yes"
	}
	return "no"
}

func EventUpdateFromState(s *state.State, hash crypto.Hash, token crypto.Token) VoteUpdateEventView {
	update, ok := s.Proposals.UpdateEvent[hash]
	if !ok {
		return VoteUpdateEventView{}
	}
	old := update.Event
	head := HeaderInfo{
		Active:  "MyEvents",
		Path:    "venture / my events / ",
		EndPath: "update event " + old.StartAt.Format("2006-01-02") + " by " + LimitStringSize(old.Collective.Name, maxStringSize),
		Section: "explore",
	}
	vote := VoteUpdateEventView{
		OldDescription:  old.Description,
		OldStartAt:      old.StartAt.String(),
		OldEstimatedEnd: old.EstimatedEnd.String(),
		OldVenue:        old.Venue,
		OldOpen:         yesorno(&old.Open),
		OldPublic:       yesorno(&old.Public),
		Open:            yesorno(update.Open),
		Public:          yesorno(update.Public),
		Hash:            crypto.EncodeHash(old.Hash),
		Reasons:         update.Reasons,
		Collective:      old.Collective.Name,
		CollectiveLink:  url.QueryEscape(old.Collective.Name),
		VoteHash:        crypto.EncodeHash(hash),
		Head:            head,
		Voting:          NewDetailedVoteView(update.Votes, update.Event.Managers, s),
	}
	if update.Description != nil {
		vote.Description = *update.Description
	}
	if update.StartAt != nil {
		vote.StartAt = update.StartAt.String()
	}
	if update.EstimatedEnd != nil {
		vote.StartAt = update.StartAt.String()
	}
	if update.Venue != nil {
		vote.StartAt = *update.Venue
	}
	if old.Managers.IsMember(token) {
		vote.Managing = true
		vote.Head = HeaderInfo{
			Active:  "MyEvents",
			Path:    "venture / my events / ",
			EndPath: old.StartAt.Format("2006-01-02") + " by " + LimitStringSize(old.Collective.Name, maxStringSize),
			Section: "venture",
		}
	} else {
		vote.Head = HeaderInfo{
			Active:  "Events",
			Path:    "explore / events / ",
			EndPath: old.StartAt.Format("2006-01-02") + " by " + LimitStringSize(old.Collective.Name, maxStringSize),
			Section: "explore",
		}
	}

	return vote
}

type EventDetailView struct {
	Live               bool
	Description        string
	StartAt            time.Time
	EstimatedEnd       time.Time
	Collective         NameLink
	MemberOfCollective bool
	Venue              string
	Open               bool
	Public             bool
	ManagerMajority    int
	Managers           []MemberDetailView
	Checkedin          []CheckInDetails
	Votes              DetailedVoteView
	Managing           bool
	Hash               string
	Greeted            []MemberDetailView
	MyGreeting         string
	Head               HeaderInfo
	EventReasons       string
}

func PendingEventFromState(s *state.State, i *index.Index, hash crypto.Hash) *EventDetailView {
	event, ok := s.Proposals.CreateEvent[hash]
	if !ok {
		return &EventDetailView{
			Head: HeaderInfo{},
		}
	}
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / connections / collectives / " + LimitStringSize(event.Collective.Name, maxStringSize) + " / ",
		EndPath: "create event",
		Section: "venture",
	}

	view := EventDetailView{
		Head:            head,
		Live:            event.Live,
		Description:     event.Description,
		StartAt:         event.StartAt,
		Collective:      NameLinker(event.Collective.Name),
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		ManagerMajority: event.Managers.Majority,
		Managers:        make([]MemberDetailView, 0),
		Votes:           NewDetailedVoteView(event.Votes, event.Collective, s),
		Hash:            crypto.EncodeHash(hash),
		EventReasons:    event.EventReasons,
	}
	return &view
}

func CancelEventFromState(s *state.State, i *index.Index, hash crypto.Hash) *EventDetailView {
	cancel, ok := s.Proposals.CancelEvent[hash]
	if !ok {
		return nil
	}
	event := cancel.Event
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / connections / collectives / " + LimitStringSize(event.Collective.Name, maxStringSize) + " / ",
		EndPath: "create event",
		Section: "venture",
	}

	view := EventDetailView{
		Head:            head,
		Live:            event.Live,
		Description:     event.Description,
		StartAt:         event.StartAt,
		Collective:      NameLinker(event.Collective.Name),
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		ManagerMajority: event.Managers.Majority,
		Managers:        make([]MemberDetailView, 0),
		Votes:           NewDetailedVoteView(cancel.Votes, event.Collective, s),
		Hash:            crypto.EncodeHash(hash),
		EventReasons:    cancel.Reasons,
	}
	return &view
}

func EventsFromState(state *state.State) EventsListView {
	head := HeaderInfo{
		Active:  "Events",
		Path:    "explore / ",
		EndPath: "events",
		Section: "explore",
	}
	view := EventsListView{
		Head:   head,
		Events: make([]EventsView, 0),
	}
	for _, event := range state.Events {
		itemView := EventsView{
			Hash: crypto.EncodeHash(event.Hash),
			Live: event.Live,

			Description: event.Description,
			StartAt:     event.StartAt,
			Collective:  NameLinker(event.Collective.Name),
			Public:      event.Public,
		}
		view.Events = append(view.Events, itemView)
	}
	return view
}

type CheckInDetails struct {
	Handle       NameLink
	Reasons      string
	EphemeralKey string
}

func EventDetailFromState(s *state.State, i *index.Index, hash crypto.Hash, token crypto.Token, ephemeral crypto.PrivateKey) *EventDetailView {
	event, ok := s.Events[hash]
	if !ok {
		event = s.Proposals.GetEvent(hash)
		if event == nil {
			return nil
		}
	}
	view := EventDetailView{
		Live:               event.Live,
		Description:        event.Description,
		StartAt:            event.StartAt,
		Collective:         NameLinker(event.Collective.Name),
		MemberOfCollective: event.Collective.IsMember(token),
		EstimatedEnd:       event.EstimatedEnd,
		Venue:              event.Venue,
		Open:               event.Open,
		Public:             event.Public,
		Checkedin:          make([]CheckInDetails, 0),
		ManagerMajority:    event.Managers.Majority,
		Managers:           make([]MemberDetailView, 0),
		Managing:           event.Managers.IsMember(token),
		Hash:               crypto.EncodeHash(hash),
	}
	if view.Managing {
		view.Head = HeaderInfo{
			Active:  "MyEvents",
			Path:    "venture / my events / ",
			EndPath: event.StartAt.Format("2006-01-02") + " by " + LimitStringSize(event.Collective.Name, maxStringSize),
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Events",
			Path:    "explore / events / ",
			EndPath: event.StartAt.Format("2006-01-02") + " by " + LimitStringSize(event.Collective.Name, maxStringSize),
			Section: "explore",
		}
	}
	for token, _ := range event.Managers.ListOfMembers() {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Managers = append(view.Managers, MemberDetailView{Handle: handle, Link: url.QueryEscape(handle)})
		}
	}
	for token, greet := range event.Checkin {
		if handle, ok := s.Members[crypto.Hasher(token[:])]; ok {
			if greet != nil && greet.Action != nil {
				view.Greeted = append(view.Greeted, MemberDetailView{Handle: handle, Link: url.QueryEscape(handle)})
				// if its me, de-crypt message
				if greet.Action.CheckedIn.Equal(token) {
					dhCipher := dh.ConsensusCipher(ephemeral, greet.Action.EphemeralToken)
					if secretKey, err := dhCipher.Open(greet.Action.SecretKey); err == nil {
						cipher := crypto.CipherFromKey(secretKey)
						if content, err := cipher.Open(greet.Action.PrivateContent); err == nil {
							view.MyGreeting = string(content)
						}
					}
				}
			} else {
				bytes, _ := greet.EphemeralKey.MarshalText()
				reasons := event.CheckinReasons[token]
				view.Checkedin = append(view.Checkedin, CheckInDetails{Handle: NameLinker(handle), EphemeralKey: string(bytes), Reasons: reasons})
			}
		}
	}
	println(view.Description)
	return &view
}

func EventUpdateDetailFromState(s *state.State, i *index.Index, hash crypto.Hash, token crypto.Token) *EventDetailView {
	event, ok := s.Events[hash]
	if !ok {
		event = s.Proposals.GetEvent(hash)
		if event == nil {
			return nil
		}
	}
	head := HeaderInfo{
		Active:  "MyEvents",
		Path:    "venture / my events / " + event.StartAt.Format("2006-01-02") + " by " + LimitStringSize(event.Collective.Name, maxStringSize) + " / ",
		EndPath: "update",
		Section: "venture",
	}
	view := EventDetailView{
		StartAt:         event.StartAt,
		Live:            event.Live,
		Description:     event.Description,
		Collective:      NameLinker(event.Collective.Name),
		EstimatedEnd:    event.EstimatedEnd,
		Venue:           event.Venue,
		Open:            event.Open,
		Public:          event.Public,
		ManagerMajority: event.Managers.Majority,
		Managers:        make([]MemberDetailView, 0),
		Managing:        event.Managers.IsMember(token),
		Hash:            crypto.EncodeHash(hash),
		Head:            head,
	}
	for token, _ := range event.Managers.ListOfMembers() {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Managers = append(view.Managers, MemberDetailView{Handle: handle, Link: url.QueryEscape(handle)})
		}
	}
	return &view
}

// Members template struct

type MembersView struct {
	Hash   string
	Handle string
	Link   string
}

type MembersListView struct {
	Members []MembersView
	Head    HeaderInfo
}

type MemberDetailView struct {
	Handle string
	Link   string
}

type MemberDetailViewPage struct {
	Detail MemberDetailView
	Head   HeaderInfo
}

func MembersFromState(state *state.State) MembersListView {
	head := HeaderInfo{
		Active:  "Members",
		Path:    "explore / ",
		EndPath: "members",
		Section: "explore",
	}
	view := MembersListView{
		Head:    head,
		Members: make([]MembersView, 0),
	}
	for hash, member := range state.Members {
		hashText, _ := hash.MarshalText()
		itemView := MembersView{
			Hash:   string(hashText),
			Handle: member,
			Link:   url.QueryEscape(member),
		}
		view.Members = append(view.Members, itemView)
	}
	return view
}

func MemberDetailFromState(state *state.State, handle string) *MemberDetailViewPage {
	_, ok := state.MembersIndex[handle]
	if !ok {
		return nil
	}
	detail := MemberDetailView{
		Handle: handle,
		Link:   url.QueryEscape(handle),
	}
	head := HeaderInfo{
		Active:  "Members",
		Path:    "explore / members / ",
		EndPath: LimitStringSize(handle, maxStringSize),
		Section: "explore",
	}
	view := MemberDetailViewPage{
		Detail: detail,
		Head:   head,
	}
	return &view
}

// Central Connections

type LastAction struct {
	Type        string
	Handle      string
	TimeOfInstr string
	// TimeOfInstr time.Time
}

type LastReference struct {
	Author      string
	TimeOfInstr string
	// TimeOfInstr time.Time
}

type CentralCollectives struct {
	Name     string
	Link     string
	NBoards  int
	NStamps  int
	NEvents  int
	LastSelf LastAction
	LastAny  LastAction
}

type CentralBoards struct {
	Name     string
	Link     string
	NPins    int
	NEditors int
	LastSelf LastAction
	LastAny  LastAction
}

type CentralEvents struct {
	Hash         string
	DateCol      string //data e horario mais nome do coletivo
	NCheckins    int
	NPenCheckins int
}

type CentralEdits struct {
	Title       string
	CreatedAt   time.Time
	NReferences int
	LastRef     LastReference
}

type ConnectionsListView struct {
	Head         HeaderInfo
	Collectives  []CentralCollectives
	NCollectives int
	Boards       []CentralBoards
	NBoards      int
	Events       []CentralEvents
	NEvents      int
	Edits        []CentralEdits
	NEdits       int
}

func ConnectionsFromState(state *state.State, indexer *index.Index, token crypto.Token, genesisTime time.Time) ConnectionsListView {
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / ",
		EndPath: "connections",
		Section: "venture",
	}
	view := ConnectionsListView{
		Head:        head,
		Collectives: make([]CentralCollectives, 0),
		Boards:      make([]CentralBoards, 0),
		Events:      make([]CentralEvents, 0),
		Edits:       make([]CentralEdits, 0),
	}

	// Check collectives user is a member of and get their info
	memberscol := indexer.CollectivesOnMember(token)
	for _, collectiveName := range memberscol {
		collective := state.Collectives[crypto.Hasher([]byte(collectiveName))]
		if collective == nil {
			continue
		}
		nboards := len(indexer.BoardsOnCollective(collective))
		nstamps := len(indexer.StampsOnCollective(collective))
		nevents := len(indexer.EventsOnCollective(collective))
		item := CentralCollectives{
			Name:    collective.Name,
			Link:    url.QueryEscape(collective.Name),
			NBoards: nboards,
			NStamps: nstamps,
			NEvents: nevents,
		}
		lastaction := indexer.LastMemberActionOnCollective(token, collective.Name)
		if lastaction != nil {
			actionTime := genesisTime.Add(time.Second * time.Duration(lastaction.Epoch))
			item.LastSelf = LastAction{
				Type:        lastaction.Description,
				Handle:      state.Members[crypto.HashToken(token)],
				TimeOfInstr: PrettyDuration(time.Since(actionTime)),
			}
		}

		recent := indexer.GetLastAction(crypto.Hasher([]byte(collectiveName)))
		if recent != nil {
			actionTime := genesisTime.Add(time.Second * time.Duration(recent.Epoch))
			item.LastAny = LastAction{
				Type:        recent.Description,
				Handle:      state.Members[crypto.HashToken(recent.Author)],
				TimeOfInstr: PrettyDuration(time.Since(actionTime)),
			}
		}
		view.Collectives = append(view.Collectives, item)
	}
	view.NCollectives = len(view.Collectives)

	// Check boards user is an editor at and get their info
	membersboard := indexer.BoardsOnMember(token)
	for _, board := range membersboard {
		hashedboard := crypto.Hasher([]byte(board))
		item := CentralBoards{
			Name:     board,
			Link:     url.QueryEscape(board),
			NPins:    len(state.Boards[hashedboard].Pinned),
			NEditors: len(state.Boards[hashedboard].Editors.ListOfMembers()),
		}
		lastSelfPin := indexer.LastManagerPinOnBoard(token, board)
		if lastSelfPin != nil {
			selfPinTime := genesisTime.Add(time.Second * time.Duration(lastSelfPin.Epoch))
			item.LastSelf = LastAction{
				Type:        lastSelfPin.Description,
				Handle:      state.Members[crypto.HashToken(token)],
				TimeOfInstr: PrettyDuration(time.Since(selfPinTime)),
			}
		}
		lastPin := indexer.LastPinOnBoard(token, board)
		if lastPin != nil {
			pinTime := genesisTime.Add(time.Second * time.Duration(lastPin.Epoch))
			item.LastAny = LastAction{
				Type:        lastPin.Description,
				Handle:      state.Members[crypto.HashToken(lastPin.Author)],
				TimeOfInstr: PrettyDuration(time.Since(pinTime)),
			}
		}
		view.Boards = append(view.Boards, item)
	}
	view.NBoards = len(view.Boards)

	// Check events user is a manager on and get their info
	membersevent := indexer.EventsOnMember(token)
	for _, eventhash := range membersevent {
		event := state.Events[eventhash]
		eventname := event.StartAt.Format("2006-01-02") + " by " + event.Collective.Name
		ncheckins := len(event.Checkin)
		ngreets := 0
		for _, greet := range event.Checkin {
			if greet != nil {
				ngreets += 1
			}
		}
		item := CentralEvents{
			Hash:         crypto.EncodeHash(eventhash),
			DateCol:      eventname,
			NCheckins:    ncheckins,
			NPenCheckins: ncheckins - ngreets,
		}
		view.Events = append(view.Events, item)
	}
	view.NEvents = len(view.Events)
	// Check edits user has proposed and get their info
	for _, edit := range state.Edits {
		if edit.Authors.IsMember(token) {
			item := CentralEdits{
				Title: edit.Draft.Title,
			}
			view.Edits = append(view.Edits, item)
		}
	}
	view.NEdits = len(view.Edits)
	return view
}
