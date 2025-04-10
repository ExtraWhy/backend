package main

import (
	"context"
	"fmt"

	"github.com/ExtraWhy/internal-libs/models/player"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var names = []string{"Kucheto", "Kalniq",
	"Bonbonev", "Extramena", "Shto?",
	"Donev", "Gandalf", "Krasena", "Depresiyqta",
	"Mimito debelata", "Mishkata"}

func main() {

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://cryptowincryptowin:EfK0weUUe7t99Djx@cluster0.w07rcmn.mongodb.net/?appName=Cluster0").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	coll := client.Database("cryptowin").Collection("players")

	//CREATE
	for i := 1; i < 11; i++ {
		pl := player.Player{Name: names[i], Id: uint64(i), Money: 9999}

		result, err := coll.InsertOne(context.TODO(), pl)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Inserted document with ID: %v", result.InsertedID)
	}

}
