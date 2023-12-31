package actions

import (
	"reflect"
	"testing"

	"github.com/freehandle/breeze/crypto"
)

var (
	signin = &Signin{
		Epoch:   26,
		Author:  crypto.Token{},
		Reasons: "signin test",
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
