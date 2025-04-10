package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ExtraWhy/internal-libs/models/player"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

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
	//	pl := player.Player{Name: "Lubaka", Id: 1, Money: 9999}

	//	result, err := coll.InsertOne(context.TODO(), pl)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("Inserted document with ID: %v", result.InsertedID)

	var result player.Player
	err = coll.FindOne(context.TODO(), bson.M{"id": 1}).Decode(&result)
	if err != nil {
		fmt.Printf("FindOne failed: %v", err)
		return
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found document: %s", jsonData)

	result.Money *= 2
	updt := bson.M{"$set": bson.M{"money": result.Money}}
	_, err = coll.UpdateOne(context.TODO(), bson.M{"id": result.Id}, updt)
	if err != nil {
		fmt.Printf("UpdateOne failed: %v", err)
		return
	}

	err = coll.FindOne(context.TODO(), bson.M{"id": 1}).Decode(&result)
	if err != nil {
		fmt.Printf("FindOne failed: %v", err)
		return
	}
	jsonData, err = json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Update document: %s", jsonData)
	//find all

	cursor, err := coll.Find(context.TODO(), bson.M{})
	for cursor.Next(context.TODO()) {
		var elem player.Player
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(elem)

	}

}
