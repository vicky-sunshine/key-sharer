package main

import (
	"flag"
	"keysharer/controllers"
	"keysharer/models"
	"keysharer/views"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// TemplateRenderer is a custom html/template renderer for Echo framework

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
		log.Fatalf("config file not found", err)
	}
	err = vConfig.ReadConfig(f)
	if err != nil {
		log.Fatalf("config file parse fail", err)
	}

	us, err := models.NewUserService(
		vConfig.GetString("database.type"),
		vConfig.GetString("database.args"),
		vConfig.GetString("database.pepper"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer us.Close()
	err = us.AutoMigrate()
	if err != nil {
		log.Fatalf("user service migrate fail", err)
	}

	usersC := controllers.NewUsers(us)
	staticC := controllers.NewStatic()

	web := echo.New()
	web.Use(middleware.Logger())
	web.Use(middleware.Recover())

	web.Renderer = views.NewTemplateRenderer("views/layouts/*.tmpl")

	// setup handlers
	web.GET("/", staticC.Home)
	web.POST("/user", usersC.CreateUser)
	web.POST("/login", usersC.Login)

	log.Fatal(web.Start(vConfig.GetString("server.port")))
}
