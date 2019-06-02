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
	username := c.Param("username")
	email := c.Param("email")
	user := models.User{Username: username, Email: email}

	if err := u.us.Create(&user); err != nil {
		c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

func (u *Users) GetUser(c echo.Context) error {
	username := c.Param("username")

	user, err := u.us.ByUsername(username)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, user)
}

func (u *Users) UpdateUser(c echo.Context) error {
	username := c.Param("username")
	email := c.Param("email")

	user := models.User{Username: username, Email: email}
	err := u.us.Update(&user)

	if err != nil {
		c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, user)
}

func (u *Users) DeleteUser(c echo.Context) error {
	username := c.Param("username")

	err := u.us.Delete(username)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
