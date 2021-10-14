package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/muchlist/mini_pos/configs/roles"
	"github.com/muchlist/mini_pos/dao/merchant_dao"
	"github.com/muchlist/mini_pos/dao/outlet_dao"
	"github.com/muchlist/mini_pos/dao/product_dao"
	"github.com/muchlist/mini_pos/dao/user_dao"
	"github.com/muchlist/mini_pos/db"
	"github.com/muchlist/mini_pos/handler"
	"github.com/muchlist/mini_pos/middleware"
	"github.com/muchlist/mini_pos/service/merchant_serv"
	"github.com/muchlist/mini_pos/service/outlet_serv"
	"github.com/muchlist/mini_pos/service/product_serv"
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

	// Outlet Domain
	outletDao := outlet_dao.New(db.DB)
	outletService := outlet_serv.NewOutletService(outletDao)
	outletHandler := handler.NewOutletHandler(outletService)

	// Product Domain
	productDao := product_dao.New(db.DB)
	productService := product_serv.NewProductService(productDao)
	productHandler := handler.NewProductHandler(productService)

	app.Use(logger.New())

	app.Static("/image/products", "./static/image/products")

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

	// Outlet Endpont
	api.Get("/outlets/:id", middleware.NormalAuth(), outletHandler.Get)
	api.Get("/outlets", middleware.NormalAuth(), outletHandler.Find)
	api.Get("/current-outlet", middleware.NormalAuth(), outletHandler.GetCurrentOutlet)
	api.Post("/outlets", middleware.NormalAuth(roles.RoleOwner), outletHandler.CreateOutlet)
	api.Put("/outlets/:id", middleware.NormalAuth(roles.RoleOwner), outletHandler.Edit)
	api.Delete("/outlets/:id", middleware.NormalAuth(roles.RoleOwner), outletHandler.Delete)

	// Product Endpont
	api.Get("/products/:id", middleware.NormalAuth(), productHandler.Get)
	api.Get("/products", middleware.NormalAuth(), productHandler.Find)
	api.Post("/products", middleware.NormalAuth(roles.RoleOwner), productHandler.CreateProduct)
	api.Put("/products/:id", middleware.NormalAuth(roles.RoleOwner), productHandler.Edit)
	api.Delete("/products/:id", middleware.NormalAuth(roles.RoleOwner), productHandler.Delete)
	api.Post("/set-price", middleware.NormalAuth(roles.RoleOwner), productHandler.SetCustomPrice)
	api.Post("/products-image/:id", middleware.NormalAuth(roles.RoleOwner), productHandler.UploadImage)

}
