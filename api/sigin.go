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

const wellcomeBody = `
Your email address has been recently assigned to a synergy account for the handle %s. 

If you did not persue this action, it is safe to just ignore this message. 

Otherwise, we are happy to welcome you to Synergy, the flagship of the FreeHandle project! 

https://github.com/freehandle/synergy is designed for collaboration and collective construction. It runs on top of the breeze network and on top of handles social protocol, both also part of the FreeHandle project.

By providing this email address upon signing into synergy, the following actions were triggered and performed on your behalf:

- A safe wallet account was assigned to the handle %s. Through this wallet you can manage which applications will be authorized to sign instructions on your behalf, much like granting a power of attorney.

- The synergy app was granted this power of attorney to sign instructions on your behalf (i.e. on behalf of the handle %s). You are free to revoke this power of attorney at any time with no harm to your account, data or wallet. Should you choose to do that, you are also free to grant it back to synergy (or any other application, for that matter) at any time in the future.

- Your synergy account was created and associated to the handle %s and you may now freely use and collaborate on the synergy app. It is important to know that any instructions performed on the app will be processed on the breeze network, and cannot be futurally erased, as they will be permanently incorporated into the breeze blockchain.

- By choosing to provide an email address you secured the possibility of password recovery. This email address will be used for that purpose if by any chance you loose access to %s's Synergy account. 

Thank you for joining Synergy! 

#FreeOurHandles`

const resetBody = `

Your password reset request for synergy was received. To reset your password please follow the link below

https://%v

Thanks for #FreeingOurHandles!
`

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

type SMTPGmail struct {
	Password string
	From     string
}

func (s *SMTPGmail) Send(to, subject, body string) bool {
	auth := smtp.PlainAuth("", s.From, s.Password, "smtp.gmail.com")
	emailMsg := fmt.Sprintf("To: %s\r\n"+"Subject: %s\r\n"+"\r\n"+"%s\r\n", to, subject, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, s.From, []string{to}, []byte(emailMsg))
	if err != nil {
		log.Printf("email sending error: %v", err)
		return false
	}
	return true
}

type Mailer interface {
	Send(to, subject, body string) bool
}

type Signerin struct {
	Handle string
	//Token       crypto.Token
	Email       string
	TimeStamp   time.Time
	FingerPrint string
}

func NewSigninManager(passwords PasswordManager, mail Mailer, attorney *AttorneyGeneral) *SigninManager {
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
		mail:          mail,
		Granted:       make(map[string]crypto.Token),
	}
}

type SigninManager struct {
	safe          *safe.Safe // for optional direct onboarding
	pending       []*Signerin
	passwords     PasswordManager
	AttorneyToken crypto.Token
	mail          Mailer
	Attorney      *AttorneyGeneral
	Granted       map[string]crypto.Token
}

func (s *SigninManager) RequestReset(user crypto.Token, email, domain string) bool {
	if !s.passwords.HasWithEmail(user, email) {
		return false
	}
	reset := s.passwords.AddReset(user, email)
	url := fmt.Sprintf("%s/r/%s", domain, reset)
	if reset == "" {
		return false
	}
	go s.mail.Send(email, "Synergy password reset", fmt.Sprintf(resetBody, url))
	return true
}

func (s *SigninManager) Reset(user crypto.Token, url, password string) bool {
	return s.passwords.DropReset(user, url, password)
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
	go s.mail.Send(email, "Synergy Protocol Welcome", fmt.Sprintf(wellcomeBody, handle, handle, handle, handle, handle))
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
	/*auth := smtp.PlainAuth("", "freemyhandle@gmail.com", s.Password, "smtp.gmail.com")
	to := []string{email}
	body := fmt.Sprintf(msg, handle, s.AttorneyToken, fingerprint)
	emailMsg := fmt.Sprintf(emailSigninMessage, email, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "freemyhandle@gmail.com", to, []byte(emailMsg))
	if err != nil {
		log.Printf("email sending error: %v", err)
	}*/
	body := fmt.Sprintf(msg, handle, s.AttorneyToken, fingerprint)
	s.mail.Send(email, "Synergy Protocol Signin", body)
	//fmt.Println(emailMsg)
}

func (s *SigninManager) sendPasswordEmail(handle, email, password string) {
	/*auth := smtp.PlainAuth("", "freemyhandle@gmail.com", s.Password, "smtp.gmail.com")
	to := []string{email}
	body := fmt.Sprintf(passwordMessage, handle, password)
	emailMsg := fmt.Sprintf(emailPasswordMessage, email, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "freemyhandle@gmail.com", to, []byte(emailMsg))
	if err != nil {
		log.Printf("email sending error: %v", err)
	}*/
	body := fmt.Sprintf(passwordMessage, handle, password)
	s.mail.Send(email, "Synergy Protocol Password", body)
	//fmt.Println(emailMsg)
}
