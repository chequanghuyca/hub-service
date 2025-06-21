package ginemail

import (
	"hub-service/common"
	"hub-service/component/appctx"
	"hub-service/helper"
	emailmodel "hub-service/module/email/model"
	storageemail "hub-service/module/email/storage"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @Summary Response to portfolio inquiry
// @Description Send a response email for portfolio inquiry
// @Tags email
// @Accept json
// @Produce json
// @Param response body emailmodel.EmailResponsePortfolio true "Response data"
// @Success 200 {string} string "Email sent successfully"
// @Router /api/email/response-portfolio [post]
func ResponseEmailPortfolio(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		godotenv.Load()

		var req emailmodel.EmailResponsePortfolio

		if err := c.ShouldBind(&req); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		dataSendContact := helper.MailResponseData{
			Name:    req.Name,
			MyPhone: os.Getenv("SYSTEM_PHONE_NUMBER"),
			MyEmail: os.Getenv("SYSTEM_EMAIL"),
		}

		err := storageemail.SingleSendEmail(
			req.Email,
			helper.GetSubjectMailResponse(),
			helper.GetBodyMailResponse(dataSendContact),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, emailmodel.ErrSendEmail(err))
		}

		errSendMe := storageemail.ResponseMeEmail(req.Message)

		if errSendMe != nil {
			c.JSON(http.StatusInternalServerError, emailmodel.ErrSendEmail(errSendMe))
		}

		c.JSON(http.StatusOK, emailmodel.SendEmaiSuccess)
	}
}

func ResponseMeMail(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data emailmodel.EmailRequest

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		// The purpose seems to be sending a notification email to the system.
		// The message will be the body from the request.
		err := storageemail.ResponseMeEmail(data.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// This response is for the client that triggered the action
		c.JSON(http.StatusOK, emailmodel.SendEmailResponse{Status: "success", Message: "Notification sent successfully!"})
	}
}
