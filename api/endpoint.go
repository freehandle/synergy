package api

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/synergy/social/actions"
	"github.com/freehandle/synergy/social/index"
	"github.com/freehandle/synergy/social/state"
)

const maxStringSize = 50

type HeaderInfo struct {
	UserName   string
	UserHandle string
	Active     string
	Path       string
	EndPath    string
	Section    string
	Error      string
	ServerName string
}

type ServerName struct {
	Head       HeaderInfo
	ServerName string
}

// Drafts template struct

type DraftsView struct {
	Title       string
	Authors     []AuthorDetail
	Hash        string
	Description string
	Keywords    []string
	ServerName  string
}

type DraftsListView struct {
	Drafts     []DraftsView
	Head       HeaderInfo
	ServerName string
}

type AuthorDetail struct {
	Name       string
	Link       string
	Collective bool
}

type ReferenceDetail struct {
	Title  string
	Author string
	Date   string
}

// co-autor, stamp, pin, version, release
type DraftVoteAction struct {
	Kind       string
	OnBehalfOf string // collective or board editor
	Hash       string
}

type NameLink struct {
	Name string
	Link string
}

func NameLinker(name string) NameLink {
	return NameLink{
		Name: name,
		Link: url.QueryEscape(name),
	}
}

type DraftEditView struct {
	Date       string
	Authors    []AuthorDetail
	Hash       string
	ServerName string
}

type DraftDetailView struct {
	Title       string
	Date        string
	Description string
	Keywords    []string
	Hash        string
	//Content      string
	Authors      []AuthorDetail
	References   []ReferenceDetail
	PreviousHash string
	Pinned       []NameLink
	Edited       bool
	Released     bool
	Stamps       []NameLink
	Votes        []DraftVoteAction
	Policy       Policy
	Authorship   bool
	Head         HeaderInfo
	Content      string
	Edits        []DraftEditView
	ServerName   string
}

type EditDetailedView struct {
	DraftTitle string
	DraftHash  string
	Reasons    string
	Hash       string
	Authors    []AuthorDetail
	Votes      []DraftVoteAction
	Head       HeaderInfo
	ServerName string
}

func EditDetailFromState(s *state.State, i *index.Index, hash crypto.Hash, token crypto.Token) *EditDetailedView {
	edit, ok := s.Edits[hash]
	if !ok {
		edit, ok = s.Proposals.Edit[hash]
		if !ok {
			return nil
		}
	}
	head := HeaderInfo{
		Active:  "MyDrafts",
		Path:    "venture / my drafts / " + edit.Draft.Title + " / ",
		EndPath: "edits",
		Section: "venture",
	}
	view := EditDetailedView{
		DraftTitle: edit.Draft.Title,
		DraftHash:  crypto.EncodeHash(edit.Draft.DraftHash),
		Reasons:    edit.Reasons,
		Hash:       crypto.EncodeHash(edit.Edit),
		Authors:    AuthorList(edit.Authors, s),
		Votes:      make([]DraftVoteAction, 0),
		Head:       head,
	}
	pending := i.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			if s.Proposals.Kind(pendingHash) == state.EditProposal {
				vote := DraftVoteAction{
					Kind:       "Authorship",
					OnBehalfOf: s.Proposals.OnBehalfOf(pendingHash),
					Hash:       crypto.EncodeHash(pendingHash),
				}
				view.Votes = append(view.Votes, vote)
			}
		}
	}
	return &view
}

func membersToHandles(members map[crypto.Token]struct{}, state *state.State) []string {
	handles := make([]string, 0)
	for member, _ := range members {
		handle, ok := state.Members[crypto.Hasher(member[:])]
		if ok {
			handles = append(handles, handle)
		}
	}
	return handles
}

func hashesToString(hashes []crypto.Hash) []string {
	output := make([]string, 0)
	for _, hash := range hashes {
		text, err := hash.MarshalText()
		if err != nil {
			output = append(output, string(text))
		}
	}
	return output
}

func PinList(pin []*state.Board) []NameLink {
	list := make([]NameLink, 0)
	if len(pin) == 0 {
		return list
	}
	for _, p := range pin {
		list = append(list, NameLinker(p.Name))
	}
	return list
}

func StampList(stamps []*state.Stamp) []NameLink {
	list := make([]NameLink, 0)
	if len(stamps) == 0 {
		return list
	}
	for _, p := range stamps {
		list = append(list, NameLinker(p.Reputation.Name))
	}
	return list
}

