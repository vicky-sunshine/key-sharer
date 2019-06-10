package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func NewStatic() *Static {
	return &Static{}
}

type Static struct{}

func (s *Static) Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home", nil)
}
