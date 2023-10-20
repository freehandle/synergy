package actions

import (
	"reflect"
	"testing"

	"github.com/lienkolabs/breeze/crypto"
)

var (
	signin = &Signin{
		Epoch:   26,
		Author:  crypto.Token{},
		Reasons: "signin test",
		Handle:  "first_handle",
	}
)

func TestSignin(t *testing.T) {
	s := ParseSignIn(signin.Serialize())
	if s == nil {
		t.Error("Could not parse actions Signin")
		return
	}
	if !reflect.DeepEqual(s, signin) {
		t.Error("Parse and Serialize not working for actions Signin")
	}
}
