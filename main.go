package main

import "github.com/gin-gonic/gin"

// @title Email Service Golang
// @version 1.0
// @description This is an email service API built with Golang and Gin.
func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.Run()
}
