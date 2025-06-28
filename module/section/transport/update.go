package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/section/biz"
	"hub-service/module/section/model"
	"hub-service/module/section/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateSection godoc
// @Summary Update a section
// @Description Update an existing section by ID. Only admin and super_admin can access this endpoint.
// @Tags sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Section ID (MongoDB ObjectID)"
// @Param section body model.SectionUpdate true "Section data to update"
// @Success 200 {object} common.Response{data=boolean} "Success"
// @Failure 400 {object} common.AppError "Bad request or invalid ID format"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 403 {object} common.AppError "Forbidden - Only admin and super_admin can access"
// @Failure 404 {object} common.AppError "Section not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/sections/{id} [patch]
func UpdateSection(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data model.SectionUpdate
		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewUpdateSectionBiz(store)

		if err := business.UpdateSection(c.Request.Context(), id, &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
