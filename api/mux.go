package api

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/freehandle/breeze/crypto"

	"github.com/freehandle/cb/vault"
	"github.com/freehandle/synergy/social/actions"
	"github.com/freehandle/synergy/social/index"
	"github.com/freehandle/synergy/social/state"
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

type ServerConfig struct {
	Vault       *vault.SecureVault
	Attorney    crypto.Token
	Ephemeral   crypto.Token
	Passwords   PasswordManager
	CookieStore *CookieStore
	//Gateway       social.Gatewayer
	Indexer       *index.Index
	Gateway       chan []byte
	State         *state.State
	GenesisTime   time.Time
	EmailPassword string
	Port          int
	Path          string
}

//type AuthorAction struct {
//	author crypto.Token
//	action actions.Action
//}

func NewGeneralAttorneyServer(config ServerConfig) (*AttorneyGeneral, chan error) {
	finalize := make(chan error, 2)

	attorneySecret, ok := config.Vault.Secrets[config.Attorney]
	if !ok {
		finalize <- fmt.Errorf("attorney secret key not found in vault")
		return nil, finalize
	}
	ephemeralSecret, ok := config.Vault.Secrets[config.Ephemeral]
	if !ok {
		finalize <- fmt.Errorf("ephemeral secret key not found in vault")
		return nil, finalize
	}

	attorney := AttorneyGeneral{
		//epoch:       config.State.Epoch, TODO: epoch get out of struct
		pk:            attorneySecret,
		Token:         config.Attorney,
		wallet:        attorneySecret,
		pending:       make(map[crypto.Hash]actions.Action),
		gateway:       config.Gateway,
		state:         config.State,
		indexer:       config.Indexer,
		session:       config.CookieStore,
		emailPassword: config.EmailPassword,
		//sessionend:   make(map[uint64][]string),
		genesisTime:  config.GenesisTime,
		ephemeralpub: config.Ephemeral,
		ephemeralprv: ephemeralSecret,
	}
	if config.Path == "" {
		config.Path = "./"
	}
	templatesPath := fmt.Sprintf("%v/api/templates", config.Path)
	attorney.signin = NewSigninManager(config.Passwords, config.EmailPassword, &attorney)
	attorney.templates = template.New("root")
	files := make([]string, len(templateFiles))

	for n, file := range templateFiles {
		files[n] = fmt.Sprintf("%v/%v.html", templatesPath, file)
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	attorney.templates = t

	//blockEvent := config.Gateway.Register()
	//send := make(chan *AuthorAction) // cria canal

	/*go func() {
		for {
			select {
			case attorney.epoch = <-blockEvent:
			case action := <-send: // usado aqui mas quem sabe dele???
				dressed := attorney.DressAction(action.action, action.author)
				attorney.gateway <- dressed
				//config.Gateway.Action()
			}
		}
	}()
	*/

	staticPath := fmt.Sprintf("%v/api/static/", config.Path)
	go NewServer(&attorney, config.Port, staticPath, finalize)

	return &attorney, finalize
}

func NewServer(attorney *AttorneyGeneral, port int, staticPath string, finalize chan error) {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(staticPath))
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
	mux.HandleFunc("/login", attorney.LoginHandler)
	mux.HandleFunc("/signin", attorney.SigninHandler)
	mux.HandleFunc("/signout", attorney.SignoutHandler)
	mux.HandleFunc("/credentials", attorney.CredentialsHandler)
	mux.HandleFunc("/newuser", attorney.NewUserHandler)
	// mux.HandleFunc("/member/votes", attorney.VotesHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
	}
	finalize <- srv.ListenAndServe()
}
