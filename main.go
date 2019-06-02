package main

import (
	"flag"
	"keysharer/controllers"
	"keysharer/models"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

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

	us, err := models.NewUserService(vConfig.GetString("database.type"), vConfig.GetString("database.args"))
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	usersC := controllers.NewUsers(us)

	web := echo.New()
	web.Use(middleware.Logger())
	web.Use(middleware.Recover())

	// setup handlers
	web.GET("/user/:username", usersC.GetUser)
	web.POST("/user/:username/:email", usersC.CreateUser)
	web.PUT("/user/:username/:email", usersC.UpdateUser)
	web.DELETE("/user/:username", usersC.DeleteUser)

	log.Fatal(web.Start(vConfig.GetString("server.port")))
}
