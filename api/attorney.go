package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/freehandle/breeze/consensus/messages"
	breeze "github.com/freehandle/breeze/protocol/actions"
	"github.com/freehandle/safe"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
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
	Token         crypto.Token
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
	serverName    string
	safe          *safe.Safe // optional link to safe for direct onbboarding
}

func (a *AttorneyGeneral) IncorporateGrantPower(handle string, grant *attorney.GrantPowerOfAttorney) {
	if grant != nil {
		a.signin.GrantAttorney(grant.Author, handle, string(grant.Fingerprint))
	}
}

func (A *AttorneyGeneral) RegisterAxeDataBase(axe state.HandleProvider) {
	A.state.Axe = axe
}

func (a *AttorneyGeneral) IncorporateRevokePower(handle string) {
	//TODO: implement interface for a user to revoke power of attorney
}

func (a *AttorneyGeneral) Incorporate(action []byte) {
	if err := a.state.Action(action); err != nil {
		//fmt.Println("Error incorporating action:", err, action)
	}

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
		fmt.Println("Dressed action:", dressed)
		a.gateway <- append([]byte{messages.MsgAction}, dressed...)
		//a.gateway.Action(dressed)
	}
}

// Dress a giving action with current epoch, attorney´s author
// attorneys signature, attorneys wallet and wallet signature
func (a *AttorneyGeneral) DressAction(action actions.Action, author crypto.Token) []byte {
	bytes := SynergyToBreeze(action.Serialize(), a.state.Epoch)
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
func SynergyToBreeze(action []byte, epoch uint64) []byte {
	if action == nil {
		log.Print("PANIC BUG: SynergyToBreeze called with nil action ")
		return nil
	}
	bytes := []byte{0, breeze.IVoid}                     // Breeze Void instruction version 0
	util.PutUint64(epoch, &bytes)                        // epoch (synergy)
	bytes = append(bytes, 1, 1, 0, 0, attorney.VoidType) // synergy protocol code + axe Void instruction code
	bytes = append(bytes, action[8:]...)                 //
	return bytes
}
