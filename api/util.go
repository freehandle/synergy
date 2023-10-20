package api

/*type ViewerWithHash func(*state.State, *index.Index, crypto.Hash, crypto.Token) any

type ViewerWithString func(*state.State, *index.Index, crypto.Token, string) any

type Viewer func(*state.State, *index.Index, crypto.Token) any

func (a *AttorneyGeneral) HandleWithHash(template, path string, viewer ViewerWithHash) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		hash := getHash(r.URL.Path, path)
		author := a.Author(r)
		view := viewer(a.state, a.indexer, hash, author)
		if err := a.templates.ExecuteTemplate(w, template, view); err != nil {
			log.Println(err)
		}
	}
	a.mux.HandleFunc(path, handler)
}

func (a *AttorneyGeneral) HandleWithString(template, path string, viewer ViewerWithString) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		item := strings.Replace(r.URL.Path, path, "", 1)
		author := a.Author(r)
		view := viewer(a.state, a.indexer, author, item)
		if err := a.templates.ExecuteTemplate(w, template, view); err != nil {
			log.Println(err)
		}
	}
	a.mux.HandleFunc(path, handler)
}

func (a *AttorneyGeneral) HandleSimple(template, path string, viewer Viewer) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		author := a.Author(r)
		view := viewer(a.state, a.indexer, author)
		if err := a.templates.ExecuteTemplate(w, template, view); err != nil {
			log.Println(err)
		}
	}
	a.mux.HandleFunc(path, handler)
}
*/
