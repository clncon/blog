package controllers

import (
	"blog/models"
	"blog/system"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Post struct {
	Page *models.Page
	Tags []*models.Tag
}

func Page(c *gin.Context) {
	id := c.DefaultQuery("id", "1")
	tagCounts, _ := models.ListTagCount()
	page, err := models.GetPage(id)
	var tags []*models.Tag
	if err == nil {
		rows, _ := models.DB.Raw("select * from tags t left join tag_pages tp on t.id=tp.tag_id where tp.page_id=?", page.ID).Rows()
		for rows.Next() {
			var tag models.Tag
			models.DB.ScanRows(rows, &tag)
			tags = append(tags, &tag)
		}
	}
	c.HTML(200, "blog/page.html", gin.H{
		"page":      page,
		"tags":      tags,
		"tagCounts": tagCounts,
	})

}
func Blog(c *gin.Context) {
	var posts []*Post
	tagCounts, _ := models.ListTagCount()
	pageSize := system.GetConfiguration().PageSize
	total := models.Total()
	remainder := total % pageSize
	totalPage := total / pageSize
	if remainder != 0 || totalPage == 0 {
		totalPage++
	}
	val := c.DefaultQuery("pageNum", "1")
	var currentPage int
	temp, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		c.JSON(500, "no format number!!!")
		return
	}
	currentPage = int(temp)

	//logrus.Info(currentPage)
	if currentPage < 1 {
		currentPage = 1
	}
	if currentPage > totalPage {
		currentPage = totalPage
	}
	pages, _ := models.ListPage(currentPage, pageSize)
	for i := 0; i < len(pages); i++ {
		var post Post
		page := pages[i]
		post.Page = page
		rows, _ := models.DB.Raw("select * from tags t left join tag_pages tp on t.id=tp.tag_id where tp.page_id=?", page.ID).Rows()
		var tags []*models.Tag
		for rows.Next() {
			var tag models.Tag
			models.DB.ScanRows(rows, &tag)
			tags = append(tags, &tag)
		}
		logrus.Info(len(tags))
		post.Tags = tags
		posts = append(posts, &post)
	}
	c.HTML(200, "blog/index.html", gin.H{
		"posts":       posts,
		"totalPage":   totalPage,
		"currentPage": currentPage,
		"tagCounts":   tagCounts,
	})
}
