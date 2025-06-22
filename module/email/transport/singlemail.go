package ginemail

import (
	"net/http"

	"hub-service/core/appctx"
	emailmodel "hub-service/module/email/model"
	storagemail "hub-service/module/email/storage"

	"github.com/gin-gonic/gin"
)

// @Summary Send a single email
// @Description Send an email to a single recipient
// @Tags email
// @Accept json
// @Produce json
// @Param email body emailmodel.EmailRequest true "Email data"
// @Success 200 {string} string "Email sent successfully"
// @Router /api/email/single [post]
func SingleMail(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data emailmodel.EmailRequest

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		err := storagemail.SingleSendEmail(data.To, data.Subject, data.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, emailmodel.ErrSendEmail(err))
			return
		}

		c.JSON(http.StatusOK, emailmodel.SendEmailResponse{Status: "success", Message: "Email sent successfully!"})
	}
}
