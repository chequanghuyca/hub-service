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
