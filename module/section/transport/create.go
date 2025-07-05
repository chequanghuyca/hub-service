package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/section/biz"
	"hub-service/module/section/model"
	"hub-service/module/section/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateSection godoc
// @Summary Create a new section
// @Description Create a new section for a challenge. Only admin and super_admin can access this endpoint.
// @Tags sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param section body model.SectionCreate true "Section data"
// @Success 200 {object} common.Response{data=string} "Successfully created. Returns the ID of the new section."
// @Failure 400 {object} common.AppError "Bad request"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 403 {object} common.AppError "Forbidden - Only admin and super_admin can access"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/sections/create [post]
func CreateSection(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var section model.SectionCreate
		if err := c.ShouldBind(&section); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewCreateSectionBiz(store)

		if err := business.CreateSection(c.Request.Context(), &section); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(section.ID))
	}
}
