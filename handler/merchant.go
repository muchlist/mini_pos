package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/service/merchant_serv"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"github.com/muchlist/mini_pos/utils/sfunc"
	"strconv"
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

func (m *MerchantHandler) EditMerchant(c *fiber.Ctx) error {
	var req dto.MerchantEditReq
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := m.service.Edit(c.Context(), req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

func (m *MerchantHandler) DeleteMerchant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("parameter id diperlukan")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := m.service.Delete(c.Context(), id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("Merchant dengan id %d berhasil dihapus", id)})
}

func (m *MerchantHandler) GetMerchant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error(), "data": nil})
	}

	response, apiErr := m.service.Get(c.Context(), id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// StrToInt merubah str ke int dengan nilai default
func StrToInt(text string, defaultReturn int) int {
	number := defaultReturn
	if text != "" {
		var err error
		number, err = strconv.Atoi(text)
		if err != nil {
			number = defaultReturn
		}
	}
	return number
}

func (m *MerchantHandler) FindMerchant(c *fiber.Ctx) error {
	search := c.Query("search")
	limit := sfunc.StrToInt(c.Query("limit"), 100)
	cursor := sfunc.StrToInt(c.Query("cursor"), 0)

	response, apiErr := m.service.FindMerchant(c.Context(), search, limit, cursor)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}
