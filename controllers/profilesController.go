package controllers

import (
	"MurmurationsAllocator/config"
	"MurmurationsAllocator/database"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProfiles(c *gin.Context) {
	coll := database.DB.Database(config.Conf.Mongo.DBName).Collection("profiles")
	filter := bson.D{{"geolocation", bson.D{{"$exists", true}}}}
	opts := options.Find().SetProjection(bson.D{{"geolocation", 1}, {"primary_url", 1}})
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	var profiles []map[string]interface{}
	if err = cursor.All(context.TODO(), &profiles); err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	c.JSON(200, profiles)
}
