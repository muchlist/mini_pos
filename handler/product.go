package handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/service/product_serv"
	"github.com/muchlist/mini_pos/utils/mjwt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"github.com/muchlist/mini_pos/utils/sfunc"
	"github.com/muchlist/mini_pos/wrap"
	"time"
)

func NewProductHandler(productService product_serv.ProductServiceAssumer) *ProductHandler {
	return &ProductHandler{
		service: productService,
	}
}

type ProductHandler struct {
	service product_serv.ProductServiceAssumer
}

func (u *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	var product dto.ProductCreateRequest
	if err := c.BodyParser(&product); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if err := product.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	createdID, apiErr := u.service.CreateProduct(c.Context(), *claims, dto.ProductModel{
		MerchantID:      claims.Merchant,
		Code:            dto.UppercaseString(product.Code),
		Name:            dto.UppercaseString(product.Name),
		MasterBuyPrice:  product.MasterBuyPrice,
		MasterSellPrice: product.MasterSellPrice,
		Image:           "",
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  fmt.Sprintf("Product dengan ID %d berhasil dibuat", createdID),
			Error: nil,
		})
}

func (u *ProductHandler) Edit(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	productID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	var req dto.ProductEditRequest

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

	req.ID = productID

	productEdited, apiErr := u.service.EditProduct(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  productEdited,
			Error: nil,
		})
}

func (u *ProductHandler) Delete(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	productID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if claims.Identity == productID {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	apiErr := u.service.DeleteProduct(c.Context(), *claims, productID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  fmt.Sprintf("product %d berhasil dihapus", productID),
			Error: nil,
		})
}

func (u *ProductHandler) SetCustomPrice(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	var req dto.ProductPriceRequest
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

	result, apiErr := u.service.SetCustomPrice(c.Context(), *claims, req)

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  result,
			Error: nil,
		})
}

func (u *ProductHandler) Get(c *fiber.Ctx) error {
	productID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}
	outletID := sfunc.StrToInt(c.Query("outlet"), 0)

	product, apiErr := u.service.Get(c.Context(), productID, outletID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  product,
		Error: nil,
	})
}

func (u *ProductHandler) Find(c *fiber.Ctx) error {
	limit := sfunc.StrToInt(c.Query("limit"), 10)
	offset := sfunc.StrToInt(c.Query("offset"), 0)
	search := c.Query("search")
	outlet := sfunc.StrToInt(c.Query("outlet"), 0)

	productList, apiErr := u.service.FindProducts(c.Context(), search, limit, offset, outlet)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if productList == nil {
		productList = []dto.ProductModel{}
	}
	return c.JSON(wrap.Resp{
		Data:  productList,
		Error: nil,
	})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (u *ProductHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	// cek apakah ID cctv && branch ada
	_, apiErr := u.service.Get(c.Context(), id, 0)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%d%v", id, time.Now().Unix())
	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "products", randomName)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	cctvResult, apiErr := u.service.SetImagePath(c.Context(), id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": cctvResult})
}
