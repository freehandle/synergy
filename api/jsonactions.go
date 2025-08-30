package api

import (
	"encoding/json"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/crypto/dh"
	"github.com/freehandle/synergy/social/actions"
)

/*
	actions
		BoardEditor
		CancelEvent
		CheckinEvent
		CreateBoard
		CreateCollective
		CreateEvent
		Draft
		Edit
		GreetCheckinEvent
		ImprintStamp
		Pin
		React
		ReleaseDraft
		RemoveMember
		RequestMembership
		UpdateBoard
		UpdateCollective
		UpdateEvent
		Vote
*/

type Policy struct {
	Majority      int `json:"majority"`
	SuperMajority int `json:"superMajority"`
}

type Action struct {
	Action string `json:"action"`
	ID     int    `json:"id"`
}

func JSONType(data []byte) string {
	var a Action
	json.Unmarshal(data, &a)
	return a.Action
}

type MultiGreetCheckinEvent struct {
	Action         string
	ID             int
	Reasons        string
	EventHash      crypto.Hash
	CheckedIn      map[crypto.Token]crypto.Token
	PrivateContent string
}

func (a MultiGreetCheckinEvent) ToAction() ([]actions.Action, error) {
	all := make([]actions.Action, 0)
	for token, ephemeral := range a.CheckedIn {
		action := actions.GreetCheckinEvent{
			Reasons:   a.Reasons,
			EventHash: a.EventHash,
			CheckedIn: token,
		}
		key := crypto.NewCipherKey()
		cipher := crypto.CipherFromKey(key)
		action.PrivateContent = cipher.Seal([]byte(a.PrivateContent))
		prv, pub := dh.NewEphemeralKey()
		//bytes, _ := a.EphemeralKey.MarshalText()
		dhCipher := dh.ConsensusCipher(prv, ephemeral)
		action.EphemeralToken = pub
		action.SecretKey = dhCipher.Seal(key)
		all = append(all, &action)
	}
	return all, nil
}

type GreetCheckinEvent struct {
	Action         string       `json:"action"`
	ID             int          `json:"id"`
	Reasons        string       `json:"reasons"`
	EventHash      crypto.Hash  `json:"eventHash"`
	CheckedIn      crypto.Token `json:"checkedIn"`
	EphemeralKey   crypto.Token
	PrivateContent string
}

func (a GreetCheckinEvent) ToAction() ([]actions.Action, error) {
	action := actions.GreetCheckinEvent{
		Reasons:   a.Reasons,
		EventHash: a.EventHash,
		CheckedIn: a.CheckedIn,
	}
	key := crypto.NewCipherKey()
	cipher := crypto.CipherFromKey(key)
	action.PrivateContent = cipher.Seal([]byte(a.PrivateContent))
	prv, pub := dh.NewEphemeralKey()
	//bytes, _ := a.EphemeralKey.MarshalText()
	dhCipher := dh.ConsensusCipher(prv, a.EphemeralKey)
	action.EphemeralToken = pub
	action.SecretKey = dhCipher.Seal(key)
	return []actions.Action{&action}, nil
}

type BoardEditor struct {
	Action  string       `json:"action"`
	ID      int          `json:"id"`
	Reasons string       `json:"reasons"`
	Board   string       `json:"board"`
	Editor  crypto.Token `json:"editor"`
	Insert  bool         `json:"insert"`
}

func (a BoardEditor) ToAction() ([]actions.Action, error) {
	action := actions.BoardEditor{
		Reasons: a.Reasons,
		Board:   a.Board,
		Editor:  a.Editor,
		Insert:  a.Insert,
	}
	return []actions.Action{&action}, nil
}

type CancelEvent struct {
	Action  string      `json:"action"`
	ID      int         `json:"id"`
	Reasons string      `json:"reasons"`
	Hash    crypto.Hash `json:"hash"`
}

func (a CancelEvent) ToAction() ([]actions.Action, error) {
	action := actions.CancelEvent{
		Reasons: a.Reasons,
		Hash:    a.Hash,
	}
	return []actions.Action{&action}, nil
}

