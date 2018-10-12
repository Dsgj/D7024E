package api

import (
	"D7024E/dht"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
)

func Routes(r chi.Router, k *dht.Kademlia) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/files/{key}", func(w http.ResponseWriter, r *http.Request) {
		fetchFile(w, r, k)
	})
	r.Post("/files", func(w http.ResponseWriter, r *http.Request) {
		storeFile(w, r, k)
	})
	r.Patch("/files/{key}/pin", func(w http.ResponseWriter, r *http.Request) {
		pinFile(w, r, k)
	})
	r.Patch("/files/{key}/unpin", func(w http.ResponseWriter, r *http.Request) {
		unpinFile(w, r, k)
	})
}

func fetchFile(w http.ResponseWriter, r *http.Request, k *dht.Kademlia) {
	key := chi.URLParam(r, "key")
	data := k.FetchFile(key)
	if data == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	w.Write(data)
}

func storeFile(w http.ResponseWriter, r *http.Request, k *dht.Kademlia) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	key := k.StoreFile(body)
	w.Write([]byte(key))
}

func pinFile(w http.ResponseWriter, r *http.Request, k *dht.Kademlia) {
	w.Write([]byte("pin file"))
}

func unpinFile(w http.ResponseWriter, r *http.Request, k *dht.Kademlia) {
	w.Write([]byte("unpin file"))
}
