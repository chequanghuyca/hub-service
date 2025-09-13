package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/translation/biz"
	translationmodel "hub-service/module/translation/model"
	"hub-service/module/translation/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTranslation godoc
// @Summary Get a translation by ID
// @Description Get a translation with all its sentences. All authenticated users can access this endpoint.
// @Tags translations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Translation ID" example("62b4c3789196e8a159933552")
// @Success 200 {object} common.Response{data=translationmodel.TranslationWithSentences} "Success"
// @Failure 400 {object} common.AppError "Bad request - Invalid translation ID"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Translation not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/translations/{id} [get]
func GetTranslation(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get translation ID from URL parameter
		translationIDStr := c.Param("id")

		// Validate and convert translation ID
		translationID, err := primitive.ObjectIDFromHex(translationIDStr)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewGetTranslationBiz(store)

		result, err := business.GetTranslation(c.Request.Context(), translationID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

// GetTranslationWithProgress godoc
// @Summary Get translation with user progress
// @Description Get a translation with user's progress and scores. All authenticated users can access this endpoint.
// @Tags translations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Translation ID" example("62b4c3789196e8a159933552")
// @Success 200 {object} common.Response{data=translationmodel.TranslationWithUserProgress} "Success"
// @Failure 400 {object} common.AppError "Bad request - Invalid translation ID"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Translation not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/translations/{id}/progress [get]
func GetTranslationWithProgress(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get translation ID from URL parameter
		translationIDStr := c.Param("id")

		// Validate and convert translation ID
		translationID, err := primitive.ObjectIDFromHex(translationIDStr)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		// Get user ID from context
		userID := c.MustGet("user_id").(primitive.ObjectID)

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewGetTranslationBiz(store)

		result, err := business.GetTranslationWithUserProgress(c.Request.Context(), translationID, userID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

// GetUserTranslationScores godoc
// @Summary Get user translation scores
// @Description Get all translation summaries for a user. All authenticated users can access this endpoint.
// @Tags translations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID" example("62b4c3789196e8a159933552")
// @Success 200 {object} common.Response{data=translationmodel.GetUserTranslationScoresResponse} "Success"
// @Failure 400 {object} common.AppError "Bad request - Invalid user ID"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/translations/user/{user_id}/scores [get]
func GetUserTranslationScores(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from URL parameter
		userIDStr := c.Param("user_id")

		// Validate and convert user ID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewGetTranslationBiz(store)

		result, err := business.GetUserTranslationScores(c.Request.Context(), userID)
		if err != nil {
			panic(err)
		}

		response := translationmodel.GetUserTranslationScoresResponse{
			Summaries: result,
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(response))
	}
}
