package controllers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/guythatdrinkscoffee/ShortyURL/repository"
	"net/http"
	"time"
)

type Redirect struct {
	Repository repository.Repository
}

func NewRedirect(repo repository.Repository) Redirect {
	return Redirect{repo}
}

func (r *Redirect) Redirect(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	hash, ok := vars["hash"]

	if !ok {
		http.Error(writer, "The link was not found", http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	shortUrl, err := r.Repository.Find(ctx, hash)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(writer, request, shortUrl.OriginalURL, http.StatusTemporaryRedirect)
}