func References(r []crypto.Hash, s *state.State, genesis time.Time) []ReferenceDetail {
	references := make([]ReferenceDetail, 0)
	for _, hash := range r {
		if draft, ok := s.Drafts[hash]; ok {
			date := genesis.Add(time.Duration(draft.Date) * time.Second)
			reference := ReferenceDetail{
				Title:  draft.Title,
				Author: authorsEtAll(draft.Authors, s),
				Date:   fmt.Sprintf("%v", date.Year()),
			}
			references = append(references, reference)
		}
	}
	return references
}

func authorsEtAll(c state.Consensual, s *state.State) string {
	authors := AuthorList(c, s)
	if len(authors) == 0 {
		return ""
	}
	N := len(authors)
	tail := ""
	if len(authors) > 3 {
		N = 3
		tail = " et al."
	}
	authorlist := make([]string, N)
	for n := 0; n < N; n++ {
		authorlist[n] = authors[n].Name
	}

	return fmt.Sprintf("%v%v", strings.Join(authorlist, ","), tail)
}

func AuthorList(c state.Consensual, s *state.State) []AuthorDetail {
	if c == nil {
		return []AuthorDetail{}
	}
	name := c.CollectiveName()
	if name != "" {
		author := AuthorDetail{
			Name:       name,
			Link:       url.QueryEscape(name),
			Collective: true,
		}
		return []AuthorDetail{author}
	}
	authors := make([]AuthorDetail, 0)
	for token, _ := range c.ListOfMembers() {
		if handle, ok := s.Members[crypto.Hasher(token[:])]; ok {
			authors = append(authors, AuthorDetail{Name: handle, Link: url.QueryEscape(handle)})
		}
	}
	return authors
}

func DraftsFromState(state *state.State) DraftsListView {
	head := HeaderInfo{
		Active:  "Drafts",
		Path:    "explore / ",
		EndPath: "drafts",
		Section: "explore",
	}
	view := DraftsListView{
		Head:   head,
		Drafts: make([]DraftsView, 0),
	}
	for _, draft := range state.Drafts {
		hash, _ := draft.DraftHash.MarshalText()
		itemView := DraftsView{
			Title:       draft.Title,
			Hash:        string(hash),
			Authors:     AuthorList(draft.Authors, state),
			Description: draft.Description,
			Keywords:    draft.Keywords,
		}
		view.Drafts = append(view.Drafts, itemView)
	}
	return view
}

func DraftDetailFromState(s *state.State, i *index.Index, hash crypto.Hash, token crypto.Token, genesis time.Time) *DraftDetailView {
	draft, ok := s.Drafts[hash]
	if !ok {
		draft, ok = s.Proposals.Draft[hash]
		if !ok {
			return nil
		}
	}
	date := genesis.Add(time.Duration(draft.Date) * time.Second)
	hashText, _ := hash.MarshalText()
	view := DraftDetailView{
		Title:       draft.Title,
		Date:        PrettyDate(date),
		Description: draft.Description,
		Keywords:    draft.Keywords,
		Authors:     AuthorList(draft.Authors, s),
		References:  References(draft.References, s, genesis),
		Pinned:      PinList(draft.Pinned),
		Votes:       make([]DraftVoteAction, 0),
		Authorship:  draft.Authors.IsMember(token),
		Hash:        string(hashText),
	}
	if view.Authorship {
		view.Head = HeaderInfo{
			Active:  "MyDrafts",
			Path:    "venture / my drafts / ",
			EndPath: LimitStringSize(draft.Title, maxStringSize),
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Drafts",
			Path:    "explore / drafts / ",
			EndPath: LimitStringSize(draft.Title, maxStringSize),
			Section: "explore",
		}
	}
	pending := i.GetVotes(token)
	if len(pending) > 0 {
		for pendingHash := range pending {
			pendingHashText, _ := pendingHash.MarshalText()
			vote := DraftVoteAction{
				OnBehalfOf: s.Proposals.OnBehalfOf(pendingHash),
				Hash:       string(pendingHashText),
			}
			switch s.Proposals.Kind(pendingHash) {
			case state.DraftProposal:
				if pending, ok := s.Proposals.Draft[pendingHash]; ok && pending.DraftHash.Equal(hash) {
					vote.Kind = "Authorship"
					view.Votes = append(view.Votes, vote)
				}
			case state.ReleaseDraftProposal:
				if pending, ok := s.Proposals.ReleaseDraft[pendingHash]; ok && pending.Draft.DraftHash.Equal(hash) {
					vote.Kind = "Release"
					view.Votes = append(view.Votes, vote)
				}
			case state.PinProposal:
				if pending, ok := s.Proposals.Pin[pendingHash]; ok && pending.Draft.DraftHash.Equal(hash) {
					vote.Kind = "Pin"
					view.Votes = append(view.Votes, vote)
				}
			case state.ImprintStampProposal:
				if pending, ok := s.Proposals.ImprintStamp[pendingHash]; ok && pending.Release.Draft.DraftHash.Equal(hash) {
					vote.Kind = "Stamp"
					view.Votes = append(view.Votes, vote)
				}
			}
		}
	}
	if release, ok := s.Releases[draft.DraftHash]; ok {
		view.Stamps = StampList(release.Stamps)
		view.Released = true
	}
	if len(draft.Edits) > 0 {
		view.Edited = true
	}
	if draft.PreviousVersion != nil {
		text, _ := draft.PreviousVersion.DraftHash.MarshalText()
		view.PreviousHash = string(text)
	}
	view.Policy.Majority, view.Policy.SuperMajority = draft.Authors.GetPolicy()

	if draft.DraftType == "txt" {
		if media, ok := s.Media[hash]; ok {
			view.Content = string(media)
		}
	} else if draft.DraftType == "md" {
		if media, ok := s.Media[hash]; ok {
			view.Content = mdToHTML(media)
		}
	}

	view.Edits = make([]DraftEditView, 0)
	for _, edit := range draft.Edits {
		date := genesis.Add(time.Duration(edit.Date) * time.Second)
		editView := DraftEditView{
			Date:    PrettyDate(date),
			Authors: AuthorList(edit.Authors, s),
			Hash:    crypto.EncodeHash(edit.Edit),
		}
		view.Edits = append(view.Edits, editView)
	}
	return &view
}

