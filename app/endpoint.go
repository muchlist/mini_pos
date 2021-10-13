package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/muchlist/mini_pos/dao/merchant_dao"
	"github.com/muchlist/mini_pos/db"
	"github.com/muchlist/mini_pos/handler"
	"github.com/muchlist/mini_pos/service/merchant_serv"
	"github.com/muchlist/mini_pos/utils/mcrypt"
)

func prepareEndPoint(app *fiber.App) {

	// Utils
	cryptoUtils := mcrypt.NewCrypto()

	// Merchant Domain
	merchantDao := merchant_dao.New(db.DB)
	merchantService := merchant_serv.NewMerchantService(merchantDao, cryptoUtils)
	merchantHandler := handler.NewMerchantHandler(merchantService)

	app.Use(logger.New())

	// url mapping
	api := app.Group("/api/v1")

	// Merchant Endpoint
	api.Post("/merchant", merchantHandler.CreateMerchant)
	api.Get("/merchant/:id", merchantHandler.GetMerchant)
	api.Get("/merchant", merchantHandler.FindMerchant)
	api.Put("/merchant/:id", merchantHandler.EditMerchant)
	api.Delete("/merchant/:id", merchantHandler.DeleteMerchant)
}
