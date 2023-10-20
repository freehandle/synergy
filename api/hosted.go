package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/util"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/index"
	"github.com/lienkolabs/synergy/social/state"
)

const cookieName = "synergySession"

const cookieLifeItemSeconds = 60 * 60 * 24 * 7 // 1 week

func newCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(value),
		MaxAge:   cookieLifeItemSeconds,
		Secure:   true,
		HttpOnly: true,
	}
}

type AttorneyGeneral struct {
	epoch       uint64
	pk          crypto.PrivateKey
	credentials PasswordManager
	wallet      crypto.PrivateKey
	pending     map[crypto.Hash]actions.Action
	//gateway       social.Gatewayer
	gateway       chan []byte
	state         *state.State
	templates     *template.Template
	indexer       *index.Index
	session       *CookieStore
	emailPassword string
	//session      map[string]crypto.Token
	//sessionend   map[uint64][]string
	genesisTime  time.Time
	ephemeralprv crypto.PrivateKey
	ephemeralpub crypto.Token
}

func (a *AttorneyGeneral) Incorporate(action []byte) {
	a.state.Action(action)
}

func (a *AttorneyGeneral) SetEpoch(epoch uint64) {
	a.epoch = uint64(epoch)
	a.state.Epoch = epoch
}

func (a *AttorneyGeneral) DestroySession(token crypto.Token, cookie string) {
	a.session.Unset(token, cookie)
}

func (a *AttorneyGeneral) CreateSession(handle string) string {
	token, ok := a.state.MembersIndex[handle]
	if !ok {
		return ""
	}
	seed := make([]byte, 32)
	if n, err := rand.Read(seed); n != 32 || err != nil {
		log.Printf("unexpected error in cookie generation:%v", err)
		return ""
	}
	cookie := hex.EncodeToString(seed)
	a.session.Set(token, cookie, a.epoch)
	return cookie
}

func (a *AttorneyGeneral) Author(r *http.Request) crypto.Token {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return crypto.ZeroToken
	}
	if token, ok := a.session.Get(cookie.Value); ok {
		return token
	}
	return crypto.ZeroToken
}

func (a *AttorneyGeneral) Handle(r *http.Request) string {
	author := a.Author(r)
	handle := a.state.Members[crypto.HashToken(author)]
	return handle
}

func (a *AttorneyGeneral) Send(all []actions.Action, author crypto.Token) {
	for _, action := range all {
		dressed := a.DressAction(action, author)
		a.gateway <- dressed
		//a.gateway.Action(dressed)
	}
}

// Dress a giving action with current epoch, attorney´s author
// attorneys signature, attorneys wallet and wallet signature
func (a *AttorneyGeneral) DressAction(action actions.Action, author crypto.Token) []byte {
	bytes := action.Serialize()
	dress := []byte{0}
	util.PutUint64(a.epoch, &dress)
	util.PutToken(author, &dress)
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

func (a *AttorneyGeneral) Confirmed(hash crypto.Hash) {
	delete(a.pending, hash)
}
