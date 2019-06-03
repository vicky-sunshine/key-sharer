package main

import (
	"flag"
	"keysharer/controllers"
	"keysharer/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	us, err := models.NewUserService(
		vConfig.GetString("database.type"),
		vConfig.GetString("database.args"),
		vConfig.GetString("database.pepper"),
	)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	usersC := controllers.NewUsers(us)

	web := gin.Default()
	web.Use(gin.Logger())
	web.Use(gin.Recovery())

	// setup handlers
	web.POST("/user", usersC.CreateUser)
	web.POST("/login", usersC.Login)

	log.Fatal(web.Run(vConfig.GetString("server.port")))
}
