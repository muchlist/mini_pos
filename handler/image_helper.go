package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/mini_pos/utils/logger"
	"github.com/muchlist/mini_pos/utils/mjwt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	jpgExtension  = ".jpg"
	pngExtension  = ".png"
	jpegExtension = ".jpeg"
)

// saveImage return path to save in db
func saveImage(c *fiber.Ctx, claims mjwt.CustomClaim, folder string, imageName string) (string, rest_err.APIError) {
	file, err := c.FormFile("image")
	if err != nil {
		apiErr := rest_err.NewAPIError("File gagal di upload", http.StatusBadRequest, "bad_request", []interface{}{err.Error()})
		logger.Info(fmt.Sprintf("u: %s | formfile | %s", claims.Name, err.Error()))
		return "", apiErr
	}

	fileName := file.Filename
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	if !(fileExtension == jpgExtension || fileExtension == pngExtension || fileExtension == jpegExtension) {
		apiErr := rest_err.NewBadRequestError("Ektensi file tidak di support")
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, apiErr.Error()))
		return "", apiErr
	}

	if file.Size > 2*1024*1024 { // 1 MB
		apiErr := rest_err.NewBadRequestError("Ukuran file tidak dapat melebihi 2MB")
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, apiErr.Error()))
		return "", apiErr
	}

	// rename image
	// path := filepath.Join("static", "image", folder, imageName + fileExtension)
	// pathInDB := filepath.Join("image", folder, imageName + fileExtension)
	path := fmt.Sprintf("static/image/%s/%s", folder, imageName+fileExtension)
	pathInDB := fmt.Sprintf("image/%s/%s", folder, imageName+fileExtension)

	err = c.SaveFile(file, path)
	if err != nil {
		logger.Error(fmt.Sprintf("%s gagal mengupload file", claims.Name), err)
		apiErr := rest_err.NewInternalServerError("File gagal di upload", err)
		return "", apiErr
	}

	return pathInDB, nil
}
