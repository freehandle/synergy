package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
)

func FormToI(r *http.Request, field string) int {
	value, _ := strconv.Atoi(r.FormValue(field))
	return value
}

func FormToB(r *http.Request, field string) byte {
	value, _ := strconv.Atoi(r.FormValue(field))
	byteValue := byte(value)
	return byteValue
}

func FormToHash(r *http.Request, field string) crypto.Hash {
	hash := crypto.DecodeHash(r.FormValue(field))
	return hash
}

func FormToToken(r *http.Request, field string, handles map[string]crypto.Token) crypto.Token {
	token := handles[r.FormValue(field)]
	return token
}

func FormToEphemeralToken(r *http.Request, field string) crypto.Token {
	return crypto.TokenFromString(r.FormValue(field))
}

func FormToTokenArray(r *http.Request, field string, handles map[string]crypto.Token) []crypto.Token {
	h := strings.Split(r.FormValue(field), ",")
	tokens := make([]crypto.Token, 0)
	for _, handle := range h {
		if token, ok := handles[handle]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func FormToHashArray(r *http.Request, field string) []crypto.Hash {
	h := strings.Split(r.FormValue(field), ",")
	hashes := make([]crypto.Hash, 0)
	for _, caption := range h {
		var hash crypto.Hash
		if hash.UnmarshalText([]byte(caption)) == nil {
			hashes = append(hashes, hash)
		}
	}
	return hashes
}

func FormToStringArray(r *http.Request, field string) []string {
	words := strings.Split(r.FormValue(field), ",")
	return words
}

func FormToBool(r *http.Request, field string) bool {
	return r.FormValue(field) == "on"

}

func FormToPolicy(r *http.Request) Policy {
	return Policy{
		Majority:      FormToI(r, "policyMajority"),
		SuperMajority: FormToI(r, "policySupermajority"),
	}
}

func FormToTime(r *http.Request, field string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04", r.FormValue(field))
	return t
}

func BoardEditorForm(r *http.Request, handles map[string]crypto.Token, token crypto.Token) BoardEditor {
	action := BoardEditor{
		Action:  "BoardEditor",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
		Board:   r.FormValue("board"),
		Editor:  FormToToken(r, "editor", handles),
		Insert:  FormToBool(r, "insert"),
	}
	if action.Insert {
		action.Editor = token
	}
	return action
}

func CancelEventForm(r *http.Request) CancelEvent {
	action := CancelEvent{
		Action:  "CancelEvent",
		Reasons: r.FormValue("reasons"),
		ID:      FormToI(r, "id"),
		Hash:    FormToHash(r, "hash"),
	}
	return action
}

func CheckinEventForm(r *http.Request, ephemeralToken crypto.Token) CheckinEvent {
	action := CheckinEvent{
		Action:         "CheckinEvent",
		ID:             FormToI(r, "id"),
		EphemeralToken: ephemeralToken,
		Reasons:        r.FormValue("reasons"),
		EventHash:      FormToHash(r, "eventhash"),
	}
	return action
}

func CreateBoardForm(r *http.Request) CreateBoard {
	action := CreateBoard{
		Action:      "CreateBoard",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		OnBehalfOf:  r.FormValue("onBehalfOf"),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Keywords:    strings.Split(r.FormValue("keywords"), ","),
		PinMajority: FormToI(r, "pinMajority"),
	}
	return action
}

func CreateCollectiveForm(r *http.Request) CreateCollective {
	action := CreateCollective{
		Action:      "CreateCollective",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Policy:      FormToPolicy(r),
	}
	return action
}

func CreateEventForm(r *http.Request, handles map[string]crypto.Token, token crypto.Token) CreateEvent {
	action := CreateEvent{
		Action:          "CreateEvent",
		ID:              FormToI(r, "id"),
		Reasons:         r.FormValue("reasons"),
		OnBehalfOf:      r.FormValue("onBehalfOf"),
		StartAt:         FormToTime(r, "startAt"),
		EstimatedEnd:    FormToTime(r, "estimatedEnd"),
		Description:     r.FormValue("description"),
		Venue:           r.FormValue("venue"),
		Open:            FormToBool(r, "open"),
		Public:          FormToBool(r, "public"),
		ManagerMajority: FormToI(r, "managerMajority"),
	}
	if s := r.FormValue("managers"); s == "" {
		action.Managers = []crypto.Token{token}
	} else {
		action.Managers = FormToTokenArray(r, "managers", handles)
	}
	return action

}

func DraftForm(r *http.Request, handles map[string]crypto.Token, file []byte, ext string) Draft {
	action := Draft{
		Action:        "Draft",
		ID:            FormToI(r, "id"),
		Reasons:       r.FormValue("reasons"),
		OnBehalfOf:    r.FormValue("onBehalfOf"),
		CoAuthors:     FormToTokenArray(r, "coAuthors", handles),
		Title:         r.FormValue("title"),
		Description:   r.FormValue("description"),
		Keywords:      FormToStringArray(r, "keywords"),
		ContentType:   FileType(r.FormValue("fileName")),
		File:          file,
		PreviousDraft: FormToHash(r, "PreviousDraft"),
		References:    FormToHashArray(r, "references"),
	}
	return action
}

func EditForm(r *http.Request, handles map[string]crypto.Token, file []byte, ext string) Edit {
	action := Edit{
		Action:      "Edit",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		OnBehalfOf:  r.FormValue("onBehalfOf"),
		CoAuthors:   FormToTokenArray(r, "coAuthors", handles),
		EditedDraft: FormToHash(r, "editedDraft"),
		ContentType: FileType(r.FormValue("fileName")),
		File:        file,
	}
	return action
}

func GreetCheckinEventForm(r *http.Request, handles map[string]crypto.Token) MultiGreetCheckinEvent {
	fmt.Println(r.Form)
	action := MultiGreetCheckinEvent{
		Action:         "GreetCheckinEvent",
		ID:             FormToI(r, "id"),
		Reasons:        r.FormValue("reasons"),
		PrivateContent: r.FormValue("privateContent"),
		EventHash:      FormToHash(r, "eventhash"),
	}
	action.CheckedIn = make(map[crypto.Token]crypto.Token)
	for key, value := range r.Form {
		if strings.HasPrefix(key, "check_") && len(value) > 0 {
			ephemeral := crypto.TokenFromString(strings.ReplaceAll(key, "check_", ""))
			handle := value[0]
			fmt.Println(ephemeral, handle)
			handleUnescaped, _ := url.QueryUnescape(handle)
			if token, ok := handles[handleUnescaped]; ok {
				action.CheckedIn[token] = ephemeral
			}
		}
	}
	return action
}

func ImprintStampForm(r *http.Request) ImprintStamp {
	action := ImprintStamp{
		Action:     "ImprintStamp",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
		Hash:       FormToHash(r, "hash"),
	}
	return action
}

func PinForm(r *http.Request) Pin {
	action := Pin{
		Action:  "Pin",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
		Board:   r.FormValue("boardName"),
		Draft:   FormToHash(r, "draft"),
		Pin:     FormToBool(r, "pin"),
	}
	return action
}

func ReactForm(r *http.Request) React {
	action := React{
		Action:     "React",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
		Hash:       FormToHash(r, "hash"),
		Reaction:   byte(FormToI(r, "reaction")),
	}
	return action
}

func ReleaseDraftForm(r *http.Request) ReleaseDraft {
	action := ReleaseDraft{
		Action:      "ReleaseDraft",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		ContentHash: FormToHash(r, "contentHash"),
	}
	text, _ := json.Marshal(action)
	fmt.Println(string(text))
	return action
}

func RemoveMemberForm(r *http.Request, handles map[string]crypto.Token) RemoveMember {
	action := RemoveMember{
		Action:     "RemoveMember",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
		Member:     FormToToken(r, "member", handles),
	}
	return action
}

func RequestMembershipForm(r *http.Request) RequestMembership {
	action := RequestMembership{
		Action:     "RequestMembership",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		Collective: r.FormValue("collective"),
		Include:    FormToBool(r, "include"),
	}
	text, _ := json.Marshal(action)
	log.Println(string(text))
	return action
}

func UpdateBoardForm(r *http.Request) UpdateBoard {
	action := UpdateBoard{
		Action:  "UpdateBoard",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
	}
	if s := r.FormValue("board"); s != "" {
		action.Board, _ = url.QueryUnescape(s)
	}
	if s := r.FormValue("description"); s != "" {
		action.Description = &s
	}
	if s := r.FormValue("keywords"); s != "" {
		keywords := FormToStringArray(r, "keywords")
		action.Keywords = &keywords
	}
	if s := r.FormValue("pinMajority"); s != "" {
		majority := byte(FormToI(r, "pinMajorty"))
		action.PinMajority = &majority
	}
	text, _ := json.Marshal(action)
	log.Println(string(text))
	return action
}

func UpdateCollectiveForm(r *http.Request) UpdateCollective {
	action := UpdateCollective{
		Action:     "UpdateCollective",
		ID:         FormToI(r, "id"),
		Reasons:    r.FormValue("reasons"),
		OnBehalfOf: r.FormValue("onBehalfOf"),
	}

	if s := r.FormValue("description"); s != "" {
		action.Description = &s
	}
	if s := r.FormValue("majority"); s != "" {
		majority := FormToB(r, "majority")
		action.Majority = &majority
	}
	if s := r.FormValue("supermajority"); s != "" {
		supermajority := FormToB(r, "supermajority")
		action.SuperMajority = &supermajority
	}
	return action
}

func UpdateEventForm(r *http.Request, handles map[string]crypto.Token) UpdateEvent {
	action := UpdateEvent{
		Action:    "UpdateEvent",
		ID:        FormToI(r, "id"),
		Reasons:   r.FormValue("reasons"),
		EventHash: FormToHash(r, "eventHash"),
	}
	if s := r.FormValue("description"); s != "" {
		action.Description = &s
	}
	if s := r.FormValue("venue"); s != "" {
		action.Venue = &s
	}
	if s := r.FormValue("open"); s != "" {
		open := FormToBool(r, "open")
		action.Open = &open
	}
	if s := r.FormValue("public"); s != "" {
		public := FormToBool(r, "public")
		action.Public = &public
	}
	if s := r.FormValue("managerMajority"); s != "" {
		majority := FormToI(r, "managerMajority")
		action.ManagerMajority = &majority
	}
	if s := r.FormValue("managers"); s != "" {
		managers := FormToTokenArray(r, "managers", handles)
		action.Managers = &managers
	}
	return action
}

func VoteForm(r *http.Request) Vote {
	action := Vote{
		Action:  "Vote",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
		Hash:    FormToHash(r, "hash"),
		Approve: FormToBool(r, "approve"),
	}
	return action
}
