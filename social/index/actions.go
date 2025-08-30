package index

import (
	"fmt"
	"net/url"
	"time"

	"github.com/freehandle/breeze/crypto"

	"github.com/freehandle/synergy/social/actions"
	"github.com/freehandle/synergy/social/state"
)

func (i *Index) ActionToObjects(action actions.Action) []crypto.Hash {
	switch v := action.(type) {
	case *actions.ImprintStamp:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf)), v.Hash}
	case *actions.CreateEvent:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.CancelEvent:
		return []crypto.Hash{v.Hash}
	case *actions.UpdateEvent:
		return []crypto.Hash{v.EventHash}
	case *actions.CheckinEvent:
		return []crypto.Hash{v.EventHash}
	case *actions.GreetCheckinEvent:
		return []crypto.Hash{v.EventHash}
	case *actions.CreateBoard:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.UpdateBoard:
		hash := crypto.Hasher([]byte(v.Board))
		return []crypto.Hash{hash}
	case *actions.Pin:
		return []crypto.Hash{crypto.Hasher([]byte(v.Board)), v.Draft}
	case *actions.BoardEditor:
		return []crypto.Hash{crypto.Hasher([]byte(v.Board))}
	case *actions.Draft:
		if v.OnBehalfOf != "" {
			return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
		}
		return nil
	case *actions.ReleaseDraft:
		return []crypto.Hash{v.ContentHash}
	case *actions.Edit:
		if v.OnBehalfOf != "" {
			return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf)), v.EditedDraft}
		}
		return []crypto.Hash{v.EditedDraft}
	case *actions.React:
		if v.OnBehalfOf != "" {
			return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf)), v.Hash}
		}
		return []crypto.Hash{v.Hash}
	case *actions.CreateCollective:
		return []crypto.Hash{crypto.ZeroHash}
	case *actions.UpdateCollective:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.RequestMembership:
		return []crypto.Hash{crypto.Hasher([]byte(v.Collective))}
	case *actions.RemoveMember:
		return []crypto.Hash{crypto.Hasher([]byte(v.OnBehalfOf))}
	case *actions.Signin:
		return []crypto.Hash{crypto.ZeroHash}
	}
	return nil
}

func fmtHandle(handle string) string {
	return fmt.Sprintf("<a href=\"./member/%v\">%v</a>", url.QueryEscape(handle), handle)
}

func fmtCollective(collective string) string {
	return fmt.Sprintf("<a href=\"./collective/%v\">%v</a>", url.QueryEscape(collective), collective)
}

func fmtBoard(board string) string {
	return fmt.Sprintf("<a href=\"./board/%v\">%v</a>", url.QueryEscape(board), board)
}

func fmtDraft(draft string, hash crypto.Hash) string {
	if len(draft) > 40 {
		draft = draft[:40] + "..."
	}
	return fmt.Sprintf("<a href=\"./draft/%v\">&ldquo;%v&rdquo;</a>", crypto.EncodeHash(hash), draft)
}

func fmtEvent(date time.Time, hash crypto.Hash) string {
	return fmt.Sprintf("<a href=\"./event/%v\">%v</a>", crypto.EncodeHash(hash), date.Format("Mon Jan 2 at 15:04 MST"))
}

func fmtAuthors(authors state.Consensual, s *state.State) string {
	if authors == nil {
		return ""
	}
	if authors.CollectiveName() != "" {
		return fmt.Sprintf("em nome de %v", fmtCollective(authors.CollectiveName()))
	}
	authorsCaption := ""
	count := 0
	for token, _ := range authors.ListOfMembers() {
		if handle, ok := s.Members[crypto.HashToken(token)]; ok {
			if count == 0 {
				authorsCaption = fmtHandle(handle)
			} else if count == 1 {
				authorsCaption = fmt.Sprintf("%v e %v", authorsCaption, fmtHandle(handle))
			} else {
				authorsCaption = fmt.Sprintf("%v et al.", authorsCaption)
				break
			}
			count += 1
		}
	}
	if len(authorsCaption) == 0 {
		return ""
	}
	return fmt.Sprintf("by %v", authorsCaption)
}

