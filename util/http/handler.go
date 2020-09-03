package http

import "github.com/gin-gonic/gin"

func NewHttpServ() *gin.Engine {
	return gin.Default()
}