// Edits template struct

type EditsView struct {
	Authors    []AuthorDetail
	Reasons    string
	Hash       string
	ServerName string
}

type EditsListView struct {
	DraftTitle string
	DraftHash  string
	Edits      []EditsView
	Head       HeaderInfo
	ServerName string
}

func EditsFromState(s *state.State, drafthash crypto.Hash) EditsListView {
	draft, ok := s.Drafts[drafthash]
	if !ok {
		return EditsListView{}
	}
	head := HeaderInfo{
		Active:  "MyDrafts",
		Path:    "venture / my drafts / " + LimitStringSize(draft.Title, maxStringSize) + " / ",
		EndPath: "edits",
		Section: "venture",
	}
	view := EditsListView{
		DraftTitle: draft.Title,
		DraftHash:  crypto.EncodeHash(draft.DraftHash),
		Edits:      make([]EditsView, 0),
		Head:       head,
	}
	for _, edit := range draft.Edits {
		itemView := EditsView{
			Authors: AuthorList(edit.Authors, s),
			Reasons: edit.Reasons,
			Hash:    crypto.EncodeHash(edit.Edit),
		}
		view.Edits = append(view.Edits, itemView)
	}
	return view
}

// Votes template struct

type VotesView struct {
	Action            string
	Scope             string
	ScopeLink         string
	Hash              string
	Handler           string
	ObjectType        string
	ObjectLink        string
	ObjectCaption     string
	ComplementType    string
	ComplementLink    string
	ComplementCaption string
	Reasons           string
	ServerName        string
}

type VotesListView struct {
	Votes      []VotesView
	Head       HeaderInfo
	ServerName string
}

type VoteDetailView struct {
	Hash       string
	ServerName string
}

