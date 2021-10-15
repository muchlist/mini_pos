package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/service/merchant_serv"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"github.com/muchlist/mini_pos/utils/sfunc"
	"github.com/muchlist/mini_pos/wrap"
)

func NewMerchantHandler(merchantService merchant_serv.MerchantServiceAssumer) *MerchantHandler {
	return &MerchantHandler{
		service: merchantService,
	}
}

type MerchantHandler struct {
	service merchant_serv.MerchantServiceAssumer
}

// CreateMerchant menambahkan merchant
// @Summary create merchant and owner user
// @Description Menambahkan merchant juga akan membuat user dengan role owner
// @ID merchant-create
// @Accept json
// @Produce json
// @Tags Merchant
// @Security bearerAuth
// @Param ReqBody body dto.MerchantCreateReq true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.MerchantCreateRes}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /merchant [post]
func (m *MerchantHandler) CreateMerchant(c *fiber.Ctx) error {
	var req dto.MerchantCreateReq
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}
	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	response, apiErr := m.service.CreateMerchant(c.Context(), req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(wrap.Resp{
		Data:  response,
		Error: nil,
	})
}

// EditMerchant
// @Summary edit merchant
// @Description melakukan perubahan data pada merchant
// @ID merchant-edit
// @Accept json
// @Produce json
// @Tags Merchant
// @Security bearerAuth
// @Param id path int true "Merchant ID"
// @Param ReqBody body dto.MerchantEditReq true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.Merchant}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /merchant/{id} [put]
func (m *MerchantHandler) EditMerchant(c *fiber.Ctx) error {
	var req dto.MerchantEditReq
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}
	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	response, apiErr := m.service.Edit(c.Context(), req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  response,
		Error: nil,
	})
}

// DeleteMerchant menghapus merchant
// @Summary delete merchant by ID
// @Description menghapus merchant berdasarkan ID
// @ID merchant-delete
// @Accept json
// @Produce json
// @Tags Merchant
// @Security bearerAuth
// @Param id path int true "Merchant ID"
// @Success 200 {object} wrap.RespMsgExample
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /merchant/{id} [delete]
func (m *MerchantHandler) DeleteMerchant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("parameter id diperlukan")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	apiErr := m.service.Delete(c.Context(), id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  fmt.Sprintf("Merchant dengan id %d berhasil dihapus", id),
			Error: nil,
		})
}

// GetMerchant menampilkan merchant berdasarkan id
// @Summary get user by ID
// @Description menampilkan user berdasarkan userID
// @ID merchant-get
// @Accept json
// @Produce json
// @Tags Merchant
// @Security bearerAuth
// @Param id path int true "Merchant ID"
// @Success 200 {object} wrap.Resp{data=dto.Merchant}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /merchant/{id} [get]
func (m *MerchantHandler) GetMerchant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("parameter id diperlukan")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	response, apiErr := m.service.Get(c.Context(), id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  response,
		Error: nil,
	})
}

// FindMerchant menampilkan list merchant
// @Summary find merchant
// @Description menampilkan daftar merchant
// @ID merchant-find
// @Accept json
// @Produce json
// @Tags Merchant
// @Security bearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset cursor untuk skip data sebanyak offsite"
// @Param search query string false "Search apabila di isi akan melakukan pencarian berdasarkan nama merchant"
// @Success 200 {object} wrap.Resp{data=[]dto.Merchant}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /merchant [get]
func (m *MerchantHandler) FindMerchant(c *fiber.Ctx) error {
	search := c.Query("search")
	limit := sfunc.StrToInt(c.Query("limit"), 100)
	cursor := sfunc.StrToInt(c.Query("cursor"), 0)

	response, apiErr := m.service.FindMerchant(c.Context(), search, limit, cursor)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  response,
		Error: nil,
	})
}
