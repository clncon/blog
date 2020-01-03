package controllers

import (
	"blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)
var msg map[string]string
func TagHtml(c *gin.Context){
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK,"tag.html",gin.H{
		"user":user,
	})
}
func ListTag(c *gin.Context){
    tags:=models.MustListTag()
    c.JSON(200,tags)
}

func AddTag(c *gin.Context){
	msg = make(map[string]string)
	var user *models.User
	u, _:= c.Get(CONTEXT_USER_KEY)
	user = u.(*models.User)
	tagName:= c.PostForm("tagName")
	isUsing:= c.PostForm("isUsing")
	flag,err:=strconv.ParseBool(isUsing)
	if err != nil {
		msg["message"]="addTag 接口出错，请查看日志"
		c.JSON(500,msg)
	}else{
		tag:=&models.Tag{BaseModel:models.BaseModel{Creator:user.Email},Name: tagName, IsUsing: flag}
		tag.Insert()
		msg["message"]="添加tag成功"
		c.JSON(200,msg)
	}


}