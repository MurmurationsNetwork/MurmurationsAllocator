package database

import (
	"context"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var DB *mongo.Client

func ConnectMongo() {
	var err error
	credential := options.Credential{
		Username: config.Conf.Mongo.USERNAME,
		Password: config.Conf.Mongo.PASSWORD,
	}
	clientOptions := options.Client().ApplyURI(config.Conf.Mongo.HOST).SetAuth(credential)
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
