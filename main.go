package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/guythatdrinkscoffee/ShortyURL/controllers"
	"github.com/guythatdrinkscoffee/ShortyURL/internal"
	"github.com/guythatdrinkscoffee/ShortyURL/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting server application...")

	connectionCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/shorty"))

	if err != nil {
		log.Fatalln("Failed to initialize a MongoDB Client")
	}

	err = client.Connect(connectionCtx)

	if err != nil {
		log.Fatalln("Failed to connect to the database...")
	}

	defer func() {
		if err = client.Disconnect(connectionCtx); err != nil {
			log.Fatalln("Failed to disconnect from the database...")
		}
	}()

	collection := client.Database("shorty").Collection("urls")

	db := internal.NewDB(collection)
	repo := repository.NewRepository(db)
	mapperController := controllers.NewMapper(repo)
	redirectController := controllers.NewRedirect(repo)

	//Mux Config
	router := mux.NewRouter()

	mappingRoute := router.Methods(http.MethodPost).Subrouter()
	mappingRoute.HandleFunc("/map", mapperController.MapURL)
	mappingRoute.Use(mapperController.ValidateURL)

	redirectRoute := router.Methods(http.MethodGet).Subrouter()
	redirectRoute.HandleFunc("/{hash}", redirectController.Redirect)

	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln("Failed to start server...")
	}
}
