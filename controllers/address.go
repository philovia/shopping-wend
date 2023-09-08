package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mreym/shopping/models"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "addres_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "internal server error")
		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}

		} else {
			c.IndentedJSON(400, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}

}

func EditHomeAddress(UserCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/type")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "internal server error")
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House},
			{Key: "address.0.street_name", Value: editaddress.Street},
			{Key: "address.0.city_name", Value: editaddress.City},
			{Key: "address.0.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully update the home address")

	}
}

func EditWorkAddress(UserCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			return
		}
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userObjectID}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "address.1.house_name", Value: editAddress.House},
				{Key: "address.1.street_name", Value: editAddress.Street},
				{Key: "address.1.city_name", Value: editAddress.City},
				{Key: "address.1.pin_code", Value: editAddress.Pincode},
			}},
		}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something Went Wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully updated the work address"})
	}
}

// func EditWorkAddress() gin.HandlerFunc {

// 	return func(c *gin.Context) {
// 		user_id := c.Query("id")
// 		if user_id == "" {
// 			c.Header("Content-Type", "application/type")
// 			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
// 			c.Abort()
// 			return
// 		}
// 		user_id, err := primitive.ObjectIDFromHex(user_id)
// 		if err != nil {
// 			c.IndentedJSON(500, "internal server error")
// 		}
// 		var editaddress models.Address
// 		if err := c.BindJSON(&editaddress); err != nil {
// 			c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
// 		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House},
// 			{Key: "adress.1.street_name", Value: editaddress.Street},
// 			{Key: "adress.1.city_name", Value: editaddress.City},
// 			{Key: "address.1.pin_code", Value: editaddress.Pincode}}}}
// 		_, err = UserCollection.UpdateOne(ctx, filter, update)
// 		if err != nil {
// 			c.IndentedJSON(500, "Something Wrong")
// 			return
// 		}
// 		defer cancel()
// 		ctx.Done()
// 		c.IndentedJSON(200, "Successfully updated the work address")
// 	}

// }

func DeleteAddress(UserCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Search Index"})
			return
		}

		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userObjectID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: []models.Address{}}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wrong Command"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted"})
	}
}
