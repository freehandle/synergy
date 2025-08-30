package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
)

func FormToI(r *http.Request, field string) int {
	if r == nil {
		log.Print("PANIC BUG: FormToI called with nil request ")
		return 0
	}
	value, _ := strconv.Atoi(r.FormValue(field))
	return value
}

func FormToB(r *http.Request, field string) byte {
	if r == nil {
		log.Print("PANIC BUG: FormToB called with nil request ")
		return 0
	}
	value, _ := strconv.Atoi(r.FormValue(field))
	byteValue := byte(value)
	return byteValue
}

func FormToHash(r *http.Request, field string) crypto.Hash {
	if r == nil {
		log.Print("PANIC BUG: FormToHash called with nil request ")
		return crypto.ZeroHash
	}
	hash := crypto.DecodeHash(r.FormValue(field))
	return hash
}

func FormToToken(r *http.Request, field string, handles map[string]crypto.Token) crypto.Token {
	if r == nil {
		log.Print("PANIC BUG: FormToToken called with nil request ")
		return crypto.ZeroToken
	}
	if handles == nil {
		log.Print("PANIC BUG: FormToToken called with nil handles ")
		return crypto.ZeroToken
	}
	token := handles[r.FormValue(field)]
	return token
}

func FormToEphemeralToken(r *http.Request, field string) crypto.Token {
	if r == nil {
		log.Print("PANIC BUG: FormToEphemeralToken called with nil request ")
		return crypto.ZeroToken
	}
	return crypto.TokenFromString(r.FormValue(field))
}

