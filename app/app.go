package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/configs"
	"github.com/muchlist/mini_pos/utils/logger"
	"log"
	"os"
	"os/signal"
)

func RunApp() {
	// Init config, logger dan db
	configs.InitConfig()
	logger.InitLogger()

	// membuat fiber app
	app := fiber.New()

	// gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// <--- prepareEndPoint disini

	// blocking and listen for fiber
	if err := app.Listen(":3500"); err != nil {
		logger.Error("error fiber listen", err)
		log.Panic()
	}

	// cleanup app
	fmt.Println("Running cleanup tasks...")
}
