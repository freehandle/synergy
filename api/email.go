package api

import (
	"fmt"
	"log"
	"net/smtp"
)

var emailMessage = "To: %v\r\n" + "Subject: Synergy Protocol Sigin\r\n" + "\r\n" + "%v\r\n"

const message = `Your email was associated to a Synergy account for the handle %v.

If you did not pursue such action, just ignore this email. Otherwise, please follow the instructions below to activate your account.

Fiirst you need to grant power of attorney to the application on axé/breeze network. Go to your wallet associated to the handle and grant power of attorney to

%v

using the fingerprint

%v

Once this message is incorporated into the network you will receive another email confirming your account activation and providing you with a provisory password to access your account.

Thank you for joining Synergy! #FreeOurHandles
`

func (a *AttorneyGeneral) sendEmail(handle, email, fingerprint, password string) {
	auth := smtp.PlainAuth("", "freemyhandle@gmail.com", password, "smtp.gmail.com")
	to := []string{email}
	body := fmt.Sprintf(message, handle, a.pk.PublicKey(), fingerprint)
	emailMsg := fmt.Sprintf(emailMessage, email, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "freemyhandle@gmail.com", to, []byte(emailMsg))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("email sent")
}
