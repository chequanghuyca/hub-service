package emailmodel

import (
	"fmt"
	"hub-service/common"
)

type EmailRequest struct {
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

type ListDataEmail []struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type MultipleEmailRequest struct {
	AddressesTo ListDataEmail `json:"listRecipient" binding:"required"`
	Subject     string        `json:"subject" binding:"required"`
	Body        string        `json:"body" binding:"required"`
}

type EmailResponsePortfolio struct {
	Email   string `json:"email" binding:"required"`
	Message string `json:"message" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

func ErrSendEmail(err error) error {
	return common.NewCustomError(
		err,
		"failed to send email",
		"ErrSendEmail",
	)
}

func ErrSendMailToUser(email string, err error) error {
	return common.NewCustomError(
		err,
		fmt.Sprintf("failed to send email: %s", email),
		fmt.Sprintf("ErrSendMailToUser %s", email),
	)
}

var (
	SendEmaiSuccess = common.SimpleSuccessResponse("Email sent successfully")
)

// API Response Models for Email Module
type SendEmailResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type EmailResponsePortfolioResponse struct {
	Status string                 `json:"status"`
	Data   EmailResponsePortfolio `json:"data"`
}