type CheckinEvent struct {
	Action         string       `json:"action"`
	ID             int          `json:"id"`
	EphemeralToken crypto.Token `json:"ephemeralKey"`
	Reasons        string       `json:"reasons"`
	EventHash      crypto.Hash  `json:"eventHash"`
}

func (a CheckinEvent) ToAction() ([]actions.Action, error) {
	action := actions.CheckinEvent{
		EphemeralToken: a.EphemeralToken,
		Reasons:        a.Reasons,
		EventHash:      a.EventHash,
	}
	return []actions.Action{&action}, nil
}

type CreateBoard struct {
	Action      string   `json:"action"`
	ID          int      `json:"id"`
	Reasons     string   `json:"reasons"`
	OnBehalfOf  string   `json:"onBehalfOf"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	PinMajority int      `json:"pinMajority"`
}

func (a CreateBoard) ToAction() ([]actions.Action, error) {
	action := actions.CreateBoard{
		Reasons:     a.Reasons,
		OnBehalfOf:  a.OnBehalfOf,
		Name:        a.Name,
		Description: a.Description,
		Keywords:    a.Keywords,
		PinMajority: byte(a.PinMajority),
	}
	return []actions.Action{&action}, nil
}

type CreateCollective struct {
	Action      string `json:"action"`
	ID          int    `json:"id"`
	Reasons     string `json:"reasons"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Policy      Policy `json:"policy"`
}

func (a CreateCollective) ToAction() ([]actions.Action, error) {
	action := actions.CreateCollective{
		Reasons:     a.Reasons,
		Name:        a.Name,
		Description: a.Description,
		Policy:      actions.Policy(a.Policy),
	}
	return []actions.Action{&action}, nil
}

type CreateEvent struct {
	Action          string         `json:"action"`
	ID              int            `json:"id"`
	Reasons         string         `json:"reasons"`
	OnBehalfOf      string         `json:"onBehalfOf"`
	StartAt         time.Time      `json:"startAt"`
	EstimatedEnd    time.Time      `json:"estimatedEnd"`
	Description     string         `json:"description"`
	Venue           string         `json:"venue"`
	Open            bool           `json:"open"`
	Public          bool           `json:"public"`
	ManagerMajority int            `json:"managerMajority"`
	Managers        []crypto.Token `json:"managers,omitempty"`
}

func (a CreateEvent) ToAction() ([]actions.Action, error) {
	action := actions.CreateEvent{
		Reasons:         a.Reasons,
		OnBehalfOf:      a.OnBehalfOf,
		StartAt:         a.StartAt,
		EstimatedEnd:    a.EstimatedEnd,
		Description:     a.Description,
		Venue:           a.Venue,
		Open:            a.Open,
		Public:          a.Public,
		ManagerMajority: byte(a.ManagerMajority),
		Managers:        a.Managers,
	}
	return []actions.Action{&action}, nil
}

type Draft struct {
	Action        string         `json:"action"`
	ID            int            `json:"id"`
	Reasons       string         `json:"reasons"`
	OnBehalfOf    string         `json:"onBeahlfOf,omitempty"`
	CoAuthors     []crypto.Token `json:"coAuthors,omitempty"`
	Policy        *Policy        `json:"policy"`
	Title         string         `json:"title"`
	Keywords      []string       `json:"keywords"`
	Description   string         `json:"description"`
	ContentType   string         `json:"contentType"`
	File          []byte         `json:"filePath"`
	PreviousDraft crypto.Hash    `json:"previousDraft,omitempty"`
	References    []crypto.Hash  `json:"references,omitempty"`
}

func (a Draft) ToAction() ([]actions.Action, error) {
	truncated := splitBytes(a.File)
	allActions := make([]actions.Action, len(truncated.Parts))
	allActions[0] = &actions.Draft{
		Reasons:       a.Reasons,
		OnBehalfOf:    a.OnBehalfOf,
		CoAuthors:     a.CoAuthors,
		Policy:        (*actions.Policy)(a.Policy),
		Title:         a.Title,
		Keywords:      a.Keywords,
		Description:   a.Description,
		ContentType:   a.ContentType,
		ContentHash:   truncated.Hash,
		NumberOfParts: byte(len(truncated.Parts)),
		Content:       truncated.Parts[0],
		PreviousDraft: a.PreviousDraft,
		References:    a.References,
	}
	for n := 1; n < len(truncated.Parts); n++ {
		allActions[n] = &actions.MultipartMedia{
			Hash: truncated.Hash,
			Part: byte(n),
			Of:   byte(len(truncated.Parts)),
			Data: truncated.Parts[n],
		}
	}
	return allActions, nil
}

