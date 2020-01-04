package controllers

import (
	"blog/helpers"
	"blog/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SigninGet(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func SignupGet(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}
func ToUserPage(c *gin.Context) {
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(200, "user.html", gin.H{
		"user": user,
	})
}
func ListUser(c *gin.Context) {
	users, _ := models.ListUser()
	c.JSON(200, users)
}
func LogoutGet(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

func SignupPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)

	defer writeJSON(c, res)
	email := c.PostForm("email")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	user := &models.User{
		Email:     email,
		Telephone: telephone,
		Password:  password,
		IsAdmin:   true,
	}

	if len(user.Email) == 0 || len(user.Password) == 0 {
		res["message"] = "email or password cannot be null"
		return
	}
	user.Password = helpers.Md5(user.Email + user.Password)
	err = user.Insert()
	if err != nil {
		res["message"] = "email already exists"
		return
	}

	res["success"] = true
}

func SigninPost(c *gin.Context) {
	var (
		err  error
		user *models.User
	)

	username := c.PostForm("username")
	password := c.PostForm("password")
	logrus.Info(username + ":" + password)
	if username == "" || password == "" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"message": "invalid username or password",
		})
		return
	}
	user, err = models.GetUserByUsername(username)

	if err != nil || user.Password != helpers.Md5(username+password) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"message": "invalid username or password",
		})
		return
	}
	if user.LockState {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"message": "Your account have been locked",
		})
		return
	}

	s := sessions.Default(c)
	s.Clear()
	s.Set(SESSION_KEY, user.ID)
	s.Save()
	if user.IsAdmin {
		c.Redirect(http.StatusMovedPermanently, "/admin/index")
	} else {
		c.Redirect(http.StatusMovedPermanently, "/")

	}
}
