package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/section/biz"
	"hub-service/module/section/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListSimpleSection godoc
// @Summary Get all sections with only id and title
// @Description Get all sections with only id and title, optionally filtered by title search. No pagination. All authenticated users can access this endpoint.
// @Tags sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param title query string false "Search by section title (case-insensitive)"
// @Success 200 {object} common.Response{data=[]model.SectionSimple} "Success"
// @Failure 400 {object} common.AppError "Bad request"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/sections/simple [get]
func ListSimpleSection(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get search title from query parameter
		title := c.Query("title")

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewListSimpleSectionBiz(store)

		result, err := business.ListSimpleSection(c.Request.Context(), title)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, nil, nil))
	}
}
