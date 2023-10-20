package api

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/crypto/dh"
	"github.com/lienkolabs/breeze/util"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/index"
	"github.com/lienkolabs/synergy/social/state"
)

var templateFiles []string = []string{
	"main",
	"boards", "board", "collectives", "collective", "draft", "drafts", "edits", "events",
	"event", "member", "members", "votes", "newdraft2", "edit",
	"createboard", "votecreateboard", "updateboard", "voteupdateboard", "updateevent",
	"updatecollective", "voteupdatecollective", "createevent", "voteupdateevent", "editview",
	"createcollective", "connections", "updates", "news", "pending", "mymedia", "myevents",
	"detailedvote", "votecreateevent", "votecancelevent", "login", "signin",
}

type Attorney struct {
	author       crypto.Token
	pk           crypto.PrivateKey
	ephemeralprv crypto.PrivateKey
	ephemeralpub crypto.Token
	wallet       crypto.PrivateKey
	pending      map[crypto.Hash]actions.Action
	epoch        uint64
	gateway      social.Gatewayer
	state        *state.State
	templates    *template.Template
	indexer      *index.Index
	genesisTime  time.Time
}

func NewAttorneyServer(pk crypto.PrivateKey, token crypto.Token, port int, gateway social.Gatewayer, indexer *index.Index) *Attorney {
	attorney := Attorney{
		author:  token,
		pk:      pk,
		wallet:  pk,
		pending: make(map[crypto.Hash]actions.Action),
		gateway: gateway,
		state:   gateway.State(),
		epoch:   0,
		indexer: indexer,
	}
	attorney.genesisTime = time.Date(2023, 9, 20, 0, 0, 0, 0, time.UTC)
	attorney.ephemeralprv, attorney.ephemeralpub = dh.NewEphemeralKey()
	blockEvent := gateway.Register()
	send := make(chan actions.Action)
	go func() {
		for {
			select {
			case epoch := <-blockEvent:
				attorney.epoch = epoch
			case action := <-send:
				gateway.Action(attorney.DressAction(action))
			}
		}
	}()

	go func() {
		attorney.templates = template.New("root")
		files := make([]string, len(templateFiles))
		for n, file := range templateFiles {
			files[n] = fmt.Sprintf("./api/templates/%v.html", file)
		}
		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}
		attorney.templates = t

		mux := http.NewServeMux()

		fs := http.FileServer(http.Dir("./api/static"))
		mux.Handle("/static/", http.StripPrefix("/static/", fs)) //

		mux.HandleFunc("/api", attorney.ApiHandler)
		mux.HandleFunc("/", attorney.MainHandler)
		mux.HandleFunc("/boards", attorney.BoardsHandler)
		mux.HandleFunc("/board/", attorney.BoardHandler)
		mux.HandleFunc("/collectives", attorney.CollectivesHandler)
		mux.HandleFunc("/collective/", attorney.CollectiveHandler)
		mux.HandleFunc("/drafts", attorney.DraftsHandler)
		mux.HandleFunc("/draft/", attorney.DraftHandler)
		mux.HandleFunc("/edits/", attorney.EditsHandler)
		mux.HandleFunc("/events", attorney.EventsHandler)
		mux.HandleFunc("/event/", attorney.EventHandler)
		mux.HandleFunc("/members", attorney.MembersHandler)
		mux.HandleFunc("/member/", attorney.MemberHandler)
		mux.HandleFunc("/votes/", attorney.VotesHandler)
		mux.HandleFunc("/requestmembership/", attorney.RequestMemberShipVoteHandler)
		mux.HandleFunc("/newdraft", attorney.NewDraft2Handler)
		mux.HandleFunc("/edit", attorney.NewEditHandler)
		mux.HandleFunc("/editview/", attorney.EditViewHandler)
		mux.HandleFunc("/media/", attorney.MediaHandler)
		mux.HandleFunc("/uploadfile", attorney.UploadHandler)
		mux.HandleFunc("/createboard", attorney.CreateBoardHandler)
		mux.HandleFunc("/votecreateboard/", attorney.VoteCreateBoardHandler)
		mux.HandleFunc("/updateboard/", attorney.UpdateBoardHandler)
		mux.HandleFunc("/voteupdateboard/", attorney.UpdateBoardHandler)
		mux.HandleFunc("/updatecollective/", attorney.UpdateCollectiveHandler)
		mux.HandleFunc("/voteupdatecollective/", attorney.VoteUpdateCollectiveHandler)
		mux.HandleFunc("/updateevent/", attorney.UpdateEventHandler)
		mux.HandleFunc("/votecancelevent/", attorney.VoteCancelEventHandler)
		mux.HandleFunc("/votecreateevent/", attorney.VoteCreateEventHandler)
		mux.HandleFunc("/createevent", attorney.CreateEventHandler)
		mux.HandleFunc("/voteupdateevent/", attorney.VoteUpdateEventHandler)
		mux.HandleFunc("/news", attorney.NewsHandler)
		mux.HandleFunc("/connections/", attorney.ConnectionsHandler)
		mux.HandleFunc("/updates", attorney.UpdatesHandler)
		mux.HandleFunc("/pending", attorney.PendingActionsHandler)
		mux.HandleFunc("/createcollective/", attorney.CreateCollectiveHandler)
		mux.HandleFunc("/mymedia", attorney.MyMediaHandler)
		mux.HandleFunc("/myevents", attorney.MyEventsHandler)
		mux.HandleFunc("/detailedvote/", attorney.DetailedVoteHandler)
		mux.HandleFunc("/reload", attorney.ReloadTemplates)
		// mux.HandleFunc("/member/votes", attorney.VotesHandler)

		srv := &http.Server{
			Addr:         fmt.Sprintf(":%v", port),
			Handler:      mux,
			WriteTimeout: 2 * time.Second,
		}
		srv.ListenAndServe()

	}()

	return nil
}

func (a *Attorney) Send(all []actions.Action) {
	for _, action := range all {
		dressed := a.DressAction(action)
		a.gateway.Action(dressed)
	}
}

// Dress a giving action with current epoch, attorney´s author
// attorneys signature, attorneys wallet and wallet signature
func (a *Attorney) DressAction(action actions.Action) []byte {
	bytes := action.Serialize()
	dress := []byte{0}
	util.PutUint64(a.epoch, &dress)
	util.PutToken(a.author, &dress)
	dress = append(dress, 0, 1, 0, 0) // axe void synergy
	dress = append(dress, bytes[8+crypto.TokenSize:]...)
	util.PutToken(a.pk.PublicKey(), &dress)
	signature := a.pk.Sign(dress)
	util.PutSignature(signature, &dress)
	util.PutToken(a.wallet.PublicKey(), &dress)
	util.PutUint64(0, &dress) // zero fee
	signature = a.wallet.Sign(dress)
	util.PutSignature(signature, &dress)
	return dress
}

func (a *Attorney) Confirmed(hash crypto.Hash) {
	delete(a.pending, hash)
}
