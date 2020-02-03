package controllers

import (
	"blog/helpers"
	"blog/models"
	"blog/system"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alimoeeny/gooauth2"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type GithubUserInfo struct {
	AvatarURL         string      `json:"avatar_url"`
	Bio               interface{} `json:"bio"`
	Blog              string      `json:"blog"`
	Company           interface{} `json:"company"`
	CreatedAt         string      `json:"created_at"`
	Email             interface{} `json:"email"`
	EventsURL         string      `json:"events_url"`
	Followers         int         `json:"followers"`
	FollowersURL      string      `json:"followers_url"`
	Following         int         `json:"following"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	GravatarID        string      `json:"gravatar_id"`
	Hireable          interface{} `json:"hireable"`
	HTMLURL           string      `json:"html_url"`
	ID                int         `json:"id"`
	Location          interface{} `json:"location"`
	Login             string      `json:"login"`
	Name              interface{} `json:"name"`
	OrganizationsURL  string      `json:"organizations_url"`
	PublicGists       int         `json:"public_gists"`
	PublicRepos       int         `json:"public_repos"`
	ReceivedEventsURL string      `json:"received_events_url"`
	ReposURL          string      `json:"repos_url"`
	SiteAdmin         bool        `json:"site_admin"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	Type              string      `json:"type"`
	UpdatedAt         string      `json:"updated_at"`
	URL               string      `json:"url"`
}

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
	c.Redirect(http.StatusSeeOther, "/")
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
		IsAdmin:   false,
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

func Oauth2Callback(c *gin.Context) {
	var (
		userInfo *GithubUserInfo
		user     *models.User
	)
	code := c.Query("code")
	state := c.Query("state")

	//validate state
	session := sessions.Default(c)
	if len(state) == 0 || state != session.Get(SESSION_GITHUB_STATE) {
		c.Abort()
		return
	}
	//remove state from session
	session.Delete(SESSION_GITHUB_STATE)
	session.Save()

	//exchage accesstoken by code
	token, err := exchangeTokenByCode(code)
	if err != nil {
		logrus.Error(err)
		c.Redirect(http.StatusMovedPermanently, "/sign")
		return
	}
	//get github userinfo by accesstoken
	userInfo, err = getGithubUserInfoByAccessToken(token)
	if err != nil {
		logrus.Error(err)
		c.Redirect(http.StatusMovedPermanently, "/signin")
		return
	}

	sessionUser, exists := c.Get(CONTEXT_USER_KEY)
	if exists {
		//已经登陆了
		user, _ = sessionUser.(*models.User)
		_, err1 := models.IsGithubIdExists(userInfo.Login, user.ID)
		if err1 != nil {
			if user.IsAdmin {
				user.GithubLoginId = userInfo.Login
			}
			user.AvatarUrl = userInfo.AvatarURL
			user.GithubUrl = userInfo.HTMLURL
			err = user.UpdateGithubUserInfo()
		} else {
			err = errors.New("this github loginId has bound another account.")
		}
	} else {
		user, err = models.IsExists(userInfo.Email.(string))
		logrus.Info(user.ID)
		if err == nil {
			if user.IsAdmin {
				user.GithubLoginId = userInfo.Login
			}
			user.AvatarUrl = userInfo.AvatarURL
			user.GithubUrl = userInfo.HTMLURL
			user.NickName = userInfo.Name.(string)
			user.UpdateGithubUserInfo()
		} else {
			user = &models.User{
				GithubLoginId: userInfo.Login,
				AvatarUrl:     userInfo.AvatarURL,
				GithubUrl:     userInfo.HTMLURL,
				Email:         userInfo.Email.(string),
				NickName:      userInfo.Name.(string),
				IsAdmin:       false,
			}
			user, err = user.FirstOrCreate()
		}
		if err == nil {
			if user.LockState {
				err = errors.New("Your account have been locked.")
				HandleMessage(c, "Your account have been locked.")
				return
			}
		}
	}

	if err == nil {
		s := sessions.Default(c)
		s.Clear()
		s.Set(SESSION_KEY, user.ID)
		s.Save()
		if user.IsAdmin {
			c.Redirect(http.StatusMovedPermanently, "/admin/index")
		} else {
			c.Redirect(http.StatusMovedPermanently, "/")
		}
		return
	}
}

func exchangeTokenByCode(code string) (accessToken string, err error) {
	var (
		transport *oauth.Transport
		token     *oauth.Token
	)

	transport = &oauth.Transport{Config: &oauth.Config{
		ClientId:     system.GetConfiguration().GithubClientId,
		ClientSecret: system.GetConfiguration().GithubClientSecret,
		Scope:        system.GetConfiguration().GithubScope,
		TokenURL:     system.GetConfiguration().GithubTokenUrl,
		RedirectURL:  system.GetConfiguration().GithubRedirectURL,
	}}
	token, err = transport.Exchange(code)
	if err != nil {
		return
	}
	accessToken = token.AccessToken
	//cache token

	tokenCache := oauth.CacheFile("./request.token")
	if err := tokenCache.PutToken(token); err != nil {
		logrus.Error(err)
	}
	return
}

func getGithubUserInfoByAccessToken(token string) (*GithubUserInfo, error) {
	var (
		resp *http.Response
		body []byte
		err  error
	)

	resp, err = http.Get(fmt.Sprintf("https://api.github.com/user?access_token=%s", token))
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var userInfo GithubUserInfo
	err = json.Unmarshal(body, &userInfo)
	return &userInfo, err
}
