package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/service/merchant_serv"
	"github.com/muchlist/mini_pos/utils/rest_err"
)

func NewMerchantHandler(merchantService merchant_serv.MerchantServiceAssumer) *MerchantHandler {
	return &MerchantHandler{
		service: merchantService,
	}
}

type MerchantHandler struct {
	service merchant_serv.MerchantServiceAssumer
}

func (m *MerchantHandler) CreateMerchant(c *fiber.Ctx) error {
	var req dto.MerchantCreateReq
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := m.service.CreateMerchant(c.Context(), req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}
