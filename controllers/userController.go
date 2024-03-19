package controllers

import (
	"context"
	"login/model"
	"login/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userController struct {
	DB *mongo.Database
}

func NewUserController(db *mongo.Database) *userController {
	return &userController{DB: db}
}

func (uc *userController) GetUsers(c *gin.Context) {

	var users []model.User
	cursor, err := uc.DB.Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		shared.HandleError(c, http.StatusInternalServerError, "Error fetching users", err)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			shared.HandleError(c, http.StatusInternalServerError, "Error decoding user", err)
			return
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No users found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func (uc *userController) GetUser(c *gin.Context) {
	id := c.Param("id")
	objectID, err := shared.GetObjectIDFromHex(id)
	if err != nil {
		shared.HandleError(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var user model.User
	err = uc.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			shared.HandleError(c, http.StatusNotFound, "User not found", err)
		} else {
			shared.HandleError(c, http.StatusInternalServerError, "Error fetching user", err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *userController) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		shared.HandleError(c, http.StatusBadRequest, "Error binding JSON", err)
		return
	}

	insertResult, err := uc.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		shared.HandleError(c, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	user.ID = insertResult.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

func (uc *userController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	objectID, err := shared.GetObjectIDFromHex(id)
	if err != nil {
		shared.HandleError(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		shared.HandleError(c, http.StatusBadRequest, "Error binding JSON", err)
		return
	}

	updateResult, err := uc.DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": user})
	if err != nil {
		shared.HandleError(c, http.StatusInternalServerError, "Error updating user", err)
		return
	}
	if updateResult.ModifiedCount == 0 {
		shared.HandleError(c, http.StatusNotFound, "No changes in user", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func (uc *userController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objectID, err := shared.GetObjectIDFromHex(id)
	if err != nil {
		shared.HandleError(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	deleteResult, err := uc.DB.Collection("users").DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		shared.HandleError(c, http.StatusInternalServerError, "Error deleting user", err)
		return
	}
	if deleteResult.DeletedCount == 0 {
		shared.HandleError(c, http.StatusNotFound, "User not found", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