type Edit struct {
	Action      string         `json:"action"`
	ID          int            `json:"id"`
	Reasons     string         `json:"reasons"`
	OnBehalfOf  string         `json:"onBeahlfOf,omitempty"`
	CoAuthors   []crypto.Token `json:"coAuthors,omitempty"`
	EditedDraft crypto.Hash    `json:"editedDraft"`
	ContentType string         `json:"contentType"`
	File        []byte         `json:"filePath"`
}

func (a Edit) ToAction() ([]actions.Action, error) {
	truncated := splitBytes(a.File)
	allActions := make([]actions.Action, len(truncated.Parts))
	allActions[0] = &actions.Edit{
		Reasons:       a.Reasons,
		OnBehalfOf:    a.OnBehalfOf,
		CoAuthors:     a.CoAuthors,
		ContentType:   a.ContentType,
		EditedDraft:   a.EditedDraft,
		ContentHash:   truncated.Hash,
		NumberOfParts: byte(len(truncated.Parts)),
		Content:       truncated.Parts[0],
	}
	for n := 1; n < len(truncated.Parts); n++ {
		allActions[n] = &actions.MultipartMedia{
			Hash: truncated.Hash,
			Part: byte(n) + 1,
			Of:   byte(len(truncated.Parts)),
			Data: truncated.Parts[n],
		}
	}
	return allActions, nil
}

type ImprintStamp struct {
	Action     string      `json:"action"`
	ID         int         `json:"id"`
	Reasons    string      `json:"reasons"`
	OnBehalfOf string      `json:"onBeahlfOf,omitempty"`
	Hash       crypto.Hash `json:"hash"`
}

func (a ImprintStamp) ToAction() ([]actions.Action, error) {
	action := actions.ImprintStamp{
		Reasons:    a.Reasons,
		OnBehalfOf: a.OnBehalfOf,
		Hash:       a.Hash,
	}
	return []actions.Action{&action}, nil
}

type Pin struct {
	Action  string      `json:"action"`
	ID      int         `json:"id"`
	Reasons string      `json:"reasons"`
	Board   string      `json:"board"`
	Draft   crypto.Hash `json:"draft"`
	Pin     bool        `json:"pin"`
}

func (a Pin) ToAction() ([]actions.Action, error) {
	action := actions.Pin{
		Reasons: a.Reasons,
		Board:   a.Board,
		Pin:     a.Pin,
		Draft:   a.Draft,
	}
	return []actions.Action{&action}, nil
}

type React struct {
	Action     string      `json:"action"`
	ID         int         `json:"id"`
	Reasons    string      `json:"reasons"`
	OnBehalfOf string      `json:"onBeahlfOf,omitempty"`
	Hash       crypto.Hash `json:"hash"`
	Reaction   byte        `json:"reaction"`
}

func (a React) ToAction() ([]actions.Action, error) {
	action := actions.React{
		Reasons:    a.Reasons,
		OnBehalfOf: a.OnBehalfOf,
		Hash:       a.Hash,
		Reaction:   a.Reaction,
	}
	return []actions.Action{&action}, nil
}

type ReleaseDraft struct {
	Action      string      `json:"action"`
	ID          int         `json:"id"`
	Reasons     string      `json:"reasons"`
	ContentHash crypto.Hash `json:"contentHash"`
}

func (a ReleaseDraft) ToAction() ([]actions.Action, error) {
	action := actions.ReleaseDraft{
		Reasons:     a.Reasons,
		ContentHash: a.ContentHash,
	}
	return []actions.Action{&action}, nil
}

