package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/muchlist/mini_pos/configs/roles"
	"github.com/muchlist/mini_pos/dao/merchant_dao"
	"github.com/muchlist/mini_pos/dao/user_dao"
	"github.com/muchlist/mini_pos/db"
	"github.com/muchlist/mini_pos/handler"
	"github.com/muchlist/mini_pos/middleware"
	"github.com/muchlist/mini_pos/service/merchant_serv"
	"github.com/muchlist/mini_pos/service/user_serv"
	"github.com/muchlist/mini_pos/utils/mcrypt"
	"github.com/muchlist/mini_pos/utils/mjwt"
)

func prepareEndPoint(app *fiber.App) {

	// Utils
	cryptoUtils := mcrypt.NewCrypto()
	jwt := mjwt.NewJwt()

	// Merchant Domain
	merchantDao := merchant_dao.New(db.DB)
	merchantService := merchant_serv.NewMerchantService(merchantDao, cryptoUtils)
	merchantHandler := handler.NewMerchantHandler(merchantService)

	// User Domain
	userDao := user_dao.New(db.DB)
	userService := user_serv.NewUserService(userDao, cryptoUtils, jwt)
	userHandler := handler.NewUserHandler(userService)

	app.Use(logger.New())

	// url mapping
	api := app.Group("/api/v1")

	// Merchant Endpoint
	api.Post("/merchant", merchantHandler.CreateMerchant)
	api.Get("/merchant/:id", merchantHandler.GetMerchant)
	api.Get("/merchant", merchantHandler.FindMerchant)
	api.Put("/merchant/:id", merchantHandler.EditMerchant)
	api.Delete("/merchant/:id", merchantHandler.DeleteMerchant)

	// USER Endpont
	api.Get("/users/:id", userHandler.Get)
	api.Get("/users", userHandler.Find)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/profile", middleware.NormalAuth(), userHandler.GetProfile)
	api.Post("/register", middleware.FreshAuth(roles.RoleOwner), userHandler.Register)
	api.Put("/users/:id", middleware.NormalAuth(roles.RoleOwner), userHandler.Edit)
	api.Delete("/users/:id", middleware.NormalAuth(roles.RoleOwner), userHandler.Delete)
}
