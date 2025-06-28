package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/section/biz"
	"hub-service/module/section/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListSection godoc
// @Summary Get a list of sections with pagination and user scores
// @Description Get a list of sections with pagination and user scores. All authenticated users can access this endpoint.
// @Tags sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} common.Response{data=[]model.SectionWithScore,meta=common.Paging} "Success"
// @Failure 400 {object} common.AppError "Bad request"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/sections/list [get]
func ListSection(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var paging common.Paging
		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		paging.Fulfill()

		// Get user ID from context
		userID := c.MustGet("user_id").(primitive.ObjectID)

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewListSectionBiz(store)

		result, err := business.ListSection(c.Request.Context(), &paging, userID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, nil))
	}
}
