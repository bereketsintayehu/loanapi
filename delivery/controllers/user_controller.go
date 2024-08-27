package controllers

import (
	"loan/domain"
	"loan/infrastracture"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase domain.UserUsecase
}

func NewUserController(userUsecase domain.UserUsecase) *UserController {
	return &UserController{UserUsecase: userUsecase}
}

func (uc *UserController) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceFingerprint := infrastracture.GenerateDeviceFingerprint(ipAddress, userAgent)

	LogInResponse, uerr := uc.UserUsecase.Login(&user, deviceFingerprint)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	c.JSON(200, gin.H{"tokens": LogInResponse})

}

func (uc *UserController) RefreshToken(c *gin.Context) {
	var refreshRequest domain.RefreshTokenRequest

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceFingerprint := infrastracture.GenerateDeviceFingerprint(ipAddress, userAgent)

	refreshResponse, uerr := uc.UserUsecase.RefreshToken(refreshRequest.UserID, deviceFingerprint, refreshRequest.Token)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	c.JSON(200, gin.H{"tokens": refreshResponse})
}

func (uc *UserController) GetNewVerificationEmail(c *gin.Context) {
	var req domain.NewEmailVerification
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	uerr := uc.UserUsecase.GetNewVerificationEmail(req.Email)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	c.JSON(200, gin.H{"message": "Verification email sent. Please check your email"})
}

func (uc *UserController) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := uc.UserUsecase.Register(user)
	if err.Message != "" {
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, gin.H{"message": "Registered successfully. Please check your email for account activation."})

}

func (uc *UserController) ActivateAccount(c *gin.Context) {
	token := c.Param("token")
	email := c.Param("email")

	err := uc.UserUsecase.ActivateAccount(token, email)
	if err.Message != "" {
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, gin.H{"message": "Account activated successfully"})
}

func (uc *UserController) GetMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, uerr := uc.UserUsecase.GetMyProfile(userID)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	c.JSON(200, gin.H{
		"id":        user.ID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"name":      user.Name,
		"bio":       user.Bio,
		"role":      user.Role,
		"is_active": user.IsActive,
	})
}

func (uc *UserController) PasswordReset(c *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	uerr := uc.UserUsecase.PasswordReset(req.Email)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "message": "Successfully sent password reset link to your email"})
}

func (uc *UserController) UpdatePassword(c *gin.Context) {
	var req domain.PasswordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	uerr := uc.UserUsecase.UpdatePassword(req)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset"})
}
