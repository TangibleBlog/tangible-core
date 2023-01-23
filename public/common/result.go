package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func DataResult(c *gin.Context, code int, result interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": code, "data": result})
}

func CodeResult(c *gin.Context, code int) {
	c.JSON(code, gin.H{"code": code})
}
