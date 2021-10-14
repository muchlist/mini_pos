package handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/service/outlet_serv"
	"github.com/muchlist/mini_pos/utils/mjwt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"github.com/muchlist/mini_pos/utils/sfunc"
	"github.com/muchlist/mini_pos/wrap"
	"time"
)

func NewOutletHandler(outletService outlet_serv.OutletServiceAssumer) *OutletHandler {
	return &OutletHandler{
		service: outletService,
	}
}

type OutletHandler struct {
	service outlet_serv.OutletServiceAssumer
}

func (u *OutletHandler) CreateOutlet(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	var outlet dto.OutletCreateRequest
	if err := c.BodyParser(&outlet); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if err := outlet.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	createdID, apiErr := u.service.CreateOutlet(c.Context(), *claims, dto.OutletModel{
		MerchantID: claims.Merchant,
		OutletName: dto.UppercaseString(outlet.OutletName),
		Address:    outlet.Address,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  fmt.Sprintf("Outlet dengan ID %d berhasil dibuat", createdID),
			Error: nil,
		})
}

func (u *OutletHandler) Edit(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	outletID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	var req dto.OutletEditRequest

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

	req.ID = outletID

	outletEdited, apiErr := u.service.EditOutlet(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  outletEdited,
			Error: nil,
		})
}

func (u *OutletHandler) Delete(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	outletID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if claims.Identity == outletID {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	apiErr := u.service.DeleteOutlet(c.Context(), *claims, outletID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  fmt.Sprintf("outlet %d berhasil dihapus", outletID),
			Error: nil,
		})
}

func (u *OutletHandler) Get(c *fiber.Ctx) error {
	outletID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	outlet, apiErr := u.service.GetOutletByID(c.Context(), outletID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  outlet,
		Error: nil,
	})
}

func (u *OutletHandler) Find(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}
	limit := sfunc.StrToInt(c.Query("limit"), 10)
	offset := sfunc.StrToInt(c.Query("offset"), 0)
	search := c.Query("search")

	outletList, apiErr := u.service.FindOutlets(c.Context(), *claims, search, limit, offset)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if outletList == nil {
		outletList = []dto.OutletModel{}
	}
	return c.JSON(wrap.Resp{
		Data:  outletList,
		Error: nil,
	})
}

func (u *OutletHandler) GetCurrentOutlet(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	outletID := claims.Outlet
	if outletID == 0 {
		apiErr := rest_err.NewNotFoundError("User belum berada di outlet yang spesifik")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	outlet, apiErr := u.service.GetOutletByID(c.Context(), claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  outlet,
		Error: nil,
	})
}
