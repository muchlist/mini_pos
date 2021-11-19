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

// CreateOutlet menambahkan outlets
// @Summary create outlet for merchant user
// @Description Menambahkan outlets sesuai dengan ID merchant yang melekat di user
// @ID outlet-create
// @Accept json
// @Produce json
// @Tags Outlet
// @Security bearerAuth
// @Param ReqBody body dto.OutletCreateRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=wrap.RespMsgExample}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /outlets [post]
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

// Edit
// @Summary edit outlet
// @Description melakukan perubahan data pada outlet
// @ID outlet-edit
// @Accept json
// @Produce json
// @Tags Outlet
// @Security bearerAuth
// @Param id path int true "Outlet ID"
// @Param ReqBody body dto.OutletEditRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.OutletModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /outlets/{id} [put]
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

// Delete menghapus outlet
// @Summary delete outlet by ID
// @Description menghapus outlet berdasarkan ID
// @ID outlet-delete
// @Accept json
// @Produce json
// @Tags Outlet
// @Security bearerAuth
// @Param id path int true "Outlet ID"
// @Success 200 {object} wrap.RespMsgExample
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /outlets/{id} [delete]
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

// Get menampilkan outlet berdasarkan id
// @Summary get outlet by ID
// @Description menampilkan outlet berdasarkan userID
// @ID outlet-get
// @Accept json
// @Produce json
// @Tags Outlet
// @Security bearerAuth
// @Param id path int true "Outlet ID"
// @Success 200 {object} wrap.Resp{data=dto.OutletModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /outlets/{id} [get]
func (u *OutletHandler) Get(c *fiber.Ctx) error {
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

	outlet, apiErr := u.service.GetOutletByID(c.Context(), *claims, outletID)
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

// Find menampilkan list outlet
// @Summary find outlet
// @Description menampilkan daftar outlet
// @ID outlet-find
// @Accept json
// @Produce json
// @Tags Access
// @Security bearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset cursor untuk skip data sebanyak offsite"
// @Param search query string false "Search apabila di isi akan melakukan pencarian berdasarkan nama outlet"
// @Success 200 {object} wrap.Resp{data=[]dto.OutletModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /outlets [get]
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

// GetCurrentOutlet menampilkan outlet berdasarkan user yang login
// @Summary get outlet by current user
// @Description menampilkan outlet berdasarkan user yang saat ini login
// @ID outlet-curent
// @Accept json
// @Produce json
// @Tags Outlet
// @Security bearerAuth
// @Success 200 {object} wrap.Resp{data=dto.OutletModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /current-outlet [get]
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

	outlet, apiErr := u.service.GetOutletByID(c.Context(), *claims, claims.Identity)
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
