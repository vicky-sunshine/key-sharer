package controllers

import (
	"keysharer/models"
	"net/http"

	"github.com/labstack/echo"
)

type Users struct {
	us *models.UserService
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		us: us,
	}
}

func (u *Users) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	if err := u.us.Create(&user); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// Login is used to process the login form when a user
// tries to log in as an existing user with username & password
func (u *Users) Login(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	loginedUser, err := u.us.Authenticate(user.Username, user.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, loginedUser)
}
