package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	log.Fatal(e.Start(":8010"))
}