func (i *Index) ActionToFormatedString(action actions.Action) (string, string, uint64) {
	switch v := action.(type) {
	case *actions.ImprintStamp:
		if draft, ok := i.state.Drafts[v.Hash]; ok {
			return fmt.Sprintf("%v colocou um selo %v", fmtCollective(v.OnBehalfOf), fmtDraft(draft.Title, draft.DraftHash)), "awareness", v.Epoch
		}
	case *actions.CreateEvent:
		eventhash := v.Hashed()
		if event, ok := i.state.Events[eventhash]; ok {
			isPublic := "privado"
			if event.Public {
				isPublic = "público"
			}
			isOpen := "fechado"
			if event.Open {
				isOpen = "aberto"
			}
			return fmt.Sprintf("%v marcou um evento %v, %v em %v", fmtCollective(event.Collective.Name), isPublic, isOpen, fmtEvent(v.StartAt, eventhash)), "new stuff", v.Epoch
		}
	case *actions.CancelEvent:
		if event, ok := i.state.Events[v.Hash]; ok {
			return fmt.Sprintf("%v cancelou um evento em %v", fmtCollective(event.Collective.Name), fmtEvent(event.StartAt, v.Hash)), "update", v.Epoch
		}
	case *actions.UpdateEvent:
		if event, ok := i.state.Events[v.EventHash]; ok {
			return fmt.Sprintf("%v atualizou um evento em %v", fmtCollective(event.Collective.Name), fmtEvent(event.StartAt, v.EventHash)), "update", v.Epoch
		}
	case *actions.CheckinEvent:
		if event, ok := i.state.Events[v.EventHash]; ok {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			return fmt.Sprintf("%v confirmou participação no evento %v por %v ", fmtHandle(handle), fmtEvent(event.StartAt, v.EventHash), event.Collective.Name), "people", v.Epoch
		}
	case *actions.GreetCheckinEvent:
		return "", "", 0
	case *actions.CreateBoard:
		boardhash := v.Hashed()
		if board, ok := i.state.Boards[boardhash]; ok {
			return fmt.Sprintf("%v criou um novo mural %v", fmtCollective(board.Collective.Name), fmtBoard(board.Name)), "new stuff", v.Epoch
		}
	case *actions.UpdateBoard:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			return fmt.Sprintf("%v atualizou um mural %v", fmtCollective(board.Collective.Name), fmtBoard(board.Name)), "update", v.Epoch
		}
	case *actions.Pin:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if draft, ok := i.state.Drafts[v.Draft]; ok {
				pinaction := "desafixado de"
				if v.Pin {
					pinaction = "fixado em"
				}
				return fmt.Sprintf(`%v %v %v em nome de %v`, fmtDraft(draft.Title, draft.DraftHash), pinaction, fmtBoard(board.Name), fmtCollective(board.Collective.Name)), "awareness", v.Epoch
			}
		}
	case *actions.BoardEditor:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			editor := i.state.Members[crypto.HashToken(v.Editor)]
			editorship := []string{"removido de", "de"}
			if v.Insert {
				editorship = []string{"incluido por", "por"}
			}
			return fmt.Sprintf("%v %v mural de editores de %v em nome de %v", fmtHandle(editor), editorship[0], fmtBoard(board.Name), fmtCollective(board.Collective.Name)), "people", v.Epoch
		}
	case *actions.Draft:
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			authors := fmtAuthors(draft.Authors, i.state)
			if authors != "" {
				if draft.PreviousVersion != nil {
					return fmt.Sprintf("Nova versão de %v foi publicada %v", fmtDraft(draft.Title, draft.DraftHash), authors), "update", v.Epoch
				} else {
					return fmt.Sprintf("Novo esboço %v foi publicado %v", fmtDraft(draft.Title, draft.DraftHash), authors), "new stuff", v.Epoch
				}
			}
		}
	case *actions.ReleaseDraft:
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			authors := fmtAuthors(draft.Authors, i.state)
			return fmt.Sprintf("Esboço %v foi lançado %v", fmtDraft(draft.Title, draft.DraftHash), authors), "update", v.Epoch
		}
	case *actions.Edit:
		if edit, ok := i.state.Edits[v.ContentHash]; ok {
			draft := i.state.Drafts[v.EditedDraft]
			authors := fmtAuthors(edit.Authors, i.state)
			if draft != nil && authors != "" {
				return fmt.Sprintf("Uma edição para %v foi proposta por %v", fmtDraft(draft.Title, draft.DraftHash), authors), "update", v.Epoch
			}
		}
	case *actions.React:
		if handle, ok := i.state.Members[crypto.HashToken(v.Author)]; ok {
			if collective, ok := i.state.Collectives[v.Hash]; ok {
				return fmt.Sprintf("%v reagiu ao coletivo %v. <span class=\"reaction\"> %v </span>", fmtHandle(handle), fmtCollective(collective.Name), v.Reasons), "react collective", v.Epoch
			} else if event, ok := i.state.Events[v.Hash]; ok {
				return fmt.Sprintf("%v reagiu ao %v evento por %v. <span class=\"reaction\"> %v </span>", fmtHandle(handle), fmtEvent(event.StartAt, v.Hash), fmtCollective(event.Collective.Name), v.Reasons), "react event", v.Epoch
			} else if board, ok := i.state.Boards[v.Hash]; ok {
				return fmt.Sprintf("%v reagiu ao mural %v. <span class=\"reaction\"> %v </span>", fmtHandle(handle), fmtBoard(board.Name), v.Reasons), "react board", v.Epoch
			} else if draft, ok := i.state.Drafts[v.Hash]; ok {
				authors := fmtAuthors(draft.Authors, i.state)
				return fmt.Sprintf("%v reagiu ao esboço %v de %v. <span class=\"reaction\"> %v </span>", fmtHandle(handle), fmtDraft(draft.Title, draft.DraftHash), authors, v.Reasons), "react draft", v.Epoch
			} else if edit, ok := i.state.Edits[v.Hash]; ok {
				authors := fmtAuthors(edit.Authors, i.state)
				return fmt.Sprintf("%v reagiu à edição proposta por %v para %v. <span class=\"reaction\"> %v </span>", fmtHandle(handle), authors, fmtDraft(draft.Title, draft.DraftHash), v.Reasons), "react edit", v.Epoch
			}
		}
	case *actions.CreateCollective:
		collectivehash := crypto.Hasher([]byte(v.Name))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			return fmt.Sprintf("Novo coletivo %v criado", fmtCollective(collective.Name)), "new stuff", v.Epoch
		}
	case *actions.UpdateCollective:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			return fmt.Sprintf("Coletivo %v atualizado", fmtCollective(collective.Name)), "update", v.Epoch
		}
	case *actions.RequestMembership:
		return "", "", 0
	case *actions.RemoveMember:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			member := i.state.Members[crypto.HashToken(v.Member)]
			return fmt.Sprintf("%v foi removido de %v", fmtHandle(member), fmtCollective(collective.Name)), "people", v.Epoch
		}
	case *actions.Signin:
		authorhash := crypto.HashToken(v.Author)
		if membro, ok := i.state.Members[authorhash]; ok {
			return fmt.Sprintf("%v ingressou no Motiró", fmtHandle(membro)), "people", v.Epoch
		}
	}
	return "", "", 0
}

