package api

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/freehandle/breeze/crypto"

	"github.com/freehandle/synergy/config"
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
	"detailedvote", "votecreateevent", "votecancelevent", "login", "signin", "totalsignin",
	"forgot", "reset", "resetpassword", "invite",
}

type ServerConfig struct {
	Vault       *config.SecretsVault
	Attorney    crypto.Token
	Ephemeral   crypto.Token
	Passwords   PasswordManager
	CookieStore *CookieStore
	//Gateway       social.Gatewayer
	Indexer     *index.Index
	Gateway     chan []byte
	State       *state.State
	GenesisTime time.Time
	Mail        Mailer
	Port        int
	Path        string
	ServerName  string
	Hostname    string
	Safe        int // optional link to safe for direct onbboarding
}

//type AuthorAction struct {
//	author crypto.Token
//	action actions.Action
//}

func NewGeneralAttorneyServer(cfg ServerConfig) (*AttorneyGeneral, chan error) {
	finalize := make(chan error, 2)

	attorneySecret, ok := cfg.Vault.Secrets[cfg.Attorney]
	if !ok {
		finalize <- fmt.Errorf("attorney secret key not found in vault")
		return nil, finalize
	}
	ephemeralSecret, ok := cfg.Vault.Secrets[cfg.Ephemeral]
	if !ok {
		finalize <- fmt.Errorf("ephemeral secret key not found in vault")
		return nil, finalize
	}

	attorney := AttorneyGeneral{
		//epoch:       config.State.Epoch, TODO: epoch get out of struct
		pk:      attorneySecret,
		Token:   cfg.Attorney,
		wallet:  attorneySecret,
		pending: make(map[crypto.Hash]actions.Action),
		gateway: cfg.Gateway,
		state:   cfg.State,
		indexer: cfg.Indexer,
		session: cfg.CookieStore,
		//Mail:    cfg.Mail,
		//sessionend:   make(map[uint64][]string),
		genesisTime:  cfg.GenesisTime,
		ephemeralpub: cfg.Ephemeral,
		ephemeralprv: ephemeralSecret,
		serverName:   cfg.ServerName,
		hostname:     cfg.Hostname,
		safe:         cfg.Safe,
		inviteUser:   make(map[crypto.Hash]struct{}),
	}
	if cfg.Path == "" {
		cfg.Path = "./"
	}
	templatesPath := fmt.Sprintf("%v/api/templates", cfg.Path)
	attorney.signin = NewSigninManager(cfg.Passwords, cfg.Mail, &attorney)
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

	staticPath := fmt.Sprintf("%v/api/static/", cfg.Path)
	go NewServer(&attorney, cfg.Port, staticPath, finalize, cfg.ServerName)

	return &attorney, finalize
}

func NewServer(attorney *AttorneyGeneral, port int, staticPath string, finalize chan error, servername string) {

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
	// mux.HandleFunc("/votes/", attorney.VotesHandler)
	mux.HandleFunc("/votes", attorney.VotesHandler)
	mux.HandleFunc("/newdraft", attorney.NewDraft2Handler)
	mux.HandleFunc("/edit", attorney.NewEditHandler)
	mux.HandleFunc("/editview/", attorney.EditViewHandler)
	mux.HandleFunc("/media/", attorney.MediaHandler)
	mux.HandleFunc("/uploadfile", attorney.UploadHandler)
	mux.HandleFunc("/createboard", attorney.CreateBoardHandler)
	mux.HandleFunc("/votecreateboard/", attorney.VoteCreateBoardHandler)
	mux.HandleFunc("/updateboard/", attorney.UpdateBoardHandler)
	mux.HandleFunc("/voteupdateboard/", attorney.VoteUpdateBoardHandler)
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
	//mux.HandleFunc("/signin", attorney.SigninHandler)
	mux.HandleFunc("/signin/", attorney.OnboardingHandler)
	mux.HandleFunc("/signout", attorney.SignoutHandler)
	mux.HandleFunc("/forgot", attorney.ForgotHandler)
	mux.HandleFunc("/passwordreset", attorney.ResetRequestHandler)
	mux.HandleFunc("/r/", attorney.ResetFromURLHandler)
	mux.HandleFunc("/reset", attorney.ResetHandler)
	mux.HandleFunc("/credentials", attorney.CredentialsHandler)
	mux.HandleFunc("/newuser", attorney.NewUserHandler)
	mux.HandleFunc("/onboarding", attorney.OnboardNewUserHandler)
	// mux.HandleFunc("/member/votes", attorney.VotesHandler)
	mux.HandleFunc("/resetpassword", attorney.ResetPasswordHandler)
	mux.HandleFunc("/credentialsreset", attorney.CredentialsResetHandler)
	mux.HandleFunc("/invitenewuser", attorney.InviteNewUserHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
	}
	finalize <- srv.ListenAndServe()
}
