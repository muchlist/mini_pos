package user_serv

import (
	"context"
	"github.com/muchlist/mini_pos/dao/user_dao"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/mcrypt"
	"github.com/muchlist/mini_pos/utils/mjwt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"net/http"
	"time"
)

const (
	expiredJWTToken        = 60 * 1       // 1 Hour
	expiredJWTRefreshToken = 15 * 24 * 10 // 15 days
)

type UserServiceAssumer interface {
	UserServiceAccess
	UserServiceModifier
	UserServiceReader
}

type UserServiceReader interface {
	GetUserByID(ctx context.Context, userID int) (*dto.UserModel, rest_err.APIError)
	FindUsers(ctx context.Context, search string, limit int, offset int) ([]dto.UserModel, rest_err.APIError)
}

type UserServiceAccess interface {
	Login(ctx context.Context, login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError)
	Refresh(ctx context.Context, payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
}

type UserServiceModifier interface {
	InsertUser(ctx context.Context, claims mjwt.CustomClaim, user dto.UserModel) (string, rest_err.APIError)
	EditUser(ctx context.Context, claims mjwt.CustomClaim, request dto.UserEditRequest) (*dto.UserModel, rest_err.APIError)
	DeleteUser(ctx context.Context, claims mjwt.CustomClaim, userID int) rest_err.APIError
}

func NewUserService(dao user_dao.UserDaoAssumer, crypto mcrypt.BcryptAssumer, jwt mjwt.JWTAssumer) UserServiceAssumer {
	return &userService{
		dao:    dao,
		crypto: crypto,
		jwt:    jwt,
	}
}

type userService struct {
	dao    user_dao.UserDaoAssumer
	crypto mcrypt.BcryptAssumer
	jwt    mjwt.JWTAssumer
}

// Login
func (u *userService) Login(ctx context.Context, login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError) {
	user, err := u.dao.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, rest_err.NewBadRequestError("Email atau password tidak valid")
	}

	println(user.MerchantID)

	if !u.crypto.IsPWAndHashPWMatch(login.Password, user.Password) {
		return nil, rest_err.NewUnauthorizedError("email atau password tidak valid")
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        string(user.Name),
		ExtraMinute: expiredJWTToken, // 1 Hour
		Type:        mjwt.Access,
		Fresh:       true,
		Role:        string(user.Role),
		Merchant:    user.MerchantID,
		Outlet:      user.DefOutlet,
	}

	RefreshClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        string(user.Name),
		ExtraMinute: expiredJWTRefreshToken, // 15 days
		Type:        mjwt.Refresh,
		Fresh:       false,
		Role:        string(user.Role),
		Merchant:    user.MerchantID,
		Outlet:      user.DefOutlet,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserLoginResponse{
		ID:           user.ID,
		Email:        string(user.Email),
		Name:         string(user.Name),
		MerchantID:   user.MerchantID,
		DefOutlet:    user.DefOutlet,
		Role:         string(user.Role),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expired:      time.Now().Add(time.Minute * time.Duration(expiredJWTToken)).Unix(),
	}

	return &userResponse, nil
}

// InsertUser melakukan register user oleh akun owner kepada akun akun dibawah merchant
func (u *userService) InsertUser(ctx context.Context, claims mjwt.CustomClaim, user dto.UserModel) (string, rest_err.APIError) {

	timeNow := time.Now().Unix()
	hashPassword, err := u.crypto.GenerateHash(user.Password)
	if err != nil {
		return "", err
	}

	user.Password = hashPassword // hashing password default yang diberikan
	user.CreatedAt = timeNow
	user.UpdatedAt = timeNow
	user.MerchantID = claims.Merchant // merchant ID adalah sama dengan merchant id owner

	successMsg, err := u.dao.Insert(ctx, user)
	if err != nil {
		return "", err
	}
	return successMsg, nil
}

// EditUser
func (u *userService) EditUser(ctx context.Context, claims mjwt.CustomClaim, request dto.UserEditRequest) (*dto.UserModel, rest_err.APIError) {
	editParams := dto.UserEditModel{
		WhereID:         request.ID,
		WhereMerchantID: claims.Merchant, // <--- user yang diedit harus memiliki merchant id yang sama
		Email:           dto.LowercaseString(request.Email),
		Name:            dto.UppercaseString(request.Name),
		Role:            dto.LowercaseString(request.Role),
		UpdatedAt:       time.Now().Unix(),
		DefOutlet:       request.DefOutlet,
	}

	result, err := u.dao.Edit(ctx, editParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Refresh token
func (u *userService) Refresh(ctx context.Context, payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError) {
	token, apiErr := u.jwt.ValidateToken(payload.RefreshToken)
	if apiErr != nil {
		return nil, apiErr
	}
	claims, apiErr := u.jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}

	// cek apakah tipe claims token yang dikirim adalah tipe refresh
	if claims.Type != mjwt.Refresh {
		return nil, rest_err.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{"not a refresh token"})
	}

	// mendapatkan data terbaru dari user
	user, apiErr := u.dao.GetByID(ctx, claims.Identity)
	if apiErr != nil {
		return nil, apiErr
	}

	accessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        string(user.Name),
		ExtraMinute: expiredJWTToken, // 1 Hour
		Type:        mjwt.Access,
		Fresh:       false,
		Role:        string(user.Role),
		Merchant:    user.MerchantID,
		Outlet:      user.DefOutlet,
	}

	accessToken, err := u.jwt.GenerateToken(accessClaims)
	if err != nil {
		return nil, err
	}

	userRefreshTokenResponse := dto.UserRefreshTokenResponse{
		AccessToken: accessToken,
		Expired:     time.Now().Add(time.Minute * time.Duration(expiredJWTToken)).Unix(),
	}

	return &userRefreshTokenResponse, nil
}

// DeleteUser
func (u *userService) DeleteUser(ctx context.Context, claims mjwt.CustomClaim, userID int) rest_err.APIError {
	err := u.dao.Delete(ctx, userID, claims.Merchant)
	if err != nil {
		return err
	}
	return nil
}

// GetUser mendapatkan user dari database
func (u *userService) GetUserByID(ctx context.Context, userID int) (*dto.UserModel, rest_err.APIError) {
	user, err := u.dao.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUsers
func (u *userService) FindUsers(ctx context.Context, search string, limit int, offset int) ([]dto.UserModel, rest_err.APIError) {
	userList, err := u.dao.FindWithPagination(ctx, user_dao.FindPaginationParams{
		Search: search,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return userList, nil
}
