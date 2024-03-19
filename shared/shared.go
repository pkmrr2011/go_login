package shared

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleError(c *gin.Context, statusCode int, message string, err error) {
	response := gin.H{"error": message}
	if err != nil {
		response["details"] = err.Error()
	}
	c.JSON(statusCode, response)
	c.Abort()
}

func GetObjectIDFromHex(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