type RemoveMember struct {
	Action     string       `json:"action"`
	ID         int          `json:"id"`
	Reasons    string       `json:"reasons"`
	OnBehalfOf string       `json:"onBeahlfOf,omitempty"`
	Member     crypto.Token `json:"member"`
}

func (a RemoveMember) ToAction() ([]actions.Action, error) {
	action := actions.RemoveMember{
		Reasons:    a.Reasons,
		OnBehalfOf: a.OnBehalfOf,
		Member:     a.Member,
	}
	return []actions.Action{&action}, nil
}

type RequestMembership struct {
	Action     string `json:"action"`
	ID         int    `json:"id"`
	Reasons    string `json:"reasons"`
	Collective string `json:"collective"`
	Include    bool   `json:"include"`
}

func (a RequestMembership) ToAction() ([]actions.Action, error) {
	action := actions.RequestMembership{
		Reasons:    a.Reasons,
		Collective: a.Collective,
		Include:    a.Include,
	}
	return []actions.Action{&action}, nil
}

type UpdateBoard struct {
	Action      string    `json:"action"`
	ID          int       `json:"id"`
	Reasons     string    `json:"reasons"`
	Board       string    `json:"board"`
	Description *string   `json:"description,omitempty"`
	Keywords    *[]string `json:"keywords,omitempty"`
	PinMajority *byte     `json:"pinMajority"`
}

func (a UpdateBoard) ToAction() ([]actions.Action, error) {
	action := actions.UpdateBoard{
		Reasons:     a.Reasons,
		Board:       a.Board,
		Description: a.Description,
		Keywords:    a.Keywords,
		PinMajority: a.PinMajority,
	}
	return []actions.Action{&action}, nil
}

type UpdateCollective struct {
	Action        string  `json:"action"`
	ID            int     `json:"id"`
	Reasons       string  `json:"reasons"`
	OnBehalfOf    string  `json:"onBehalfOf"`
	Description   *string `json:"description,omitempty"`
	Majority      *byte   `json:"majority,omitempty"`
	SuperMajority *byte   `json:"superMajority,omitempty"`
}

func (a UpdateCollective) ToAction() ([]actions.Action, error) {
	action := actions.UpdateCollective{
		Reasons:       a.Reasons,
		OnBehalfOf:    a.OnBehalfOf,
		Description:   a.Description,
		Majority:      a.Majority,
		SuperMajority: a.SuperMajority,
	}
	return []actions.Action{&action}, nil
}

type UpdateEvent struct {
	Action          string          `json:"action"`
	ID              int             `json:"id"`
	Reasons         string          `json:"reasons"`
	EventHash       crypto.Hash     `json:"eventHash"`
	Description     *string         `json:"description,omitempty"`
	Venue           *string         `json:"venue,omitempty"`
	Open            *bool           `json:"open,omitempty"`
	Public          *bool           `json:"public,omitempty"`
	ManagerMajority *byte           `json:"managerMajority,omitempty"`
	Managers        *[]crypto.Token `json:"managers,omitempty"`
}

func (a UpdateEvent) ToAction() ([]actions.Action, error) {
	// var byteMajority *byte
	// if a.ManagerMajority != nil {
	// 	*byteMajority = byte(*a.ManagerMajority)
	// }
	action := actions.UpdateEvent{
		Reasons:     a.Reasons,
		EventHash:   a.EventHash,
		Description: a.Description,
		Venue:       a.Venue,
		Open:        a.Open,
		Public:      a.Public,
		// ManagerMajority: byteMajority,
		ManagerMajority: a.ManagerMajority,
		Managers:        a.Managers,
	}
	return []actions.Action{&action}, nil
}

type Vote struct {
	Action  string      `json:"action"`
	ID      int         `json:"id"`
	Reasons string      `json:"reasons,omitempty"`
	Hash    crypto.Hash `json:"hash"`
	Approve bool        `json:"approve"`
}

func (a Vote) ToAction() ([]actions.Action, error) {
	action := actions.Vote{
		Reasons: a.Reasons,
		Hash:    a.Hash,
		Approve: a.Approve,
	}
	return []actions.Action{&action}, nil
}
