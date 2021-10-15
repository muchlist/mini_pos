package handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/service/user_serv"
	"github.com/muchlist/mini_pos/utils/mjwt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"github.com/muchlist/mini_pos/utils/sfunc"
	"github.com/muchlist/mini_pos/wrap"
)

func NewUserHandler(userService user_serv.UserServiceAssumer) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

type UserHandler struct {
	service user_serv.UserServiceAssumer
}

// Login login
// @Summary login
// @Description login menggunakan userID dan password untuk mendapatkan JWT Token
// @ID user-login
// @Accept json
// @Produce json
// @Tags Access
// @Param ReqBody body dto.UserLoginRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.UserLoginResponse}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /login [post]
func (u *UserHandler) Login(c *fiber.Ctx) error {
	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if login.Email == "" || login.Password == "" {
		apiErr := rest_err.NewBadRequestError("email atau password tidak boleh kosong")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	response, apiErr := u.service.Login(c.Context(), login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  response,
			Error: nil,
		})
}

// Register menambahkan user
// @Summary register user
// @Description menambahkan user pada merchant sesuai usr owner, endpoint ini membutuhkan hak akses owner, sedangkan akun owner dapat didapatkan ketika membuat Merchant
// @ID user-register
// @Accept json
// @Produce json
// @Tags Access
// @Security bearerAuth
// @Param ReqBody body dto.UserRegisterRequest true "Body raw JSON"
// @Success 200 {object} wrap.RespMsgExample
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /users [post]
func (u *UserHandler) Register(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	var user dto.UserRegisterRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	responseMsg, apiErr := u.service.InsertUser(c.Context(), *claims, dto.UserModel{
		Email:    dto.LowercaseString(user.Email),
		Name:     dto.UppercaseString(user.Name),
		Password: user.Password,
		Role:     dto.LowercaseString(user.Role),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  responseMsg,
			Error: nil,
		})
}

// Edit
// @Summary edit user
// @Description melakukan perubahan data pada user
// @ID user-edit
// @Accept json
// @Produce json
// @Tags Access
// @Security bearerAuth
// @Param id path int true "User ID"
// @Param ReqBody body dto.UserEditRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.UserModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /users/{id} [put]
func (u *UserHandler) Edit(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	var req dto.UserEditRequest

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

	req.ID = userID

	userEdited, apiErr := u.service.EditUser(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  userEdited,
			Error: nil,
		})
}

// RefreshToken
// @Summary refresh token
// @Description mendapatkan token dengan tambahan waktu expired menggunakan refresh token
// @ID user-refresh
// @Accept json
// @Produce json
// @Tags Access
// @Param ReqBody body dto.UserRefreshTokenRequest true "Body raw JSON"
// @Success 200 {object} wrap.Resp{data=dto.UserRefreshTokenResponse}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /refresh [post]
func (u *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.UserRefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	response, apiErr := u.service.Refresh(c.Context(), req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// Delete menghapus user
// @Summary delete user by ID
// @Description menghapus user berdasarkan userID
// @ID user-delete
// @Accept json
// @Produce json
// @Tags Access
// @Security bearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} wrap.RespMsgExample
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /users/{id} [delete]
func (u *UserHandler) Delete(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if claims.Identity == userID {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	apiErr := u.service.DeleteUser(c.Context(), *claims, userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(
		wrap.Resp{
			Data:  fmt.Sprintf("user %d berhasil dihapus", userID),
			Error: nil,
		})
}

// Get menampilkan user berdasarkan id
// @Summary get user by ID
// @Description menampilkan user berdasarkan userID
// @ID user-get
// @Accept json
// @Produce json
// @Tags Access
// @Security bearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} wrap.Resp{data=dto.UserModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /users/{id} [get]
func (u *UserHandler) Get(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	user, apiErr := u.service.GetUserByID(c.Context(), userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  user,
		Error: nil,
	})
}

// Find menampilkan list user
// @Summary find user
// @Description menampilkan daftar user
// @ID user-find
// @Accept json
// @Produce json
// @Tags Access
// @Security bearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset cursor untuk skip data sebanyak offsite"
// @Param search query string false "Search apabila di isi akan melakukan pencarian berdasarkan nama"
// @Success 200 {object} wrap.Resp{data=[]dto.UserModel}
// @Failure 400 {object} wrap.Resp{error=wrap.ErrorExample400}
// @Failure 500 {object} wrap.Resp{error=wrap.ErrorExample500}
// @Router /users [get]
func (u *UserHandler) Find(c *fiber.Ctx) error {
	limit := sfunc.StrToInt(c.Query("limit"), 10)
	offset := sfunc.StrToInt(c.Query("offset"), 0)
	search := c.Query("search")

	userList, apiErr := u.service.FindUsers(c.Context(), search, limit, offset)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	if userList == nil {
		userList = []dto.UserModel{}
	}
	return c.JSON(wrap.Resp{
		Data:  userList,
		Error: nil,
	})
}

// GetProfile mengembalikan user yang sedang login
// @Summary get current profile
// @Description menampilkan profile berdasarkan user yang login saat ini
// @ID user-profile
// @Accept json
// @Produce json
// @Tags Access
// @Security bearerAuth
// @Success 200 {object} wrap.Resp{data=dto.UserModel}
// @Router /profile [get]
func (u *UserHandler) GetProfile(c *fiber.Ctx) error {
	claims, ok := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	if !ok {
		apiErr := rest_err.NewInternalServerError("internal error", errors.New("claims assert failed"))
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	user, apiErr := u.service.GetUserByID(c.Context(), claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(wrap.Resp{
			Data:  nil,
			Error: apiErr,
		})
	}

	return c.JSON(wrap.Resp{
		Data:  user,
		Error: nil,
	})
}
