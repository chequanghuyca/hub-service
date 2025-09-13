package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/translation/biz"
	translationmodel "hub-service/module/translation/model"
	"hub-service/module/translation/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTranslation godoc
// @Summary Create a new translation
// @Description Create a new translation with sentences. Only admin and super_admin can access this endpoint.
// @Tags translations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param translation body translationmodel.TranslationCreate true "Translation data"
// @Success 200 {object} common.Response{data=translationmodel.Translation} "Successfully created translation"
// @Failure 400 {object} common.AppError "Bad request"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 403 {object} common.AppError "Forbidden - Only admin and super_admin can access"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/translations/create [post]
func CreateTranslation(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var translation translationmodel.TranslationCreate
		if err := c.ShouldBind(&translation); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		apiKey := appCtx.GetEnv("GEMINI_API_KEY")
		baseURL := appCtx.GetEnv("GEMINI_BASE_URL")
		business := biz.NewCreateTranslationBiz(store, apiKey, baseURL)

		result, err := business.CreateTranslation(c.Request.Context(), &translation)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
