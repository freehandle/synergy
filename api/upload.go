package api

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social/actions"
)

const maxFileSize = 10000

type TruncatedFile struct {
	Hash  crypto.Hash
	Parts [][]byte
}

func splitBytes(bytes []byte) *TruncatedFile {
	truncated := TruncatedFile{
		Hash:  crypto.Hasher(bytes),
		Parts: make([][]byte, len(bytes)/maxFileSize+1),
	}
	for n := 0; n < len(truncated.Parts); n++ {
		if (n+1)*maxFileSize >= len(bytes) {
			truncated.Parts[n] = bytes[n*maxFileSize:]
		} else {
			truncated.Parts[n] = bytes[n*maxFileSize : (n+1)*maxFileSize]
		}

	}
	return &truncated
}

func (a *Attorney) UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(maxFileSize)
	file, handle, err := r.FormFile("fileUpload")
	if err != nil {
		log.Printf("Error Retrieving the File: %v\n", err)
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("errors reading file bytes: %v\n", err)
	}
	var actionArray []actions.Action
	name := handle.Filename
	parts := strings.Split(name, ".")
	ext := parts[len(parts)-1]
	switch r.FormValue("action") {
	case "Draft":
		actionArray, err = DraftForm(r, a.state.MembersIndex, fileBytes, ext).ToAction()
	case "Edit":
		actionArray, err = EditForm(r, a.state.MembersIndex, fileBytes, ext).ToAction()
	}
	if err == nil && len(actionArray) > 0 {
		a.Send(actionArray)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AttorneyGeneral) UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(maxFileSize)
	author := a.Author(r)
	file, handle, err := r.FormFile("fileUpload")
	if err != nil {
		log.Printf("Error Retrieving the File: %v\n", err)
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("errors reading file bytes: %v\n", err)
	}
	var actionArray []actions.Action
	name := handle.Filename
	parts := strings.Split(name, ".")
	ext := parts[len(parts)-1]
	switch r.FormValue("action") {
	case "Draft":
		actionArray, err = DraftForm(r, a.state.MembersIndex, fileBytes, ext).ToAction()
	case "Edit":
		actionArray, err = EditForm(r, a.state.MembersIndex, fileBytes, ext).ToAction()
	}
	if err == nil && len(actionArray) > 0 {
		a.Send(actionArray, author)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
