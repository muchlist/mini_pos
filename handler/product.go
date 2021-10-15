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

// CreateProduct menambahkan outlets
// @Summary create product for merchant user
// @Description Menambahkan product sesuai dengan ID merchant yang melekat di user
// @ID product-create
// @Accept json
// @Produce json
// @Tags Product
// @Security bearerAuth
// @Param ReqBody body dto.ProductCreateRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=wrap.RespMsgExample}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /products [post]
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

// Edit
// @Summary edit product
// @Description melakukan perubahan data pada product
// @ID product-edit
// @Accept json
// @Produce json
// @Tags Product
// @Security bearerAuth
// @Param id path int true "Product ID"
// @Param ReqBody body dto.ProductEditRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.OutletModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /products/{id} [put]
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

// Delete menghapus product
// @Summary delete product by ID
// @Description menghapus product berdasarkan ID
// @ID product-delete
// @Accept json
// @Produce json
// @Tags Outlet
// @Security bearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} wrap.RespMsgExample
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /products/{id} [delete]
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

// SetCustomPrice
// @Summary menambahkan harga custom
// @Description menambahkan harga custom product pada outlet tertentu
// @ID product-set-price
// @Accept json
// @Produce json
// @Tags Product
// @Security bearerAuth
// @Param ReqBody body dto.ProductPriceRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.OutletModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /set-price/{id} [post]
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

// Get menampilkan product berdasarkan id
// @Summary get product by ID
// @Description menampilkan product berdasarkan userID, query outlet untuk mendapatkan harga custom pada outlet tertentu
// @ID product-get
// @Accept json
// @Produce json
// @Tags Product
// @Security bearerAuth
// @Param id path int true "Product ID"
// @Param outlet query int false "Outlet Price"
// @Success 200 {object} wrap.Resp{data=dto.ProductModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /products/{id} [get]
func (u *ProductHandler) Get(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	outletID := sfunc.StrToInt(c.Query("outlet"), 0)
	productID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	product, apiErr := u.service.Get(c.Context(), *claims, productID, outletID)
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

// Find menampilkan list product
// @Summary find product
// @Description menampilkan daftar product untuk merchant tertentu, gunakan query outlet untuk mendapatkan harga product sesuai outlet
// @ID product-find
// @Accept json
// @Produce json
// @Tags Product
// @Security bearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset cursor untuk skip data sebanyak offsite"
// @Param search query string false "Search apabila di isi akan melakukan pencarian berdasarkan nama outlet"
// @Param outlet query int false "tambahkan outlet untuk melihat harga outlet tertentu"
// @Success 200 {object} wrap.Resp{data=[]dto.OutletModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /products [get]
func (u *ProductHandler) Find(c *fiber.Ctx) error {
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
	outlet := sfunc.StrToInt(c.Query("outlet"), 0)

	productList, apiErr := u.service.FindProducts(c.Context(), *claims, product_serv.FindProductsParams{
		Search:         search,
		Limit:          limit,
		Offset:         offset,
		OutletSpecific: outlet,
	})
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

// UploadImage
// @Summary menambahkan foto pada product
// @Description menambahkan foto pada product
// @ID product-upload-photo
// @Accept json
// @Produce json
// @Tags Product
// @Security bearerAuth
// @Param image formData file true "file gambar"
// @Success 200 {object} wrap.Resp{data=dto.ProductModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /products-image/{id} [post]
func (u *ProductHandler) UploadImage(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	// cek apakah ID cctv && branch ada
	_, apiErr := u.service.Get(c.Context(), *claims, id, 0)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	randomName := fmt.Sprintf("%d%v", id, time.Now().Unix())
	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "products", randomName)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	// update path image di database
	result, apiErr := u.service.SetImagePath(c.Context(), id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  result,
		Error: nil,
	})
}
