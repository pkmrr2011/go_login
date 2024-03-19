package controllers

import (
	"login/middleware"
	"login/model"
	"login/shared"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type adminController struct{}

func NewAdminController() *adminController {
	return &adminController{}
}

func (uc *adminController) AdminLogin(c *gin.Context) {
	var auth model.Auth
	if err := c.ShouldBindJSON(&auth); err != nil {
		shared.HandleError(c, http.StatusBadRequest, "Error binding JSON", err)
		return
	}

	envEmail := os.Getenv("ADMIN_EMAIL")
	envPassword := os.Getenv("ADMIN_PASSWORD")

	if auth.Email == envEmail && auth.Password == envPassword {
		token, err := middleware.CreateToken(auth.Email)
		if err != nil {
			shared.HandleError(c, http.StatusInternalServerError, "Error creting token", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token})
		return
	} else {
		shared.HandleError(c, http.StatusUnauthorized, "Login failed", nil)
		return
	}
}
