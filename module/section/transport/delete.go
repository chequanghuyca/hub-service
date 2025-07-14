package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/section/biz"
	"hub-service/module/section/storage"
	"hub-service/module/upload/service"
	"hub-service/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteSection godoc
// @Summary Delete a section
// @Description Delete a section and all its related challenges. Only admin and super_admin can access this endpoint.
// @Tags sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Section ID" example("62b4c3789196e8a159933552")
// @Success 200 {object} common.Response{data=string} "Successfully deleted"
// @Failure 400 {object} common.AppError "Bad request - Invalid section ID"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 403 {object} common.AppError "Forbidden - Only admin and super_admin can access"
// @Failure 404 {object} common.AppError "Section not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/sections/{id} [delete]
func DeleteSection(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())

		section, err := store.GetSection(c.Request.Context(), id)
		if err != nil {
			panic(err)
		}

		if section.Image != "" {
			oldFileName := helper.ExtractFileNameFromURL(section.Image)
			if oldFileName != "" {
				r2Service, err := service.NewR2Service()
				if err == nil {
					_ = r2Service.DeleteFile(oldFileName)
				}
			}
		}

		business := biz.NewDeleteSectionBiz(store)
		if err := business.DeleteSection(c.Request.Context(), id); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
