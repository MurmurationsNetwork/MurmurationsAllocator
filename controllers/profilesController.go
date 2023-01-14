package controllers

import (
	"context"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/config"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
)

func GetProfile(c *gin.Context) {
	profileUrl, exist := c.GetQuery("profile_url")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "profile_url is missing",
		})
		return
	}

	coll := database.DB.Database(config.Conf.Mongo.DBName).Collection("profiles")
	filter := bson.D{{"profile_url", profileUrl}}
	opts := options.FindOne().SetProjection(bson.M{"_id": 0})

	var profile map[string]interface{}
	err := coll.FindOne(context.TODO(), filter, opts).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func GetProfiles(c *gin.Context) {
	coll := database.DB.Database(config.Conf.Mongo.DBName).Collection("profiles")
	filter := bson.M{}
	filter["geolocation"] = bson.D{{"$exists", true}}
	opts := options.Find().SetProjection(bson.D{{"_id", 1}, {"geolocation", 1}, {"profile_url", 1}})

	getParams := c.Request.URL.Query()
	if getParams.Has("schema") {
		schema := getParams.Get("schema")
		schemaArray := [1]string{schema}
		filter["linked_schemas"] = bson.D{{"$in", schemaArray}}
	}
	if getParams.Has("primary_url") {
		filter["primary_url"] = getParams.Get("primary_url")
	}
	if getParams.Has("last_updated") {
		lastUpdated := getParams.Get("last_updated")
		lastUpdatedInt, err := strconv.Atoi(lastUpdated)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		filter["last_updated"] = bson.D{{"$gte", lastUpdatedInt}}
	}
	if getParams.Has("status") {
		filter["status"] = getParams.Get("status")
	}

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var profiles []map[string]interface{}
	if err = cursor.All(context.TODO(), &profiles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if len(profiles) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
		return
	}

	// [lon, lat, profile_url]
	queryResults := make([][]interface{}, 0)
	for _, profile := range profiles {
		geolocation := profile["geolocation"].(map[string]interface{})
		mapResult := []interface{}{geolocation["lon"], geolocation["lat"], profile["profile_url"]}
		queryResults = append(queryResults, mapResult)
	}

	c.JSON(http.StatusOK, gin.H{"data": queryResults})
}
