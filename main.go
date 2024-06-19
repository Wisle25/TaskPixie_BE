package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/server"
)

func main() {
	config := commons.LoadConfig(".")
	app := server.CreateServer(config)

	log.Fatal(
		app.Listen(fmt.Sprintf(":%s", config.ServerPort)),
	)
}
