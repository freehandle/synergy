package actions

import (
	"reflect"
	"testing"

	"github.com/lienkolabs/breeze/crypto"
)

var (
	media = &MultipartMedia{
		Epoch:  24,
		Author: crypto.Token{},
		Hash:   crypto.Hash{},
		Part:   1,
		Of:     2,
		Data:   content,
	}
)

func TestMultipartMedia(t *testing.T) {
	m := ParseMultipartMedia(media.Serialize())
	if m == nil {
		t.Error("Could not parse actions MultipartMedia")
		return
	}
	if !reflect.DeepEqual(m, media) {
		t.Error("Parse and Serialize not working for actions MultipartMedia")
	}
}
