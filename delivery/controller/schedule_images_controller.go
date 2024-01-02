package controller

import (
	"fmt"
	"instructor-led-app/config"
	"instructor-led-app/delivery/middleware"
	"instructor-led-app/shared/common"
	"instructor-led-app/usecase"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ScheduleImageController struct {
	scheduleImageUseCase usecase.ScheduleImageUseCase
	authMiddleware       middleware.AuthMiddleware
	rg                   *gin.RouterGroup
}

func (c *ScheduleImageController) UploadImageHandler(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Println("FormFile")
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.MustGet("userID").(string)
	userIdSep := strings.Split(userId, "-")[0]

	log.Println("UserID :", userId)

	allowedExtension := []string{".jpg", ".jpeg", ".png", ".webp"}
	filename := fmt.Sprintf("assets/images/%s-%s", userIdSep, file.Filename)
	ext := filepath.Ext(filename)

	if !c.isAllowedExtension(ext, allowedExtension...) {
		log.Println("Not Allowed Ext")
		common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid file extension")
		return
	}

	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		log.Println("uploaded file")
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	startTime, endTime := ctx.Query("startTime"), ctx.Query("endTime")

	imageDto, err := c.scheduleImageUseCase.UploadImageActivity(userId, filename, startTime, endTime)
	if err != nil {
		log.Println("upload image activity")
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	common.SendSingleResponse(ctx, imageDto, "Upload Image successfully")
}

func (c *ScheduleImageController) isAllowedExtension(ext string, allowedExtensions ...string) bool {
	for _, allowedExtension := range allowedExtensions {
		if ext == allowedExtension {
			return true
		}
	}

	return false
}

func (c *ScheduleImageController) Route() {
	trainer := c.rg.Group(config.TrainerGroup)
	trainer.POST(config.UploadActivityProof, c.authMiddleware.RequireToken("trainer"), middleware.FileSizeLimitMiddleware(10<<20), c.UploadImageHandler)
}

func NewScheduleImageController(scheduleImageUseCase usecase.ScheduleImageUseCase, authMiddleware middleware.AuthMiddleware, rg *gin.RouterGroup) *ScheduleImageController {
	return &ScheduleImageController{scheduleImageUseCase, authMiddleware, rg}
}
