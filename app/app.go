package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/configs"
	"github.com/muchlist/mini_pos/db"
	_ "github.com/muchlist/mini_pos/docs"
	"github.com/muchlist/mini_pos/utils/logger"
	"log"
	"os"
	"os/signal"
)

// RunApp
// @title mini_pos API
// @version 1.0
// @description Mini Pos Api
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email whois.muchlis@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
// @host localhost:3500
// @BasePath /api/v1
func RunApp() {
	// Init config, logger dan db
	configs.InitConfig()
	logger.InitLogger()
	db.Init()
	defer db.Close()

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

	prepareEndPoint(app)

	// blocking and listen for fiber
	if err := app.Listen(":3500"); err != nil {
		logger.Error("error fiber listen", err)
		log.Panic()
	}

	// cleanup app
	fmt.Println("Running cleanup tasks...")
}
