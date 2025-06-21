package ginemail

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hub-service/component/appctx"
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
		var data emailmodel.MultipleEmailRequest

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		err := storagemail.MultipleSendEmail(data.AddressesTo, data.Subject, data.Body)

		if err != nil {
			c.JSON(http.StatusInternalServerError, emailmodel.ErrSendEmail(err))
		}

		c.JSON(http.StatusOK, emailmodel.SendEmailResponse{Status: "success", Message: "Emails sent successfully!"})
	}
}
