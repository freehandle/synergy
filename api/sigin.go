package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/safe"
	"github.com/freehandle/synergy/social/actions"
)

var emailSigninMessage = "To: %v\r\n" + "Subject: Synergy Protocol Sigin\r\n" + "\r\n" + "%v\r\n"

const signinWithoutHandleBody = `Your email was associated to a Synergy account for the handle %v.

If you did not pursue such action, just ignore this email. Otherwise, please follow the instructions below to activate your account.

This handle is not yet registered on the axé/breeze network. You need to register it first. You can find instructions on how to do it at

................

After creating a user related to the handle bla bla 

%v

using the fingerprint

%v

Once this message is incorporated into the network you will receive another email confirming your account activation and providing you with a provisory password to access your account.

Thank you for joining Synergy! #FreeOurHandles
`

const signinBody = `Your email was associated to a Synergy account for the handle %v.

If you did not pursue such action, just ignore this email. Otherwise, please follow the instructions below to activate your account.

Fiirst you need to grant power of attorney to the application on axé/breeze network. Go to your wallet associated to the handle and grant power of attorney to

%v

using the fingerprint

%v

Once this message is incorporated into the network you will receive another email confirming your account activation and providing you with a provisory password to access your account.

Thank you for joining Synergy! #FreeOurHandles
`

var emailPasswordMessage = "To: %v\r\n" + "Subject: Synergy Protocol Password\r\n" + "\r\n" + "%v\r\n"

const passwordMessage = `Your new password for the Synergy App account for the handle %v is

%v

Thank you for using Synergy App! #FreeOurHandles
`

type Signerin struct {
	Handle string
	//Token       crypto.Token
	Email       string
	TimeStamp   time.Time
	FingerPrint string
}

func NewSigninManager(passwords PasswordManager, emailsecret string, attorney *AttorneyGeneral) *SigninManager {
	if attorney == nil {
		log.Print("PANIC BUG: NewSigninManager called with nil attorney ")
		return nil
	}
	return &SigninManager{
		safe:          attorney.safe,
		pending:       make([]*Signerin, 0),
		passwords:     passwords,
		AttorneyToken: attorney.Token,
		Attorney:      attorney,
		Password:      emailsecret,
		Granted:       make(map[string]crypto.Token),
	}
}

type SigninManager struct {
	safe          *safe.Safe // for optional direct onboarding
	pending       []*Signerin
	passwords     PasswordManager
	AttorneyToken crypto.Token
	Password      string
	Attorney      *AttorneyGeneral
	Granted       map[string]crypto.Token
}

func (s *SigninManager) Check(user crypto.Token, password string) bool {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	//fmt.Println("hashed check", hashed)
	return s.passwords.Check(user, hashed)

}

func (s *SigninManager) Set(user crypto.Token, password string, email string) {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	//fmt.Println("hashed set", hashed)
	s.passwords.Set(user, hashed, email)
}

func (s *SigninManager) Has(token crypto.Token) bool {
	return s.passwords.Has(token)
}

func (s *SigninManager) OnboardSigner(handle, email, passwd string) bool {
	if s.safe == nil {
		log.Println("PANIC BUG: OnboardSigner called with nil safe")
		return false
	}
	ok, token := s.safe.SigninWithToken(handle, passwd, email)
	if !ok {
		return false
	}
	if err := s.safe.GrantPower(handle, s.AttorneyToken.Hex(), crypto.EncodeHash(crypto.HashToken(token))); err != nil {
		log.Println("error granting power of attorney", err)
		return false
	}
	s.Set(token, passwd, email)
	signin := actions.Signin{
		Epoch:   s.Attorney.state.Epoch,
		Author:  token,
		Reasons: "Synergy app sign in with approved power of attorney",
	}
	s.Attorney.Send([]actions.Action{&signin}, token)
	s.Granted[handle] = token
	return true
}

func (s *SigninManager) AddSigner(handle, email string, token *crypto.Token) {
	signer := &Signerin{}
	for _, pending := range s.pending {
		if signer.Handle == handle {
			signer = pending
		}
	}
	signer.Handle = handle
	signer.Email = email
	t, _ := crypto.RandomAsymetricKey()
	signer.FingerPrint = crypto.EncodeHash(crypto.HashToken(t))
	signer.TimeStamp = time.Now()
	if token != nil {
		go s.sendSigninEmail(signinBody, handle, email, signer.FingerPrint)
	} else {
		go s.sendSigninEmail(signinWithoutHandleBody, handle, email, signer.FingerPrint)
	}
	s.pending = append(s.pending, signer)
}

func randomPassword() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func (s *SigninManager) RevokeAttorney(handle string) {
	delete(s.Granted, handle)
}

func (s *SigninManager) GrantAttorney(token crypto.Token, handle, fingerprint string) {
	log.Println("to aqui", handle, fingerprint)
	for n, signer := range s.pending {
		if signer.Handle == handle && signer.FingerPrint == fingerprint {
			passwd := randomPassword()
			s.Set(token, passwd, signer.Email)
			signin := actions.Signin{
				Epoch:   s.Attorney.state.Epoch,
				Author:  token,
				Reasons: "Synergy app sign in with approved power of attorney",
			}
			s.Attorney.Send([]actions.Action{&signin}, token)
			go s.sendPasswordEmail(signer.Handle, signer.Email, passwd)
			s.pending = append(s.pending[:n], s.pending[n+1:]...)
			s.Granted[handle] = token
		}
	}
}

func (s *SigninManager) sendSigninEmail(msg, handle, email, fingerprint string) {
	auth := smtp.PlainAuth("", "freemyhandle@gmail.com", s.Password, "smtp.gmail.com")
	to := []string{email}
	body := fmt.Sprintf(msg, handle, s.AttorneyToken, fingerprint)
	emailMsg := fmt.Sprintf(emailSigninMessage, email, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "freemyhandle@gmail.com", to, []byte(emailMsg))
	if err != nil {
		log.Printf("email sending error: %v", err)
	}
	//fmt.Println(emailMsg)
}

func (s *SigninManager) sendPasswordEmail(handle, email, password string) {
	auth := smtp.PlainAuth("", "freemyhandle@gmail.com", s.Password, "smtp.gmail.com")
	to := []string{email}
	body := fmt.Sprintf(passwordMessage, handle, password)
	emailMsg := fmt.Sprintf(emailPasswordMessage, email, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "freemyhandle@gmail.com", to, []byte(emailMsg))
	if err != nil {
		log.Printf("email sending error: %v", err)
	}
	//fmt.Println(emailMsg)
}