func (i *Index) ActionToString(action actions.Action, status state.ConsensusState) (string, string, crypto.Token, uint64, string) {
	switch v := action.(type) {
	case *actions.ImprintStamp:
		if draft, ok := i.state.Drafts[v.Hash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v colocou um selo em %v", v.OnBehalfOf, draft.Title), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "stamp"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs selo de %v para %v", handle, v.OnBehalfOf, draft.Title), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "stamp"
			}
			if status == state.Against {
				return fmt.Sprintf("%v selo negado em nome de %v", v.OnBehalfOf, draft.Title), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "stamp"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CreateEvent:
		// hash do evento eh o hash da acao do evento
		eventhash := v.Hashed()
		if event, ok := i.state.Events[eventhash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("evento em %v criado em nome de %v", v.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(eventhash), v.Author, v.Epoch, "create event"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs um evento em %v em nome de %v", handle, v.StartAt.Format("2006-01-02"), v.OnBehalfOf), crypto.EncodeHash(eventhash), v.Author, v.Epoch, "create event"
			}
			if status == state.Against {
				return fmt.Sprintf("evento em %v negado em nome de %v", v.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(eventhash), v.Author, v.Epoch, "create event"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CancelEvent:
		// hash eh o hash do evento original
		if event, ok := i.state.Events[v.Hash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("evento em %v cancelado em nome de %v", event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.Hash), v.Author, v.Epoch, "cancel event"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs cancelamento do evento em %v em nome de %v", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.Hash), v.Author, v.Epoch, "cancel event"
			}
			if status == state.Against {
				return fmt.Sprintf("evento em %v cancelado em nome de %v", event.Collective.Name, event.Collective.Name), crypto.EncodeHash(v.Hash), v.Author, v.Epoch, "cancel event"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.UpdateEvent:
		// hash eh o hash do evento original
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("evento em %v atualizado em nome de %v", event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "update event"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs atualização para o evento em %v em nome de %v", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "update event"
			}
			if status == state.Against {
				return fmt.Sprintf("atualização do evento em %v negada em nome de %v", event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "update event"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CheckinEvent:
		// hash eh o hash do evento
		if event, ok := i.state.Events[v.EventHash]; ok {
			// verificar como status se comporta aqui
			if status == state.Favorable {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v confirmou presença no evento em %v por %v ", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "event checkin"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.GreetCheckinEvent:
		if event, ok := i.state.Events[v.EventHash]; ok {
			// verificar como status se comporta aqui
			if status == state.Favorable {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v recebeu boas-vindas para evento em %v por %v ", handle, event.StartAt.Format("2006-01-02"), event.Collective.Name), crypto.EncodeHash(v.EventHash), v.Author, v.Epoch, "event greet"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.CreateBoard:
		// hash do board eh o hash do nome do board que esta sendo criado
		boardhash := crypto.Hasher([]byte(v.Name))

		if status == state.Favorable {
			if board, ok := i.state.Boards[boardhash]; ok {
				return fmt.Sprintf("%v criado em nome de %v", board.Name, v.OnBehalfOf), crypto.EncodeHash(boardhash), v.Author, v.Epoch, "create board"
			}
		}
		if status == state.Undecided {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			return fmt.Sprintf("%v propôs criação de %v em nome de %v", handle, v.Name, v.OnBehalfOf), "", v.Author, v.Epoch, "create board"
		}
		if status == state.Against {
			return fmt.Sprintf("criação de %v foi negada por %v", v.Name, v.OnBehalfOf), crypto.EncodeHash(boardhash), v.Author, v.Epoch, "create board"
		}
		return "", "", v.Author, 0, ""
	case *actions.UpdateBoard:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v foi atualizado em nome de %v", board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, "update board"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs atualização de %v em nome de %v", handle, board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, "update board"
			}
			if status == state.Against {
				return fmt.Sprintf("atualização de %v foi negada por %v", board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, "update board"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Pin:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if draft, ok := i.state.Drafts[v.Draft]; ok {
				pinaction := []string{"desafixado de", "desafixar"}
				if v.Pin {
					pinaction = []string{"fixado em", "fixar"}
				}
				if status == state.Favorable {
					return fmt.Sprintf("%v %v %v por %v", draft.Title, pinaction[0], board.Name, board.Collective.Name), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, pinaction[1]
				}
				if status == state.Undecided {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v propôs %v %v de %v em nome de %v", handle, pinaction[1], draft.Title, board.Name, board.Collective.Name), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, pinaction[1]
				}
				if status == state.Against {
					return fmt.Sprintf("%v %v de %v em nome de %v foi negado", pinaction[1], draft.Title, board.Name, board.Collective.Name), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, pinaction[1]
				}
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.BoardEditor:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			editor := i.state.Members[crypto.HashToken(v.Editor)]
			editorship := []string{"removido de", "remoção de", "de", "editor removal"}
			if v.Insert {
				editorship = []string{"incluído como", "inclusão de", "em", "editor inclusion"}
			}
			if status == state.Favorable {
				return fmt.Sprintf("%v %v editoria de %v em nome de %v", editor, editorship[0], board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, editorship[3]
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs %v %v %v %v editoria em nome de %v", handle, editorship[1], editor, editorship[2], board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, editorship[3]
			}
			if status == state.Against {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v %v %v %v %v editoria em nome de %v foi negada", handle, editorship[1], editor, editorship[2], board.Name, board.Collective.Name), crypto.EncodeHash(hash), v.Author, v.Epoch, editorship[3]
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Draft:
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			if draft.Authors.CollectiveName() != "" {
				if status == state.Favorable {
					return fmt.Sprintf("%v foi criado em nome de %v", draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "new draft"
				}
				if status == state.Undecided {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v propôs %v em nome de %v", handle, draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "new draft"
				}
				if status == state.Against {
					return fmt.Sprintf("proposta de esboço %v em nome de %v foi negada", draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "new draft"
				}
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.ReleaseDraft:
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v foi lançado em nome de %v", draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "release"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs lançamento de %v em nome de %v", handle, draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "release"
			}
			if status == state.Against {
				return fmt.Sprintf("lançamento de %v negado em nome de %v", draft.Title, draft.Authors.CollectiveName()), crypto.EncodeHash(draft.DraftHash), v.Author, v.Epoch, "release"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Edit:
		if edit, ok := i.state.Edits[v.ContentHash]; ok {
			draft := i.state.Drafts[v.EditedDraft]
			if edit.Authors.CollectiveName() != "" {
				if status == state.Favorable {
					return fmt.Sprintf("edição de %v sugerida em nome de %v", draft.Title, edit.Authors.CollectiveName()), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
				}
				if status == state.Undecided {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v propôs edição para %v em nome de %v", handle, draft.Title, edit.Authors.CollectiveName()), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
				}
				if status == state.Against {
					return fmt.Sprintf("sugestão de edição de %v em nome de %v foi negada", draft.Title, edit.Authors.CollectiveName()), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
				}
			}
			if status == state.Favorable {
				return fmt.Sprintf("sugestão de edição para %v", draft.Title), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs edição para %v", handle, draft.Title), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
			}
			if status == state.Against {
				return fmt.Sprintf("sugestão de edição para %v negada pelos autores", draft.Title), crypto.EncodeHash(v.ContentHash), v.Author, v.Epoch, "edit"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.React:
		// reacthash := v.Hashed()
		// if ok := i.state.Reactions[reacthash]; ok {
		// if status {
		// 	return fmt.Sprintf("%v draft was released on behalf of %v", draft.Title, draft.Authors.CollectiveName()), v.Epoch
		// } else {
		// 	handle := i.state.Members[crypto.HashToken(v.Author)]
		// 	return fmt.Sprintf("%v proposed %v draft release on behalf of %v", handle, draft.Title, draft.Authors.CollectiveName()), v.Epoch
		// }
		// }
		return "", "", v.Author, 0, ""
	case *actions.CreateCollective:
		// verificar como status se comporta aqui (porque provavelmente todo coletivo é automaticamente criado)
		collectivehash := crypto.Hasher([]byte(v.Name))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v foi criado", collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "create collective"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs criação de %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "create collective"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.UpdateCollective:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v atualizado", collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "update collective"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs atualização para %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "update collective"
			}
			if status == state.Against {
				return fmt.Sprintf("atualização para %v negada", collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "update collective"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.RequestMembership:
		collectivehash := crypto.Hasher([]byte(v.Collective))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			if v.Include {
				if status == state.Favorable {
					return fmt.Sprintf("%v tornou-se membro do coletivo %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "request membership"
				}
				if status == state.Undecided {
					return fmt.Sprintf("%v candidatou inclusão ao coletivo %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "request membership"
				}
				if status == state.Against {
					return fmt.Sprintf("inclusão de %v ao coletivo %v foi negada", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "request membership"
				}
			}
			if status == state.Favorable {
				return fmt.Sprintf("%v left %v", handle, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "request membership"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.RemoveMember:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			member := i.state.Members[crypto.HashToken(v.Member)]
			if status == state.Favorable {
				return fmt.Sprintf("%v foi removido do coletivo %v", member, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "remove member"
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v pediu remoção de %v do coletivo %v", handle, member, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "remove member"
			}
			if status == state.Against {
				return fmt.Sprintf("remoção de %v do coletivo %v foi negada", member, collective.Name), crypto.EncodeHash(collectivehash), v.Author, v.Epoch, "remove member"
			}
		}
		return "", "", v.Author, 0, ""
	case *actions.Signin:
		authorhash := crypto.HashToken(v.Author)
		if _, ok := i.state.Members[authorhash]; ok {
			// verificar como status se comporta aqui (porque acho que nao tem consenso nunca ne)
			if status == state.Favorable {
				return fmt.Sprintf("%v ingressou no Motiró"), "", v.Author, v.Epoch, "sign in"
			}
		}
		return "", "", v.Author, 0, ""
	}
	return "", "", crypto.ZeroToken, 0, ""
}

func (i *Index) ActionToStringWithLinks(action actions.Action, status state.ConsensusState) (string, uint64, string) {
	switch v := action.(type) {
	case *actions.ImprintStamp:
		if draft, ok := i.state.Drafts[v.Hash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v deu um selo para %v", fmtCollective(v.OnBehalfOf), fmtDraft(draft.Title, draft.DraftHash)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs um selo de %v para %v", handle, fmtCollective(v.OnBehalfOf), fmtDraft(draft.Title, draft.DraftHash)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("selo de %v para %v foi negado", fmtCollective(v.OnBehalfOf), fmtDraft(draft.Title, draft.DraftHash)), v.Epoch, v.Reasons
			}
		}
	case *actions.CreateEvent:
		eventhash := v.Hashed()
		if status == state.Favorable {
			return fmt.Sprintf("evento em %v criado por %v", fmtEvent(v.StartAt, eventhash), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
		}
		if status == state.Undecided {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			return fmt.Sprintf("%v propôs evento em %v em nome de %v", fmtHandle(handle), fmtEvent(v.StartAt, eventhash), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
		}
		if status == state.Against {
			return fmt.Sprintf("evento em %v em nome de %v foi negado", fmtEvent(v.StartAt, eventhash), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
		}
	case *actions.CancelEvent:
		if event, ok := i.state.Events[v.Hash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("evento em %v foi cancelado por %v", fmtEvent(event.StartAt, v.Hash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs cancelar evento em %v em nome de %v", fmtHandle(handle), fmtEvent(event.StartAt, v.Hash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("cancelamento do evento em %v em nome de %v foi negado", fmtEvent(event.StartAt, v.Hash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
			}
		}
	case *actions.UpdateEvent:
		if event, ok := i.state.Events[v.EventHash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("evento em %v atualizado por %v", fmtEvent(event.StartAt, v.EventHash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs atualizar evento em %v em nome de %v", fmtHandle(handle), fmtEvent(event.StartAt, v.EventHash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("atualização do evento em %v em nome de %v foi negada", fmtEvent(event.StartAt, v.EventHash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
			}
		}
	case *actions.CheckinEvent:
		if event, ok := i.state.Events[v.EventHash]; ok {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			return fmt.Sprintf("%v confirmou participação no evento em %v por %v ", fmtHandle(handle), fmtEvent(event.StartAt, v.EventHash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
		}
	case *actions.GreetCheckinEvent:
		if event, ok := i.state.Events[v.EventHash]; ok {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			return fmt.Sprintf("%v recebeu boas-vindas para evento em %v por %v ", fmtHandle(handle), fmtEvent(event.StartAt, v.EventHash), fmtCollective(event.Collective.Name)), v.Epoch, v.Reasons
		}
	case *actions.CreateBoard:
		boardhash := v.Hashed()
		if status == state.Favorable {
			if board, ok := i.state.Boards[boardhash]; ok {
				return fmt.Sprintf("mural %v criado em nome de %v", fmtBoard(board.Name), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
			}
		}
		if status == state.Undecided {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			return fmt.Sprintf("%v propôs criação de mural %v em nome de %v", fmtHandle(handle), fmtBoard(v.Name), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
		}
		if status == state.Against {
			return fmt.Sprintf("criação do mural %v em nome de %v foi negada", fmtBoard(v.Name), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
		}
	case *actions.UpdateBoard:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("mural %v atualizado por %v", fmtBoard(board.Name), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs atualização do mural %v em nome de %v", fmtHandle(handle), fmtBoard(board.Name), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("atualização do mural %v em nome de %v foi negada", fmtBoard(board.Name), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
			}
		}
	case *actions.Pin:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			if draft, ok := i.state.Drafts[v.Draft]; ok {
				pinaction := []string{"desafixado de", "desafixar de"}
				if v.Pin {
					pinaction = []string{"fixado em", "fixar em"}
				}
				if status == state.Favorable {
					return fmt.Sprintf("%v %v %v em nome de %v", fmtDraft(draft.Title, draft.DraftHash), pinaction[0], fmtBoard(board.Name), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
				}
				if status == state.Undecided {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v propôs %v %v o esboço %v em nome de %v", fmtHandle(handle), pinaction[1], fmtBoard(board.Name), fmtDraft(draft.Title, draft.DraftHash), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
				}
				if status == state.Against {
					return fmt.Sprintf("%v %v o esboço %v em nome de %v foi negado", pinaction[1], fmtBoard(board.Name), fmtDraft(draft.Title, draft.DraftHash), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
				}
			}
		}
	case *actions.BoardEditor:
		hash := crypto.Hasher([]byte(v.Board))
		if board, ok := i.state.Boards[hash]; ok {
			editor := i.state.Members[crypto.HashToken(v.Editor)]
			editorship := []string{"removido de", "remoção de", "de", "editor removal"}
			if v.Insert {
				editorship = []string{"incluído como", "inclusão de", "para", "editor inclusion"}
			}
			if status == state.Favorable {
				return fmt.Sprintf("%v %v editoria de %v em nome de %v", fmtHandle(editor), editorship[0], fmtBoard(board.Name), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs %v %v %v editoria de %v em nome de %v", fmtHandle(handle), editorship[1], fmtHandle(editor), editorship[2], fmtBoard(board.Name), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("%v %v %v editoria de %v em nome de %v foi negada", editorship[1], fmtHandle(editor), editorship[2], fmtBoard(board.Name), fmtCollective(board.Collective.Name)), v.Epoch, v.Reasons
			}
		}
	case *actions.Draft:
		if status == state.Favorable {
			if draft, ok := i.state.Drafts[v.ContentHash]; ok {
				return fmt.Sprintf("esboço %v foi criado por %v", fmtDraft(draft.Title, draft.DraftHash), fmtAuthors(draft.Authors, i.state)), v.Epoch, v.Reasons
			}
		}
		if status == state.Undecided {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			if v.OnBehalfOf != "" {
				return fmt.Sprintf("%v propôs esboço %v em nome de %v", fmtHandle(handle), fmtDraft(v.Title, v.ContentHash), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
			}
			if len(v.CoAuthors) > 0 {
				// saida := ""
				// for _, aut := range v.CoAuthors {
				// 	autHandle := i.state.Members[crypto.HashToken(aut)]

				// }
				return fmt.Sprintf("%v propôs esboço %v em co-autoria", fmtHandle(handle), fmtDraft(v.Title, v.ContentHash)), v.Epoch, v.Reasons
			}
		}
		if status == state.Against {
			if v.OnBehalfOf != "" {
				return fmt.Sprintf("proposta de esboço %v em nome de %v foi negada", fmtDraft(v.Title, v.ContentHash), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
			}
			if len(v.CoAuthors) > 0 {
				return fmt.Sprintf("proposta de esboço %v em coautoria foi negada", fmtDraft(v.Title, v.ContentHash)), v.Epoch, v.Reasons
			}
		}
	case *actions.ReleaseDraft:
		if draft, ok := i.state.Drafts[v.ContentHash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("esboço %v foi lançado por %v", fmtDraft(draft.Title, draft.DraftHash), fmtAuthors(draft.Authors, i.state)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs lançamento de %v em nome de %v", fmtHandle(handle), fmtDraft(draft.Title, draft.DraftHash), fmtAuthors(draft.Authors, i.state)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("lançamento de %v em nome de %v foi negado", fmtDraft(draft.Title, draft.DraftHash), fmtAuthors(draft.Authors, i.state)), v.Epoch, v.Reasons
			}
		}
	case *actions.Edit:
		if status == state.Favorable {
			if edit, ok := i.state.Edits[v.ContentHash]; ok {
				draft := i.state.Drafts[v.EditedDraft]
				if edit.Authors.CollectiveName() != "" {
					if status == state.Favorable {
						return fmt.Sprintf("edição de %v foi proposta em nome de %v", fmtDraft(draft.Title, draft.DraftHash), fmtAuthors(edit.Authors, i.state)), v.Epoch, v.Reasons
					}
				}
				return fmt.Sprintf("edição de %v foi proposta", fmtDraft(draft.Title, draft.DraftHash)), v.Epoch, v.Reasons
			}
		}
		if status == state.Undecided {
			if edit, ok := i.state.Proposals.Edit[v.ContentHash]; ok {
				draft := i.state.Drafts[v.EditedDraft]
				if edit.Authors.CollectiveName() != "" {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					return fmt.Sprintf("%v propôs edição de %v em nome de %v", fmtHandle(handle), fmtDraft(draft.Title, draft.DraftHash), fmtAuthors(edit.Authors, i.state)), v.Epoch, v.Reasons
				}
				if len(v.CoAuthors) > 0 {
					handle := i.state.Members[crypto.HashToken(v.Author)]
					draft := i.state.Drafts[v.EditedDraft]
					return fmt.Sprintf("%v propôs edição de %v em co-autoria", fmtHandle(handle), fmtDraft(draft.Title, draft.DraftHash)), v.Epoch, v.Reasons
				}
			}

		}
		if status == state.Against {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			draft := i.state.Drafts[v.EditedDraft]
			if v.OnBehalfOf != "" {
				return fmt.Sprintf("edição de %v em nome de %v foi negada", fmtDraft(draft.Title, draft.DraftHash), fmtCollective(v.OnBehalfOf)), v.Epoch, v.Reasons
			}
			return fmt.Sprintf("edição proposta para %v foi por %v negada", fmtDraft(draft.Title, draft.DraftHash), fmtHandle(handle)), v.Epoch, v.Reasons
		}
	case *actions.React:
		handle := i.state.Members[crypto.HashToken(v.Author)]
		return fmt.Sprintf("%v reagiu", fmtHandle(handle)), v.Epoch, v.Reasons
	case *actions.CreateCollective:
		collectivehash := crypto.Hasher([]byte(v.Name))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v foi criado", fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs criação de %v", fmtHandle(handle), fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
		}
	case *actions.UpdateCollective:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v atualizado", fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v propôs atualização para %v", fmtHandle(handle), fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("atualização para %v foi negada", fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
		}
	case *actions.RequestMembership:
		collectivehash := crypto.Hasher([]byte(v.Collective))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			handle := i.state.Members[crypto.HashToken(v.Author)]
			if v.Include {
				if status == state.Favorable {
					return fmt.Sprintf("%v tornou-se membro do coletivo %v", fmtHandle(handle), fmtCollective(collective.Name)), v.Epoch, v.Reasons
				}
				if status == state.Undecided {
					return fmt.Sprintf("%v pediu para ingressar no coletivo %v", fmtHandle(handle), fmtCollective(collective.Name)), v.Epoch, v.Reasons
				}
				if status == state.Against {
					return fmt.Sprintf("ingresso de %v ao coletivo %v foi negado", fmtHandle(handle), fmtCollective(collective.Name)), v.Epoch, v.Reasons
				}
			}
			return fmt.Sprintf("%v left %v", fmtHandle(handle), fmtCollective(collective.Name)), v.Epoch, v.Reasons
		}
	case *actions.RemoveMember:
		collectivehash := crypto.Hasher([]byte(v.OnBehalfOf))
		if collective, ok := i.state.Collectives[collectivehash]; ok {
			member := i.state.Members[crypto.HashToken(v.Member)]
			if status == state.Favorable {
				return fmt.Sprintf("%v foi removido de %v", fmtHandle(member), fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Undecided {
				handle := i.state.Members[crypto.HashToken(v.Author)]
				return fmt.Sprintf("%v prediu a remoção de %v de %v", fmtHandle(handle), fmtHandle(member), fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
			if status == state.Against {
				return fmt.Sprintf("remoção de %v do coletivo %v foi negada", fmtHandle(member), fmtCollective(collective.Name)), v.Epoch, v.Reasons
			}
		}
	case *actions.Signin:
		authorhash := crypto.HashToken(v.Author)
		if _, ok := i.state.Members[authorhash]; ok {
			if status == state.Favorable {
				return fmt.Sprintf("%v ingressou no Motiró"), v.Epoch, v.Reasons
			}
		}
	}
	return "", 0, ""
}