func VotesFromState(s *state.State, i *index.Index, token crypto.Token) VotesListView {
	head := HeaderInfo{
		Active:  "Votes",
		Path:    "venture / ",
		EndPath: "consensus votes",
		Section: "venture",
	}
	view := VotesListView{
		Head:  head,
		Votes: make([]VotesView, 0),
	}
	votes := i.GetVotes(token)
	for hash := range votes {
		hashText, _ := hash.MarshalText()
		itemView := VotesView{
			Action:    s.Proposals.KindText(hash),
			Scope:     s.Proposals.OnBehalfOf(hash),
			ScopeLink: url.QueryEscape(s.Proposals.OnBehalfOf(hash)),
			Hash:      string(hashText),
			Reasons:   i.Reason(hash),
		}
		switch s.Proposals.Kind(hash) {
		case state.RequestMembershipProposal:
			prop := s.Proposals.RequestMembership[hash]
			handle, ok := s.Members[crypto.Hasher(prop.Request.Author[:])]
			if ok {
				itemView.ObjectCaption = handle
				itemView.ObjectLink = fmt.Sprintf("member/%v", url.QueryEscape(handle))
				itemView.ObjectType = ""
			}
			itemView.ComplementType = "collective"
			itemView.ComplementCaption = prop.Collective.Name
			itemView.ComplementLink = fmt.Sprintf("collective/%v", url.QueryEscape(prop.Collective.Name))
		case state.DraftProposal:
			itemView.Handler = "draft"
		case state.PinProposal:
			itemView.Handler = "draft"
			prop := s.Proposals.Pin[hash]
			itemView.Hash = crypto.EncodeHash(prop.Draft.DraftHash)

		case state.ImprintStampProposal:
			itemView.Handler = "draft"
			prop := s.Proposals.ImprintStamp[hash]
			itemView.Hash = crypto.EncodeHash(prop.Release.Draft.DraftHash)

		case state.ReleaseDraftProposal:
			itemView.Handler = "draft"
			prop := s.Proposals.ReleaseDraft[hash]
			itemView.Hash = crypto.EncodeHash(prop.Draft.DraftHash)

		case state.UpdateCollectiveProposal:
			itemView.Handler = "voteupdatecollective"
			//prop := s.Proposals.UpdateCollective[hash]
			//itemView.Hash = crypto.EncodeHash(prop.Hash)

		case state.RemoveMemberProposal:
			prop := s.Proposals.RemoveMember[hash]
			handle, ok := s.Members[crypto.Hasher(prop.Remove.Member[:])]
			if ok {
				itemView.ObjectCaption = handle
				itemView.ObjectLink = fmt.Sprintf("member/%v", url.QueryEscape(handle))
				itemView.ObjectType = ""
			}
		case state.EditProposal:
			itemView.Handler = "editview"
		case state.CreateBoardProposal:
			itemView.Handler = "votecreateboard"
		case state.UpdateBoardProposal:
			itemView.Handler = "voteupdateboard"
		case state.BoardEditorProposal:
			prop := s.Proposals.BoardEditor[hash]
			editor, ok := s.Members[crypto.Hasher(prop.Editor[:])]
			if ok {
				itemView.ObjectCaption = editor
				itemView.ObjectLink = fmt.Sprintf("member/%v", url.QueryEscape(editor))
				if prop.Insert {
					itemView.ObjectType = "include"
				} else {
					itemView.ObjectType = "remove"
				}
			}
			itemView.Scope = ""
			itemView.ComplementCaption = prop.Board.Name
			itemView.ComplementType = "board"
			itemView.ComplementLink = fmt.Sprintf("board/%v", url.QueryEscape(prop.Board.Name))
		case state.ReactProposal:
		case state.CreateEventProposal:
			itemView.Handler = "votecreateevent"
		case state.CancelEventProposal:
			itemView.Handler = "votecancelevent"
			//prop := s.Proposals.CancelEvent[hash]
			//itemView.Hash = crypto.EncodeHash(prop.Event.Hash)
		case state.UpdateEventProposal:
			itemView.Handler = "voteupdateevent"

		}
		view.Votes = append(view.Votes, itemView)
	}
	return view
}

type RequestMembershipView struct {
	Collective string
	Handle     string
	Hash       string
	Reasons    string
	Majority   string
	ServerName string
}

func RequestMembershipFromState(s *state.State, hash crypto.Hash) *RequestMembershipView {
	vote, ok := s.Proposals.RequestMembership[hash]
	if !ok {
		return nil
	}
	token := vote.Request.Author
	handle, ok := s.Members[crypto.Hasher(token[:])]
	if !ok {
		return nil
	}
	hashText, _ := vote.Hash.MarshalText()
	majority, _ := vote.Collective.GetPolicy()
	return &RequestMembershipView{
		Collective: vote.Collective.Name,
		Handle:     handle,
		Hash:       string(hashText),
		Reasons:    vote.Request.Reasons,
		Majority:   fmt.Sprintf("%v", majority),
	}
}

type EditVersion struct {
	DraftHash  string
	Head       HeaderInfo
	ServerName string
}

func NewEdit(s *state.State, hash crypto.Hash) *EditVersion {
	draft, ok := s.Drafts[hash]
	if !ok {
		return nil
	}
	head := HeaderInfo{
		Active:  "Draft",
		Path:    "venture / drafts / " + LimitStringSize(draft.Title, maxStringSize) + " / venture / ",
		EndPath: "edit",
		Section: "venture",
	}
	return &EditVersion{
		DraftHash: crypto.EncodeHash(draft.DraftHash),
		Head:      head,
	}
}

