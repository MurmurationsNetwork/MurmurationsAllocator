package database

import (
	"context"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"time"
)

var DB *mongo.Client

func ConnectMongo() {
	username := url.QueryEscape(config.Conf.Mongo.USERNAME)
	password := url.QueryEscape(config.Conf.Mongo.PASSWORD)
	mongoURI := "mongodb://" + username + ":" + password + "@" + config.Conf.Mongo.HOST + "/?authSource=admin&tls=true"
	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	DB, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	err = DB.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	fmt.Println("Connected to MongoDB!")
}

func DisconnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := DB.Disconnect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}
	fmt.Println("Disconnected from MongoDB!")
}
