package actions

import "github.com/freehandle/breeze/crypto"

const (
	AVote byte = iota
	ACreateCollective
	AUpdateCollective
	ARequestMembership
	ARemoveMember
	ADraft
	AEdit
	AMultipartMedia
	ACreateBoard
	AUpdateBoard
	APin
	ABoardEditor
	AReleaseDraft
	AImprintStamp
	AReact
	ASignIn
	ACreateEvent
	ACancelEvent
	AUpdateEvent
	ACheckinEvent
	AGreetCheckinEvent
	AUnknown
)

// toda constante declarada abaixo de Avote eh a de cima mais um, comeca em 0

type Action interface {
	Serialize() []byte
	Authored() crypto.Token
	Hashed() crypto.Hash
	Reasoning() string
}

// atribuindo um byte pra cada uma das acoes listadas
func ActionKind(data []byte) byte {
	if len(data) < 8+crypto.TokenSize+1 {
		return AUnknown
	}
	actionByte := data[8+crypto.TokenSize]
	if actionByte >= AUnknown {
		return AUnknown
	}
	return actionByte
}
