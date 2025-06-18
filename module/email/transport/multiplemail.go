package ginemail

import (
	"hub-service/common"
	"hub-service/component/appctx"
	"net/http"

	"github.com/gin-gonic/gin"

	emailmodel "hub-service/module/email/model"
	storagemail "hub-service/module/email/storage"
)

// @Summary Send multiple emails
// @Description Send emails to multiple recipients
// @Tags email
// @Accept json
// @Produce json
// @Param emails body emailmodel.MultipleEmailRequest true "List of email data"
// @Success 200 {string} string "Email sent successfully"
// @Router /api/email/multiple [post]
func MultipleMail(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req emailmodel.MultipleEmailRequest

		if err := c.ShouldBind(&req); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		err := storagemail.MultipleSendEmail(req.AddressesTo, req.Subject, req.Body)

		if err != nil {
			c.JSON(http.StatusInternalServerError, emailmodel.ErrSendEmail(err))
		}

		c.JSON(http.StatusOK, emailmodel.SendEmaiSuccess)
	}
}
