package controllers

import (
	"loan/domain"

	"github.com/gin-gonic/gin"
)

func (uc *UserController) GetUsers(c *gin.Context) {
	Role := c.GetString("role")
	if Role != "admin" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	users, uerr := uc.UserUsecase.GetUsers()
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	var user domain.User
	Role := c.GetString("role")
	userID := c.Param("id")
	user, uerr := uc.UserUsecase.GetMyProfile(userID)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	checkUserID := user.ID.Hex()

	if Role != "admin" || userID != checkUserID {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	deletedUser, uerr := uc.UserUsecase.DeleteUser(userID)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully", "user": deletedUser})
}
