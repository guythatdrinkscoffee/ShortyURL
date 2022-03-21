package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guythatdrinkscoffee/ShortyURL/models"
	"github.com/guythatdrinkscoffee/ShortyURL/repository"
	"net/http"
	"net/url"
	"time"
)

type Mapper struct {
	Repository repository.Repository
}

func NewMapper(repo repository.Repository) Mapper {
	return Mapper{repo}
}

type MapperKey struct{}

func (m *Mapper) ValidateURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		queryParams := request.URL.Query()
		urlQuery := queryParams.Get("url")
		aliasQuery := queryParams.Get("alias")

		//Check if the passed in url is valid
		newUrl, err := isValidUrl(urlQuery)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		//Build a payload
		payload := &models.RequestPayload{
			Url:   newUrl.String(),
			Alias: aliasQuery,
		}

		//Pass the payload through the context
		ctx := context.WithValue(request.Context(), MapperKey{}, payload)
		request = request.WithContext(ctx)

		next.ServeHTTP(writer, request)
	})
}

func (m *Mapper) MapURL(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	payload := request.Context().Value(MapperKey{}).(*models.RequestPayload)

	combined := payload.Url + payload.Alias

	key := fmt.Sprintf("%x", sha256.Sum256([]byte(combined)))[:6]
	createdAt := time.Now()

	shortUrl := &models.ShortURL{
		Hash:           key,
		OriginalURL:    payload.Url,
		DateCreated:    createdAt,
		ExpirationDate: time.Now(),
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := m.Repository.CreateURL(ctx, shortUrl)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.ResponsePayload{
		"localhost:8080/" + key,
	}

	json.NewEncoder(writer).Encode(res)
}

func isValidUrl(urlString string) (*url.URL, error) {
	validUrl, err := url.Parse(urlString)

	if err != nil {
		return nil, err
	}

	if validUrl.Scheme == "" || validUrl.Host == "" || validUrl.Path == "" {
		return nil, errors.New("invalid url")
	}

	return validUrl, nil
}
