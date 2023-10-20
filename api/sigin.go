package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/freehandle/breeze/crypto"
)

var emailSigninMessage = "To: %v\r\n" + "Subject: Synergy Protocol Sigin\r\n" + "\r\n" + "%v\r\n"

const signinMessage = `Your email was associated to a Synergy account for the handle %v.

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
	Handle      string
	Token       crypto.Token
	Email       string
	TimeStamp   time.Time
	FingerPrint crypto.Hash
}

func NewSigninManager(passwords PasswordManager, emailsecret string, attorney crypto.Token) *SigninManager {
	return &SigninManager{
		pending:       make([]*Signerin, 0),
		passwords:     passwords,
		AttorneyToken: attorney,
		Password:      emailsecret,
	}
}

type SigninManager struct {
	pending       []*Signerin
	passwords     PasswordManager
	AttorneyToken crypto.Token
	Password      string
}

func (s *SigninManager) Check(user crypto.Token, password string) bool {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	return s.passwords.Check(user, hashed)

}

func (s *SigninManager) Set(user crypto.Token, password string, email string) {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	s.passwords.Set(user, hashed, email)
}

func (s *SigninManager) Has(token crypto.Token) bool {
	return s.passwords.Has(token)
}

func (s *SigninManager) AddSigner(handle, email string, token crypto.Token) {
	signer := &Signerin{}
	for _, pending := range s.pending {
		if signer.Handle == handle {
			signer = pending
		}
	}
	signer.Handle = handle
	signer.Email = email
	t, _ := crypto.RandomAsymetricKey()
	signer.FingerPrint = crypto.HashToken(t)
	signer.TimeStamp = time.Now()
	signer.Token = token
	go s.sendSigninEmail(handle, email, crypto.EncodeHash(signer.FingerPrint))
	s.pending = append(s.pending, signer)
}

func randomPassword() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func (s *SigninManager) GrantAttorney(token crypto.Token, fingerprint crypto.Hash) {
	for n, signer := range s.pending {
		if signer.Token.Equal(token) && signer.FingerPrint.Equal(fingerprint) {
			passwd := randomPassword()
			s.Set(token, passwd, signer.Email)
			go s.sendPasswordEmail(signer.Handle, signer.Email, passwd)
			s.pending = append(s.pending[:n], s.pending[n+1:]...)
		}
	}
}

func (s *SigninManager) sendSigninEmail(handle, email, fingerprint string) {
	auth := smtp.PlainAuth("", "freemyhandle@gmail.com", s.Password, "smtp.gmail.com")
	to := []string{email}
	body := fmt.Sprintf(signinMessage, handle, s.AttorneyToken, fingerprint)
	emailMsg := fmt.Sprintf(emailSigninMessage, email, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "freemyhandle@gmail.com", to, []byte(emailMsg))
	if err != nil {
		log.Printf("email sending error: %v", err)
	}
	fmt.Println(emailMsg)
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
	fmt.Println(emailMsg)
}
