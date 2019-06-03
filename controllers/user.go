package controllers

import (
	"keysharer/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Users struct {
	us *models.UserService
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		us: us,
	}
}

func (u *Users) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := u.us.Create(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Login is used to process the login form when a user
// tries to log in as an existing user with username & password
func (u *Users) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginedUser, err := u.us.Authenticate(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loginedUser)
}
