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

// GetSection godoc
// @Summary Get a section by ID
// @Description Get a section with all its related challenges and user score. All authenticated users can access this endpoint.
// @Tags sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Section ID" example("62b4c3789196e8a159933552")
// @Success 200 {object} common.Response{data=model.SectionWithChallenges} "Success"
// @Failure 400 {object} common.AppError "Bad request - Invalid section ID"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Section not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/sections/{id} [get]
func GetSection(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get section ID from URL parameter
		sectionIDStr := c.Param("id")

		// Validate and convert section ID
		sectionID, err := primitive.ObjectIDFromHex(sectionIDStr)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		// Get user ID from context
		userID := c.MustGet("user_id").(primitive.ObjectID)

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewGetSectionBiz(store)

		result, err := business.GetSection(c.Request.Context(), sectionID, userID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