type DraftVersion struct {
	OnBehalfOf    string
	Policy        Policy
	Title         string
	Keywords      string
	Description   string
	PreviousDraft string
	References    string
	Head          HeaderInfo
	ServerName    string
}

func NewDraftVersion(s *state.State, hash crypto.Hash) *DraftVersion {
	head := HeaderInfo{
		Active:  "NewDraft",
		Path:    "venture / ",
		EndPath: "new draft",
		Section: "venture",
	}
	draft, ok := s.Drafts[hash]
	if !ok {
		return &DraftVersion{
			Head: head,
		}
	}
	majority, supermajority := draft.Authors.GetPolicy()
	return &DraftVersion{
		OnBehalfOf:    draft.Authors.CollectiveName(),
		Policy:        Policy{Majority: majority, SuperMajority: supermajority},
		Title:         draft.Title,
		Keywords:      strings.Join(draft.Keywords, ","),
		Description:   draft.Description,
		PreviousDraft: crypto.EncodeHash(hash),
		Head:          head,
	}
}

// func VoteDetailFromState(state *state.State, hash string) *VoteDetailView {
// 	ok := state.Vote(hash)
// 	if !ok {
// 		return nil
// 	}
// 	view := VoteDetailView{
// 		Hash: hash,
// 	}
// 	return &view
// }

// Boards template struct

type CollectiveUpdateView struct {
	Name             string
	Link             string
	OldDescription   string
	Description      string
	OldMajority      int
	Majority         int
	OldSuperMajority int
	SuperMajority    int
	Member           bool
	Hash             string
	Reasons          string
	Head             HeaderInfo
	Voting           DetailedVoteView
	ServerName       string
}

func CollectiveToUpdateFromState(s *state.State, name string) *CollectiveUpdateView {
	collectiveName, _ := url.QueryUnescape(name)
	col, ok := s.Collective(collectiveName)
	if !ok {
		return nil
	}
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / connections / collectives " + LimitStringSize(name, maxStringSize) + " / ",
		EndPath: "update collective",
		Section: "venture",
	}
	update := &CollectiveUpdateView{
		Name:             collectiveName,
		Link:             url.QueryEscape(name),
		OldDescription:   col.Description,
		OldMajority:      col.Policy.Majority,
		OldSuperMajority: col.Policy.SuperMajority,
		Head:             head,
	}
	return update
}

func CollectiveUpdateFromState(s *state.State, hash crypto.Hash, token crypto.Token) *CollectiveUpdateView {
	pending, ok := s.Proposals.UpdateCollective[hash]
	if !ok {
		return nil
	}
	live := pending.Collective
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / connections / collectives / ",
		EndPath: "update collective " + LimitStringSize(live.Name, maxStringSize),
		Section: "venture",
	}
	update := &CollectiveUpdateView{
		Name:             live.Name,
		Link:             url.QueryEscape(live.Name),
		OldDescription:   live.Description,
		OldMajority:      live.Policy.Majority,
		OldSuperMajority: live.Policy.SuperMajority,
		Hash:             crypto.EncodeHash(hash),
		Reasons:          pending.Update.Reasons,
		Head:             head,
		Voting:           NewDetailedVoteView(pending.Votes, pending.Collective, s),
	}
	if pending.Update.Description != nil {
		update.Description = *pending.Update.Description
	}
	if pending.Update.Majority != nil {
		update.Majority = int(*pending.Update.Majority)
	}
	if pending.Update.SuperMajority != nil {
		update.SuperMajority = int(*pending.Update.SuperMajority)
	}
	if live.IsMember(token) {
		update.Member = true
	}

	return update
}

type BoardUpdateView struct {
	Name              string
	Link              string
	Collective        string
	Description       string
	OldDescription    string
	KeywordsString    string
	OldKeywordsString string
	PinMajority       byte
	OldPinMajority    byte
	Reasons           string
	Hash              string
	Head              HeaderInfo
	Voting            DetailedVoteView
	ServerName        string
	Editorship        bool
}

func BoardToUpdateFromState(s *state.State, name string) *BoardUpdateView {
	boardName, _ := url.QueryUnescape(name)
	live, ok := s.Board(boardName)
	if !ok {
		return nil
	}
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / connections / boards / ",
		EndPath: LimitStringSize(live.Name, maxStringSize),
		Section: "venture",
	}
	update := &BoardUpdateView{
		Name:           live.Name,
		Link:           url.QueryEscape(live.Name),
		Collective:     live.Collective.Name,
		OldDescription: live.Description,
		OldPinMajority: byte(live.Editors.Majority),
		Head:           head,
	}
	if len(live.Keyword) > 0 {
		update.OldKeywordsString = strings.Join(live.Keyword, ",")
	}
	return update
}

