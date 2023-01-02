package controllers

import (
	"context"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/config"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func GetProfiles(c *gin.Context) {
	coll := database.DB.Database(config.Conf.Mongo.DBName).Collection("profiles")
	filter := bson.M{}
	filter["geolocation"] = bson.D{{"$exists", true}}
	opts := options.Find().SetProjection(bson.D{{"geolocation", 1}, {"primary_url", 1}})

	getParams := c.Request.URL.Query()
	if getParams.Has("schema") {
		schema := getParams.Get("schema")
		schemaArray := [1]string{schema}
		filter["linked_schemas"] = bson.D{{"$in", schemaArray}}
	}
	if getParams.Has("tags") {
		tags := strings.Split(getParams.Get("tags"), ",")
		filter["tags"] = bson.D{{"$all", tags}}
	}

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
