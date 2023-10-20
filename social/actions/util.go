package actions

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// up to 255 keywords
func PutKeywords(words []string, bytes *[]byte) {
	if len(words) == 0 || len(words) > 255 {
		*bytes = append(*bytes, 0)
	}
	*bytes = append(*bytes, byte(len(words)))
	for _, word := range words {
		util.PutString(word, bytes)
	}
}

func ParseKeywords(data []byte, position int) ([]string, int) {
	if position >= len(data) {
		return nil, position
	}
	length := int(data[position])
	words := make([]string, length)
	position += 1
	for n := 0; n < length; n++ {
		words[n], position = util.ParseString(data, position)
	}
	return words, position
}

func HashArrayToByteArray(hashes []crypto.Hash) []byte {
	output := make([]byte, 0, len(hashes)*crypto.Size)
	for _, hash := range hashes {
		output = append(output, hash[:]...)
	}
	return output
}

func ByteArrayToHashArray(bytes []byte) []crypto.Hash {
	if len(bytes)%crypto.Size != 0 {
		return nil
	}
	hashes := make([]crypto.Hash, len(bytes)/crypto.Size)
	for n := 0; n < len(hashes); n++ {
		copy(hashes[n][:], bytes[n*crypto.Size:(n+1)*crypto.Size])
	}
	return hashes
}

func TokenArrayToByteArray(tokens []crypto.Token) []byte {
	output := make([]byte, 0, len(tokens)*crypto.TokenSize)
	for _, token := range tokens {
		output = append(output, token[:]...)
	}
	return output
}

func ByteArrayToTokenArray(bytes []byte) []crypto.Token {
	if len(bytes)%crypto.TokenSize != 0 {
		return nil
	}
	tokens := make([]crypto.Token, len(bytes)/crypto.TokenSize)
	for n := 0; n < len(tokens); n++ {
		copy(tokens[n][:], bytes[n*crypto.TokenSize:(n+1)*crypto.TokenSize])
	}
	return tokens
}

func PutTokenArray(tokens []crypto.Token, bytes *[]byte) {
	translate := TokenArrayToByteArray(tokens)
	if translate == nil {
		translate = []byte{}
	}
	util.PutByteArray(translate, bytes)
}

func PutHashArray(hashes []crypto.Hash, bytes *[]byte) {
	translate := HashArrayToByteArray(hashes)
	if translate == nil {
		translate = []byte{}
	}
	util.PutByteArray(translate, bytes)
}

func ParseHashArray(data []byte, position int) ([]crypto.Hash, int) {
	byteArray, newPos := util.ParseByteArray(data, position)
	hashes := ByteArrayToHashArray(byteArray)
	return hashes, newPos
}

func ParseTokenArray(data []byte, position int) ([]crypto.Token, int) {
	byteArray, newPos := util.ParseByteArray(data, position)
	tokens := ByteArrayToTokenArray(byteArray)
	return tokens, newPos
}