func BoardUpdateFromState(s *state.State, hash crypto.Hash, token crypto.Token) *BoardUpdateView {
	pending, ok := s.Proposals.UpdateBoard[hash]
	if !ok {
		return nil
	}
	live := pending.Board
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / connections / boards / ",
		EndPath: "update board " + LimitStringSize(live.Name, maxStringSize),
		Section: "venture",
	}
	update := &BoardUpdateView{
		Name:           live.Name,
		Link:           url.QueryEscape(live.Name),
		Collective:     live.Collective.Name,
		OldDescription: live.Description,
		OldPinMajority: byte(live.Editors.Majority),
		Reasons:        pending.Origin.Reasons,
		Hash:           crypto.EncodeHash(pending.Hash),
		Head:           head,
		Voting:         NewDetailedVoteView(pending.Votes, pending.Board.Collective, s),
	}
	if pending.Description != nil {
		update.Description = *pending.Description
	}
	if pending.PinMajority != nil {
		update.PinMajority = *pending.PinMajority
	}

	if len(live.Keyword) > 0 {
		update.OldKeywordsString = strings.Join(live.Keyword, ",")
	}
	if pending.Keywords != nil {
		update.KeywordsString = strings.Join(*pending.Keywords, ",")
	}
	if live.Editors.IsMember(token) {
		update.Editorship = true
	}
	return update
}

type BoardsView struct {
	Name           string
	Description    string
	Hash           string
	Collective     string
	CollectiveLink string
	Link           string
	Keywords       []string
	ServerName     string
}

type BoardsListView struct {
	Boards     []BoardsView
	Head       HeaderInfo
	ServerName string
}

type VoteDetails struct {
	Caption string
	Link    string
	Reasons string
}

type DetailedVoteView struct {
	Voted      int
	Approve    []VoteDetails
	Reject     []VoteDetails
	NotCast    []VoteDetails
	Majority   int
	ServerName string
}

func NewDetailedVoteView(votes []actions.Vote, consensus state.Consensual, s *state.State) DetailedVoteView {
	majority, _ := consensus.GetPolicy()
	view := DetailedVoteView{
		Approve:  make([]VoteDetails, 0),
		Reject:   make([]VoteDetails, 0),
		NotCast:  make([]VoteDetails, 0),
		Majority: majority,
	}
	allVoters := make(map[crypto.Token]struct{})
	for keyvoter, valuevoter := range consensus.ListOfMembers() {
		allVoters[keyvoter] = valuevoter
	}
	for _, vote := range votes {
		if !consensus.IsMember(vote.Author) {
			continue
		}
		delete(allVoters, vote.Author)
		handle := s.Members[crypto.HashToken(vote.Author)]
		voteDetail := VoteDetails{
			Caption: handle,
			Link:    url.QueryEscape(handle),
			Reasons: vote.Reasons,
		}
		if vote.Approve {
			view.Approve = append(view.Approve, voteDetail)
		} else {
			view.Reject = append(view.Reject, voteDetail)
		}
	}
	view.Voted = len(view.Approve) + len(view.Reject)
	for token := range allVoters {
		handle := s.Members[crypto.HashToken(token)]
		voteDetail := VoteDetails{
			Caption: handle,
			Link:    url.QueryEscape(handle),
		}
		view.NotCast = append(view.NotCast, voteDetail)

	}
	return view
}

type BoardDetailView struct {
	Name             string
	Link             string
	Description      string
	Collective       string
	CollectiveLink   string
	Keywords         []string
	PinMajority      int
	Editors          []MemberDetailView
	Drafts           []DraftsView
	Editorship       bool
	CollectiveMember bool
	Reasons          string
	Author           string
	Hash             string
	Head             HeaderInfo
	Voting           DetailedVoteView
	ServerName       string
}

func BoardsFromState(s *state.State) BoardsListView {
	head := HeaderInfo{
		Active:  "Boards",
		Path:    "explore / ",
		EndPath: "boards",
		Section: "explore",
	}
	view := BoardsListView{
		Head:   head,
		Boards: make([]BoardsView, 0),
	}
	for _, board := range s.Boards {
		itemView := BoardsView{
			Name:           board.Name,
			Description:    board.Description,
			Hash:           crypto.EncodeHash(crypto.Hasher([]byte(board.Name))),
			Collective:     board.Collective.Name,
			CollectiveLink: url.QueryEscape(board.Collective.Name),
			Link:           url.QueryEscape(board.Name),
			Keywords:       board.Keyword,
		}
		view.Boards = append(view.Boards, itemView)
	}
	return view
}

