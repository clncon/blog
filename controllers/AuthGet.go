package controllers

import (
	"blog/helpers"
	"blog/system"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthGet(c *gin.Context){
	authType:=c.Param("authType")

	session:=sessions.Default(c)
	uuid:=helpers.UUID()
	session.Delete(SESSION_GITHUB_STATE)
	session.Set(SESSION_GITHUB_STATE,uuid)
	session.Save()

	authurl:="/signin"

	switch authType {
	case "github" :
		authurl = fmt.Sprintf(system.GetConfiguration().GithubAuthUrl,system.GetConfiguration().GithubClientId,uuid)
	case "weibo":
	case "qq":
	default:
	}
	c.Redirect(http.StatusFound,authurl)

}