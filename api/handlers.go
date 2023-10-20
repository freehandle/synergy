package api

import (
	// "fmt"
	// "log"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/state"
)

type TemplateInfo struct {
	Head           HeaderInfo
	CollectiveName string
}

type StateView struct {
	State     *state.State
	Templates map[string]*template.Template
}

func (a *Attorney) MediaHandler(w http.ResponseWriter, r *http.Request) {
	hashtext := r.URL.Path
	hashtext = strings.Replace(hashtext, "/media/", "", 1)
	hash := crypto.DecodeHash(hashtext)

	file, ok := a.state.Media[hash]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("file not found"))
		return
	}
	title := hashtext
	var ext string
	if edit, ok := a.state.Edits[hash]; ok {
		ext = edit.EditType
	} else if draft, ok := a.state.Drafts[hash]; ok {
		ext = draft.DraftType
	} else if draft, ok := a.state.Proposals.Draft[hash]; ok {
		ext = draft.DraftType
	} else if edit, ok := a.state.Proposals.Edit[hash]; ok {
		ext = edit.EditType
	}
	name := fmt.Sprintf("%v", title, ext)
	//cd := mime.FormatMediaType("attachment", map[string]string{"filename": name})
	//w.Header().Set("Content-Disposition", cd)
	//w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, name, time.Now(), bytes.NewReader(file))
}

func (a *Attorney) ReloadTemplates(w http.ResponseWriter, r *http.Request) {
	a.templates = template.New("root")
	files := make([]string, len(templateFiles))
	for n, file := range templateFiles {
		files[n] = fmt.Sprintf("./api/templates/%v.html", file)
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	a.templates = t
	a.MainHandler(w, r)
}

func (a *Attorney) MainHandler(w http.ResponseWriter, r *http.Request) {
	if err := a.templates.ExecuteTemplate(w, "main.html", ""); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CreateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	head := HeaderInfo{
		Active:  "CreateCollective",
		Path:    "venture >",
		EndPath: "create collective",
		Section: "venture",
	}
	if err := a.templates.ExecuteTemplate(w, "createcollective.html", head); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) NewEditHandler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("draftHash"))
		fmt.Println(crypto.EncodeHash(hash))
		if view := NewEdit(a.state, hash); view != nil {
			if err := a.templates.ExecuteTemplate(w, "edit.html", view); err != nil {
				log.Println(err)
			}
			return
		}
	}
}

