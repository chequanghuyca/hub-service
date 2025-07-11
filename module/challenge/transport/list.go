package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/challenge/biz"
	_ "hub-service/module/challenge/model"
	"hub-service/module/challenge/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListChallenge godoc
// @Summary List challenges
// @Description Get a list of translation challenges with pagination and search. All authenticated users can access this endpoint.
// @Tags challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param section_id query string false "Filter by section ID"
// @Param search query string false "Search in title and content (case-insensitive)"
// @Success 200 {object} common.Response{data=[]model.Challenge,meta=common.Paging} "Success"
// @Failure 400 {object} common.AppError "Bad request"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/challenges/list [get]
func ListChallenge(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var paging common.Paging
		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		paging.Fulfill()

		// Get section_id from query parameter
		sectionID := c.Query("section_id")

		// Get search from query parameter
		search := c.Query("search")

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewListChallengeBiz(store)

		result, err := business.ListChallenge(c.Request.Context(), &paging, sectionID, search)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, nil))
	}
}
