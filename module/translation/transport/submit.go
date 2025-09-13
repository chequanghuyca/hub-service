package transport

import (
	"errors"
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/translation/biz"
	translationmodel "hub-service/module/translation/model"
	"hub-service/module/translation/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubmitSentenceTranslation godoc
// @Summary Submit sentence translation
// @Description Submit a user's translation for a specific sentence. All authenticated users can access this endpoint.
// @Tags translations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Translation ID" example("62b4c3789196e8a159933552")
// @Param sentence_index path int true "Sentence index" example(0)
// @Param request body translationmodel.SubmitSentenceTranslationRequest true "Translation request (only user_translation field required)"
// @Success 200 {object} common.Response{data=translationmodel.SubmitSentenceTranslationResponse} "Success"
// @Failure 400 {object} common.AppError "Bad request - Invalid translation ID or sentence index"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Translation or sentence not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/translations/{id}/sentences/{sentence_index}/translate [post]
func SubmitSentenceTranslation(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get translation ID from URL parameter
		translationIDStr := c.Param("id")

		// Get sentence index from URL parameter
		sentenceIndexStr := c.Param("sentence_index")
		sentenceIndex, err := strconv.Atoi(sentenceIndexStr)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		// Get user ID from context
		userID := c.MustGet("user_id").(primitive.ObjectID)

		// Bind request body
		var req translationmodel.SubmitSentenceTranslationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		// Convert translation ID to ObjectID
		translationID, err := primitive.ObjectIDFromHex(translationIDStr)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		database := appCtx.GetDatabase()
		if database == nil {
			panic(common.ErrInvalidRequest(errors.New("database connection is nil")))
		}

		store := storage.NewStorage(database)
		apiKey := appCtx.GetEnv("GEMINI_API_KEY")
		baseURL := appCtx.GetEnv("GEMINI_BASE_URL")
		business := biz.NewSubmitTranslationBiz(store, apiKey, baseURL)

		result, err := business.SubmitSentenceTranslation(c.Request.Context(), translationID, sentenceIndex, req.UserTranslation, userID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
