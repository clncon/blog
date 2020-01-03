package controllers

import (
	"blog/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func ToAddPageHTML(c *gin.Context){
	tags,_:=models.ListTagForIsUsing()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK,"addArticle.html",gin.H{
		"tags":tags,
		"user":user,
	})


}
func Index(c *gin.Context){
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK,"admin/index.html",gin.H{
		"user":user,
	})
}
func ListPage(c *gin.Context){
   pages,_:=models.ListPage()
   c.JSON(200,pages)
}
func DeletePage(c *gin.Context){
	res:= gin.H{}
	res["message"]="success"
	ids:=c.PostForm("ids")
	newids:=strings.Split(ids,",")
	for i:=0;i< len(newids);i++{
		err:=models.DeletePage(newids[i])
		if err!=nil {
			res["message"]=err.Error()
			break;
		}
	}

}
func UpdatePageGet(c *gin.Context){
     id:=c.DefaultQuery("id","13");
     page,err:=models.GetPage(id)
	 tags,_:=models.ListTagForIsUsing()
     tagPages,_:=models.GetTagPage(id)
	 user, _ := c.Get(CONTEXT_USER_KEY)
     if err != nil {
     	c.Error(err)
	 }
	 logrus.Info(tags)
	 c.HTML(http.StatusOK,"editArticle.html",gin.H{
	 	"page":page,
	 	"tags":tags,
	 	"tagPages":tagPages,
	 	"user":user,
	 })
}
func UpdatePagePost(c *gin.Context){
	id:=c.PostForm("id")
	title:= c.PostForm("title")
	desc:= c.PostForm("desc")
	html:= c.PostForm("html")
	source:= c.PostForm("source")
	publish:=c.PostForm("publish")
	tags:=c.PostForm("tags")
	flag,err:=strconv.ParseBool(publish)
	if err != nil {
		msg["message"] = err.Error()
		c.JSON(500, msg)
	}else{
	    models.RemoveTagPageByPageId(id)
		page:=models.Page{Title:title,Desc:desc,Body:html,Source:source,IsPublished:flag}
		err:=models.UpdatePage(id,page)
		arrs:=strings.Split(tags,",")
		for i:=0;i< len(arrs);i++{
			tagId, _ :=strconv.ParseUint(arrs[i],10,64)
			id, _ :=strconv.ParseUint(id,10,64)
			tagPage:=models.TagPage{TagId:uint(tagId),PageId:uint(id)}
			tagPage.Insert()
		}
		if err!=nil {
			msg["message"] = err.Error()
			c.JSON(500, msg)
		}
	}
}
func AddPage(c *gin.Context){
	msg = make(map[string]string)
	var user *models.User
	u, _:= c.Get(CONTEXT_USER_KEY)
	user = u.(*models.User)
	title:= c.PostForm("title")
	desc:= c.PostForm("desc")
	html:= c.PostForm("html")
	source:= c.PostForm("source")
    publish:=c.PostForm("publish")
	tags:=c.PostForm("tags")
	flag,err:=strconv.ParseBool(publish)
	if err != nil {
		msg["message"] = "addTag 接口出错，请查看日志"
		c.JSON(500, msg)
	}else{
		page:=&models.Page{BaseModel:models.BaseModel{Creator:user.Email},Title:title,Desc:desc,Body:html,Source:source,IsPublished:flag}
		id,_:=page.Insert()
		arrs:=strings.Split(tags,",")
		for i:=0;i< len(arrs);i++{
			tagId, _ :=strconv.ParseUint(arrs[i],10,64)
            tagPage:=models.TagPage{TagId:uint(tagId),PageId:id}
            tagPage.Insert()
		}
	}


}
