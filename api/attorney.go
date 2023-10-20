package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	breeze "github.com/freehandle/breeze/protocol/actions"

	"github.com/freehandle/axe/attorney"
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/synergy/social/actions"
	"github.com/freehandle/synergy/social/index"
	"github.com/freehandle/synergy/social/state"
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

// credentials PasswordManager
type AttorneyGeneral struct {
	pk            crypto.PrivateKey
	signin        *SigninManager
	wallet        crypto.PrivateKey
	pending       map[crypto.Hash]actions.Action
	gateway       chan []byte
	state         *state.State
	templates     *template.Template
	indexer       *index.Index
	session       *CookieStore
	emailPassword string
	genesisTime   time.Time
	ephemeralprv  crypto.PrivateKey
	ephemeralpub  crypto.Token
}

func (a *AttorneyGeneral) Incorporate(action []byte) {
	a.state.Action(action)
}

func (a *AttorneyGeneral) SetEpoch(epoch uint64) {
	//a.epoch = uint64(epoch)
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
	a.session.Set(token, cookie, a.state.Epoch)
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
	bytes := SynergyToBreeze(action.Serialize())
	if bytes == nil {
		return nil
	}
	for n := 0; n < crypto.TokenSize; n++ {
		bytes[15+n] = author[n]
	}
	// put attorney
	util.PutToken(a.pk.PublicKey(), &bytes)
	signature := a.pk.Sign(bytes)
	util.PutSignature(signature, &bytes)
	// put zero token wallet
	util.PutToken(a.pk.PublicKey(), &bytes)
	util.PutUint64(0, &bytes) // zero fee
	signature = a.pk.Sign(bytes)
	util.PutSignature(signature, &bytes)
	return bytes
}

func (a *AttorneyGeneral) Confirmed(hash crypto.Hash) {
	delete(a.pending, hash)
}

// Translate synergy byte array to the head of the corresponding breeze instruction
// up to the specification of the signer in the axe void protocol
func SynergyToBreeze(action []byte) []byte {
	bytes := []byte{0, breeze.IVoid}                     // Breeze Void instruction version 0
	bytes = append(bytes, action[:8]...)                 // epoch (synergy)
	bytes = append(bytes, 1, 1, 0, 0, attorney.VoidType) // synergy protocol code + axe Void instruction code
	bytes = append(bytes, action[8:]...)                 //
	return bytes
}
