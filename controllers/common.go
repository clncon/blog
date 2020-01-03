package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SESSION_KEY          = "UserID"       // session key
	CONTEXT_USER_KEY     = "User"         // context user key
	SESSION_GITHUB_STATE = "GITHUB_STATE" // github state session key
	SESSION_CAPTCHA      = "GIN_CAPTCHA"  // captcha session key
)

func Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry,I lost myself!")
}

func HandleMessage(c *gin.Context, message string) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": message,
	})
}

func writeJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}