func (a *Attorney) NewDraft2Handler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("previousVersion"))
	}
	view := NewDraftVersion(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "newdraft2.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EditViewHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/editview/")
	view := EditDetailFromState(a.state, a.indexer, hash, a.author)
	if err := a.templates.ExecuteTemplate(w, "editview.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) NewDraftHandler(w http.ResponseWriter, r *http.Request) {
	var hash crypto.Hash
	if err := r.ParseForm(); err == nil {
		hash = crypto.DecodeHash(r.FormValue("previousVersion"))
	}
	view := NewDraftVersion(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "newdraft.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	view := BoardsFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "boards.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) BoardHandler(w http.ResponseWriter, r *http.Request) {
	boardName := r.URL.Path
	boardName = strings.Replace(boardName, "/board/", "", 1)
	view := BoardDetailFromState(a.state, boardName, a.author)
	if view == nil {
		w.Write([]byte("board not found"))
	} else if err := a.templates.ExecuteTemplate(w, "board.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CollectivesHandler(w http.ResponseWriter, r *http.Request) {
	view := ColletivesFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "collectives.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CollectiveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/collective/", "", 1)
	view := CollectiveDetailFromState(a.state, a.indexer, name, a.author)
	if err := a.templates.ExecuteTemplate(w, "collective.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) DraftsHandler(w http.ResponseWriter, r *http.Request) {
	view := DraftsFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "drafts.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) DraftHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/draft/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	view := DraftDetailFromState(a.state, a.indexer, hash, a.author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "draft.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EditsHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/edits/")
	view := EditsFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "edits.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EventsHandler(w http.ResponseWriter, r *http.Request) {
	view := EventsFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "events.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) EventHandler(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/event/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	view := EventDetailFromState(a.state, a.indexer, hash, a.author, a.ephemeralprv)
	if err := a.templates.ExecuteTemplate(w, "event.html", view); err != nil {
		log.Println(err)
	}
}

func getHash(path string, root string) crypto.Hash {
	path = strings.Replace(path, root, "", 1)
	hash := crypto.DecodeHash(path)
	return hash
}

func (a *Attorney) RequestMemberShipVoteHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/requestmembership/")
	view := RequestMembershipFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "requestmembershipvote.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VotesHandler(w http.ResponseWriter, r *http.Request) {
	view := VotesFromState(a.state, a.indexer, a.author)
	if err := a.templates.ExecuteTemplate(w, "votes.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) MembersHandler(w http.ResponseWriter, r *http.Request) {
	view := MembersFromState(a.state)
	if err := a.templates.ExecuteTemplate(w, "members.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) MemberHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	name = strings.Replace(name, "/member/", "", 1)
	view := MemberViewFromState(a.state, a.indexer, name)
	if err := a.templates.ExecuteTemplate(w, "member.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	collective := r.FormValue("collective")
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture > connections > collectives > " + collective + " > ",
		EndPath: "create board",
		Section: "venture",
	}
	info := TemplateInfo{
		Head:           head,
		CollectiveName: collective,
	}
	if err := a.templates.ExecuteTemplate(w, "createboard.html", info); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteCreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateboard/")
	view := PendingBoardFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "votecreateboard.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteCreateEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateevent/")
	view := PendingEventFromState(a.state, a.indexer, hash)
	if err := a.templates.ExecuteTemplate(w, "votecreateevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteCancelEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecancelevent/")
	view := CancelEventFromState(a.state, a.indexer, hash)
	if err := a.templates.ExecuteTemplate(w, "votecancelevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) UpdateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	collective := strings.Replace(r.URL.Path, "/updatecollective/", "", 1)
	view := CollectiveToUpdateFromState(a.state, collective)
	if err := a.templates.ExecuteTemplate(w, "updatecollective.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteUpdateCollectiveHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/voteupdatecollective/")
	view := CollectiveUpdateFromState(a.state, hash, a.author)
	if err := a.templates.ExecuteTemplate(w, "voteupdatecollective.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) UpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := strings.Replace(r.URL.Path, "/updateboard/", "", 1)
	view := BoardToUpdateFromState(a.state, board)
	if err := a.templates.ExecuteTemplate(w, "updateboard.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteUpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/votecreateboard/")
	view := BoardUpdateFromState(a.state, hash)
	if err := a.templates.ExecuteTemplate(w, "voteupdateboard.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/updateevent/")
	view := EventUpdateDetailFromState(a.state, a.indexer, hash, a.author)
	if err := a.templates.ExecuteTemplate(w, "updateevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	collective := r.FormValue("collective")
	head := HeaderInfo{
		Active:  "Connections",
		Path:    "venture > connections > collectives > " + collective + " > ",
		EndPath: "create event",
		Section: "venture",
	}
	info := TemplateInfo{
		Head:           head,
		CollectiveName: collective,
	}
	if err := a.templates.ExecuteTemplate(w, "createevent.html", info); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) VoteUpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/voteupdateevent/")
	view := EventUpdateFromState(a.state, hash, a.author)
	if err := a.templates.ExecuteTemplate(w, "voteupdateevent.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) ConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	view := ConnectionsFromState(a.state, a.indexer, a.author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "connections.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	view := UpdatesViewFromState(a.state, a.indexer, a.author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "updates.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) PendingActionsHandler(w http.ResponseWriter, r *http.Request) {
	view := PendingActionsFromState(a.state, a.indexer, a.author, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "pending.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) MyMediaHandler(w http.ResponseWriter, r *http.Request) {
	view := MyMediaFromState(a.state, a.indexer, a.author)
	if err := a.templates.ExecuteTemplate(w, "mymedia.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) MyEventsHandler(w http.ResponseWriter, r *http.Request) {
	view := MyEventsFromState(a.state, a.indexer, a.author)
	if err := a.templates.ExecuteTemplate(w, "myevents.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Attorney) NewsHandler(w http.ResponseWriter, r *http.Request) {
	view := NewActionsFromState(a.state, a.indexer, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "news.html", view); err != nil {
		log.Println(err)
	}

	// if news != nil {
	// 	fmt.Printf("new actions: %+v\n", *news)
	// } else {
	// 	fmt.Printf("no new actions\n")
	// }
}

func (a *Attorney) DetailedVoteHandler(w http.ResponseWriter, r *http.Request) {
	hash := getHash(r.URL.Path, "/detailedvote/")
	view := DetailedVoteFromState(a.state, a.indexer, hash, a.genesisTime)
	if err := a.templates.ExecuteTemplate(w, "detailedvote.html", view); err != nil {
		log.Println(err)
	}
}