func PendingBoardFromState(s *state.State, hash crypto.Hash) *BoardDetailView {
	pending, ok := s.Proposals.CreateBoard[hash]
	if !ok {
		return nil
	}
	board := pending.Board
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture / connections / collectives / " + LimitStringSize(board.Collective.Name, maxStringSize) + " / ",
		EndPath: "create board " + LimitStringSize(board.Name, maxStringSize),
		Section: "venture",
	}
	view := BoardDetailView{
		Name:           board.Name,
		Link:           url.QueryEscape(board.Name),
		Description:    board.Description,
		Collective:     board.Collective.Name,
		CollectiveLink: url.QueryEscape(board.Collective.Name),
		Keywords:       board.Keyword,
		PinMajority:    board.Editors.Majority,
		Editors:        make([]MemberDetailView, 0), // isso aqui estaria errado, nao? acho que teria que ter o propositor como editor ja
		Drafts:         make([]DraftsView, 0),
		Reasons:        pending.Origin.Reasons,
		Hash:           crypto.EncodeHash(hash),
		Head:           head,
		Voting:         NewDetailedVoteView(pending.Votes, pending.Board.Collective, s),
	}
	view.Author = s.Members[crypto.Hasher(pending.Origin.Author[:])]
	return &view
}

func BoardDetailFromState(s *state.State, name string, token crypto.Token) *BoardDetailView {
	boardName, _ := url.QueryUnescape(name)
	board, ok := s.Board(boardName)
	if !ok {
		return nil
	}
	view := BoardDetailView{
		Name:             board.Name,
		Link:             url.QueryEscape(board.Name),
		Description:      board.Description,
		Collective:       board.Collective.Name,
		CollectiveLink:   url.QueryEscape(board.Collective.Name),
		Keywords:         board.Keyword,
		PinMajority:      board.Editors.Majority,
		Editors:          make([]MemberDetailView, 0),
		Drafts:           make([]DraftsView, 0),
		Editorship:       board.Editors.IsMember(token),
		CollectiveMember: board.Collective.IsMember(token) || board.Editors.IsMember(token),
		Hash:             crypto.EncodeHash(crypto.Hasher([]byte(board.Name))),
	}
	if view.Editorship {
		view.Head = HeaderInfo{
			Active:  "Connections",
			Path:    "venture / connections / boards / ",
			EndPath: LimitStringSize(board.Name, maxStringSize),
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Boards",
			Path:    "explore / boards / ",
			EndPath: LimitStringSize(board.Name, maxStringSize),
			Section: "explore",
		}
	}
	for token, _ := range board.Editors.Members {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Editors = append(view.Editors, MemberDetailView{Handle: handle, Link: url.QueryEscape(handle)})
		}
	}
	for _, d := range board.Pinned {
		draftView := DraftsView{
			Title:       d.Title,
			Authors:     make([]AuthorDetail, 0),
			Hash:        crypto.EncodeHash(d.DraftHash),
			Description: d.Description,
			Keywords:    d.Keywords,
		}
		view.Drafts = append(view.Drafts, draftView)
	}

	return &view
}

// Collectives template struct

type CollectivesView struct {
	Name         string
	Description  string
	Participants int
	Link         string
	ServerName   string
}

type CollectivesListView struct {
	Collectives []CollectivesView
	Head        HeaderInfo
	ServerName  string
}

type CaptionLink struct {
	Caption string
	Link    string
}

type StampView struct {
	Draft            CaptionLink
	DraftAuthors     []CaptionLink
	DraftDescription string
	DraftKeywords    []string
	ServerName       string
}

type BoardOnCollectiveView struct {
	Board       CaptionLink
	Description string
	Keywords    []string
	ServerName  string
}

type EventOnCollectiveView struct {
	StartAt     string
	Hash        string
	Venue       string
	Description string
	Managers    []CaptionLink
	ServerName  string
}

type CollectiveDetailView struct {
	Name          string
	Hash          string // hash of name for the reaction funcionalities
	Link          string
	Description   string
	Majority      int
	SuperMajority int
	Members       []MemberDetailView
	Membership    bool
	Head          HeaderInfo
	Stamps        []StampView
	Boards        []BoardOnCollectiveView
	Events        []EventOnCollectiveView
	ServerName    string
}