func FormToTokenArray(r *http.Request, field string, handles map[string]crypto.Token) []crypto.Token {
	if r == nil {
		log.Print("PANIC BUG: FormToTokenArray called with nil request ")
		return nil
	}
	if handles == nil {
		log.Print("PANIC BUG: FormToTokenArray called with nil handles ")
		return nil
	}
	h := strings.Split(r.FormValue(field), ",")
	tokens := make([]crypto.Token, 0)
	for _, handle := range h {
		handle = strings.TrimSpace(handle)
		if token, ok := handles[handle]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func FormToHashArray(r *http.Request, field string) []crypto.Hash {
	if r == nil {
		log.Print("PANIC BUG: FormToHashArray called with nil request ")
		return nil
	}
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
	if r == nil {
		log.Print("PANIC BUG: FormToStringArray called with nil request ")
		return nil
	}
	words := strings.Split(r.FormValue(field), ",")
	wordst := make([]string, 0)
	for _, w := range words {
		if strings.TrimSpace(w) != "" {
			wordst = append(wordst, w)
		}
	}
	return wordst
}

func FormToBool(r *http.Request, field string) bool {
	if r == nil {
		log.Print("PANIC BUG: FormToBool called with nil request ")
		return false
	}
	return r.FormValue(field) == "on"
}

func FormToPolicy(r *http.Request) Policy {
	if r == nil {
		log.Print("PANIC BUG: FormToPolicy called with nil request ")
		return Policy{}
	}
	return Policy{
		Majority:      FormToI(r, "policyMajority"),
		SuperMajority: FormToI(r, "policySupermajority"),
	}
}

func FormToTime(r *http.Request, field string) time.Time {
	if r == nil {
		log.Print("PANIC BUG: FormToTime called with nil request ")
		return time.Time{}
	}
	t, _ := time.Parse("2006-01-02T15:04", r.FormValue(field))
	return t
}

func BoardEditorForm(r *http.Request, handles map[string]crypto.Token, token crypto.Token) BoardEditor {
	if r == nil {
		log.Print("PANIC BUG: BoardEditorForm called with nil request ")
		return BoardEditor{}
	}
	if handles == nil {
		log.Print("PANIC BUG: BoardEditorForm called with nil handles ")
		return BoardEditor{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: CancelEventForm called with nil request ")
		return CancelEvent{}
	}
	action := CancelEvent{
		Action:  "CancelEvent",
		Reasons: r.FormValue("reasons"),
		ID:      FormToI(r, "id"),
		Hash:    FormToHash(r, "hash"),
	}
	return action
}

func CheckinEventForm(r *http.Request, ephemeralToken crypto.Token) CheckinEvent {
	if r == nil {
		log.Print("PANIC BUG: CheckinEventForm called with nil request ")
		return CheckinEvent{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: CreateBoardForm called with nil request ")
		return CreateBoard{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: CreateCollectiveForm called with nil request ")
		return CreateCollective{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: CreateEventForm called with nil request ")
		return CreateEvent{}
	}
	if handles == nil {
		log.Print("PANIC BUG: CreateEventForm called with nil handles ")
		return CreateEvent{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: DraftForm called with nil request ")
		return Draft{}
	}
	if handles == nil {
		log.Print("PANIC BUG: DraftForm called with nil handles ")
		return Draft{}
	}
	if file == nil {
		log.Print("PANIC BUG: DraftForm called with nil file ")
		return Draft{}
	}
	action := Draft{
		Action:        "Draft",
		ID:            FormToI(r, "id"),
		Reasons:       r.FormValue("reasons"),
		Title:         r.FormValue("title"),
		Description:   r.FormValue("description"),
		Keywords:      FormToStringArray(r, "keywords"),
		ContentType:   FileType(r.FormValue("fileName")),
		File:          file,
		PreviousDraft: FormToHash(r, "previousDraft"),
		References:    FormToHashArray(r, "references"),
	}
	if r.FormValue("onBehalfOf") != "" {
		action.OnBehalfOf = r.FormValue("onBehalfOf")
		return action

	} else if r.FormValue("coAuthors") != "" {
		action.CoAuthors = FormToTokenArray(r, "coAuthors", handles)
		return action
	}
	return action
}

func EditForm(r *http.Request, handles map[string]crypto.Token, file []byte, ext string) Edit {
	if r == nil {
		log.Print("PANIC BUG: EditForm called with nil request ")
		return Edit{}
	}
	if handles == nil {
		log.Print("PANIC BUG: EditForm called with nil handles ")
		return Edit{}
	}
	if file == nil {
		log.Print("PANIC BUG: EditForm called with nil file ")
		return Edit{}
	}
	action := Edit{
		Action:      "Edit",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		EditedDraft: FormToHash(r, "editedDraft"),
		ContentType: FileType(r.FormValue("fileName")),
		File:        file,
	}
	if r.FormValue("onBehalfOf") != "" {
		action.OnBehalfOf = r.FormValue("onBehalfOf")
		return action

	} else if r.FormValue("coAuthors") != "" {
		action.CoAuthors = FormToTokenArray(r, "coAuthors", handles)
		return action
	}
	return action
}

func GreetCheckinEventForm(r *http.Request, handles map[string]crypto.Token) MultiGreetCheckinEvent {
	if r == nil {
		log.Print("PANIC BUG: GreetCheckinEventForm called with nil request ")
		return MultiGreetCheckinEvent{}
	}
	if handles == nil {
		log.Print("PANIC BUG: GreetCheckinEventForm called with nil handles ")
		return MultiGreetCheckinEvent{}
	}
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
			handleUnescaped, _ := url.QueryUnescape(handle)
			if token, ok := handles[handleUnescaped]; ok {
				action.CheckedIn[token] = ephemeral
			}
		}
	}
	return action
}

func ImprintStampForm(r *http.Request) ImprintStamp {
	if r == nil {
		log.Print("PANIC BUG: ImprintStampForm called with nil request ")
		return ImprintStamp{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: PinForm called with nil request ")
		return Pin{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: ReactForm called with nil request ")
		return React{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: ReleaseDraftForm called with nil request ")
		return ReleaseDraft{}
	}
	action := ReleaseDraft{
		Action:      "ReleaseDraft",
		ID:          FormToI(r, "id"),
		Reasons:     r.FormValue("reasons"),
		ContentHash: FormToHash(r, "contentHash"),
	}
	return action
}

func RemoveMemberForm(r *http.Request, handles map[string]crypto.Token) RemoveMember {
	if r == nil {
		log.Print("PANIC BUG: RemoveMemberForm called with nil request ")
		return RemoveMember{}
	}
	if handles == nil {
		log.Print("PANIC BUG: RemoveMemberForm called with nil handles ")
		return RemoveMember{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: RequestMembershipForm called with nil request ")
		return RequestMembership{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: UpdateBoardForm called with nil request ")
		return UpdateBoard{}
	}
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
		majority := FormToB(r, "pinMajority")
		action.PinMajority = &majority
	}
	text, _ := json.Marshal(action)
	log.Println(string(text))
	return action
}

func UpdateCollectiveForm(r *http.Request) UpdateCollective {
	if r == nil {
		log.Print("PANIC BUG: UpdateCollectiveForm called with nil request ")
		return UpdateCollective{}
	}
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
	if r == nil {
		log.Print("PANIC BUG: UpdateEventForm called with nil request ")
		return UpdateEvent{}
	}
	if handles == nil {
		log.Print("PANIC BUG: UpdateEventForm called with nil handles ")
		return UpdateEvent{}
	}
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
		majority := FormToB(r, "managerMajority")
		action.ManagerMajority = &majority
	}
	if s := r.FormValue("managers"); s != "" {
		managers := FormToTokenArray(r, "managers", handles)
		action.Managers = &managers
	}
	text, _ := json.Marshal(action)
	log.Println(string(text))
	return action
}

func VoteForm(r *http.Request) Vote {
	if r == nil {
		log.Print("PANIC BUG: VoteForm called with nil request ")
		return Vote{}
	}
	action := Vote{
		Action:  "Vote",
		ID:      FormToI(r, "id"),
		Reasons: r.FormValue("reasons"),
		Hash:    FormToHash(r, "hash"),
		Approve: FormToBool(r, "approve"),
	}
	return action
}
