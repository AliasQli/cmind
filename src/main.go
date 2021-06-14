package main

import (
	"conWord/src/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/related", controller.GetSingleWordRelated)
	log.Fatal(e.Start("0.0.0.0:3000"))

}
