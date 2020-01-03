package controllers

import (
   "blog/models"
   "github.com/gin-gonic/gin"
   "github.com/sirupsen/logrus"
)
type Post struct {
   Page *models.Page
   Tags []*models.Tag
}

func Page(c *gin.Context){
  id:=c.DefaultQuery("id","1")
  page,err:=models.GetPage(id)
  var tags []*models.Tag
  if err == nil {
     rows,_:=models.DB.Raw("select * from tags t left join tag_pages tp on t.id=tp.tag_id where tp.page_id=?",page.ID).Rows()
     for rows.Next() {
        var tag models.Tag
        models.DB.ScanRows(rows,&tag)
        tags=append(tags,&tag)
     }
  }
  c.HTML(200,"blog/page.html",gin.H{
     "page":page,
     "tags":tags,
  })

}
func Blog(c *gin.Context){
   var posts []*Post
   pages,_:=models.ListPage()
   for i:=0;i<len(pages);i++{
      var post Post
      page:=pages[i]
      post.Page=page
      rows,_:=models.DB.Raw("select * from tags t left join tag_pages tp on t.id=tp.tag_id where tp.page_id=?",page.ID).Rows()
      var tags []*models.Tag
      for rows.Next() {
         var tag models.Tag
         models.DB.ScanRows(rows,&tag)
         tags=append(tags,&tag)
      }
      logrus.Info(len(tags))
      post.Tags=tags
      posts=append(posts,&post)
   }
   logrus.Info(len(posts))
   c.HTML(200,"blog/index.html",gin.H{
      "posts":posts,
   })
}
