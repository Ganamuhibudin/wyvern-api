package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wyvern-api/config"
	"wyvern-api/routers"
)

func main() {
	config.LoadConfig() // load config files
	config.LoadDB()     // initiate database connection

	r := gin.Default()
	routers.Routes(r) // added all routes

	r.Run(fmt.Sprintf(":%s", config.ENV.Port))
}
