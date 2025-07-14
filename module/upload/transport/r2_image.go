package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/upload/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadR2Image godoc
// @Summary Upload image to Cloudflare R2
// @Description Upload image to R2, file sẽ được đổi tên thành UID duy nhất, trả về URL public qua Worker cho FE đọc.
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image file (jpg, jpeg, png, gif, webp, max 10MB)"
// @Success 200 {object} common.Response{data=string} "Public image URL"
// @Failure 400 {object} common.AppError "Bad request"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/upload/r2-image [post]
func UploadR2Image(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		url, err := service.UploadToR2(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInvalidRequest(err))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(url))
	}
}
