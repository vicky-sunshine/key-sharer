package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"keysharer/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// API
type API struct {
	config *viper.Viper // TODO: should config be a concrete struct for better?
	web    *echo.Echo
	db     *gorm.DB
	// TODO: maybe later add logger
}

func main() {

	// set flag
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "config file path")
	flag.Parse()

	// read environment
	vConfig := viper.New()
	vConfig.SetConfigType("yaml")
	f, err := os.Open(cfgPath)
	if err != nil {
		panic("config file not found")
	}
	err = vConfig.ReadConfig(f)
	if err != nil {
		panic("config file parse fail")
	}

	// https://github.com/jinzhu/gorm/issues/1427#issuecomment-332498453
	db, err := gorm.Open(vConfig.GetString("database.type"), vConfig.GetString("database.args"))
	if err != nil {
		panic("DB connect error")
	}

	// create api instance
	api := NewAPI(vConfig, db)

	// setup handlers
	log.Fatal(api.web.Start(api.config.GetString("server.port")))
}

func NewAPI(v *viper.Viper, db *gorm.DB) *API {
	a := &API{config: v, web: echo.New(), db: db}
	a.initMigration()
	a.registerRoute()
	return a
}

// db init
func (a *API) initMigration() {
	a.db.AutoMigrate(&models.User{})
}

// route dispather
func (a *API) registerRoute() {
	a.web.Use(middleware.Logger())
	a.web.Use(middleware.Recover())

	a.web.GET("/users", a.allUsers)
	a.web.GET("/user/:username", a.getUser)
	a.web.POST("/user/:username/:email", a.newUser)
	a.web.PUT("/user/:username/:email", a.updateUser)
	a.web.DELETE("/user/:username", a.deleteUser)

}

// TODO: should sperate controllers & models?
func (a *API) allUsers(c echo.Context) error {
	var users []models.User
	a.db.Find(&users)
	return c.JSON(http.StatusOK, users)
}

func (a *API) getUser(c echo.Context) error {
	username := c.Param("username")

	var user models.User
	a.db.Where("username=?", username).Find(&user)
	return c.JSON(http.StatusOK, user)
}

func (a *API) newUser(c echo.Context) error {
	username := c.Param("username")
	email := c.Param("email")

	a.db.Create(&models.User{Username: username, Email: email})
	return c.NoContent(http.StatusOK)

}

func (a *API) updateUser(c echo.Context) error {
	username := c.Param("username")
	email := c.Param("email")

	var user models.User
	a.db.Where("username=?", username).Find(&user)
	user.Email = email
	a.db.Save(&user)

	return c.NoContent(http.StatusOK)

}

func (a *API) deleteUser(c echo.Context) error {
	username := c.Param("username")

	var user models.User
	a.db.Where("username=?", username).Find(&user)
	a.db.Delete(&user)
	return c.String(http.StatusOK, "ok")
}