func CollectivesFromState(s *state.State) CollectivesListView {
	head := HeaderInfo{
		Active:  "Collectives",
		Path:    "explore / ",
		EndPath: "collectives",
		Section: "explore",
	}
	view := CollectivesListView{
		Head:        head,
		Collectives: make([]CollectivesView, 0),
	}
	for _, collective := range s.Collectives {
		itemView := CollectivesView{
			Name:         collective.Name,
			Description:  collective.Description,
			Participants: len(collective.Members),
			Link:         url.QueryEscape(collective.Name),
		}
		view.Collectives = append(view.Collectives, itemView)
	}
	return view
}

func CollectiveDetailFromState(s *state.State, i *index.Index, name string, token crypto.Token) *CollectiveDetailView {
	collectiveName, _ := url.QueryUnescape(name)
	collective, ok := s.Collective(collectiveName)
	if !ok {
		return nil
	}
	view := CollectiveDetailView{
		Name:          collective.Name,
		Hash:          crypto.EncodeHash(crypto.Hasher([]byte(collective.Name))),
		Link:          url.QueryEscape(collective.Name),
		Description:   collective.Description,
		Majority:      collective.Policy.Majority,
		SuperMajority: collective.Policy.SuperMajority,
		Members:       make([]MemberDetailView, 0),
		Membership:    collective.IsMember(token),
		Stamps:        make([]StampView, 0),
		Boards:        make([]BoardOnCollectiveView, 0),
		Events:        make([]EventOnCollectiveView, 0),
	}
	if view.Membership {
		view.Head = HeaderInfo{
			Active:  "Connections",
			Path:    "venture / connections / collectives / ",
			EndPath: LimitStringSize(collective.Name, maxStringSize),
			Section: "venture",
		}
	} else {
		view.Head = HeaderInfo{
			Active:  "Collectives",
			Path:    "explore / collectives / ",
			EndPath: LimitStringSize(collective.Name, maxStringSize),
			Section: "explore",
		}
	}
	for token, _ := range collective.Members {
		handle, ok := s.Members[crypto.Hasher(token[:])]
		if ok {
			view.Members = append(view.Members, MemberDetailView{Handle: handle, Link: url.QueryEscape(handle)})
		}
	}

	stamps := i.StampsOnCollective(collective)
	for _, stamp := range stamps {
		if stamp.Release != nil && stamp.Release.Draft != nil {
			draft := stamp.Release.Draft
			stampView := StampView{
				// Draft:            CaptionLink{Caption: draft.Title, Link: fmt.Sprintf("/draft/%v", crypto.EncodeHash(draft.DraftHash))},
				Draft:            CaptionLink{Caption: draft.Title, Link: url.QueryEscape(draft.DraftHash.String())},
				DraftAuthors:     make([]CaptionLink, 0),
				DraftDescription: draft.Description,
				DraftKeywords:    draft.Keywords,
			}
			for author, _ := range draft.Authors.ListOfMembers() {
				handle, ok := s.Members[crypto.HashToken(author)]
				if ok {
					// stampView.DraftAuthors = append(stampView.DraftAuthors, CaptionLink{Caption: handle, Link: fmt.Sprintf("/member/%v", handle)})
					stampView.DraftAuthors = append(stampView.DraftAuthors, CaptionLink{Caption: handle, Link: url.QueryEscape(handle)})
				}
			}
			view.Stamps = append(view.Stamps, stampView)
		}
	}

	boards := i.BoardsOnCollective(collective)
	for _, board := range boards {
		boardView := BoardOnCollectiveView{
			Board:       CaptionLink{Caption: board.Name, Link: url.QueryEscape(board.Name)},
			Description: board.Description,
			Keywords:    board.Keyword,
		}
		view.Boards = append(view.Boards, boardView)
	}

	events := i.EventsOnCollective(collective)
	for _, event := range events {
		eventView := EventOnCollectiveView{
			StartAt:     event.StartAt.Format("2006-01-02 15:04"),
			Hash:        crypto.EncodeHash(event.Hash),
			Venue:       event.Venue,
			Description: event.Description,
			Managers:    make([]CaptionLink, 0),
		}
		for manager, _ := range event.Managers.ListOfMembers() {
			handle, ok := s.Members[crypto.HashToken(manager)]
			if ok {
				eventView.Managers = append(eventView.Managers, CaptionLink{Caption: handle, Link: url.QueryEscape(handle)})
			}
		}
		view.Events = append(view.Events, eventView)
	}
	return &view
}
