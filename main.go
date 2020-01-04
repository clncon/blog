package main

import (
	"blog/controllers"
	"blog/helpers"
	"blog/models"
	"blog/system"
	"flag"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	configFilePath := flag.String("C", "conf/conf.yaml", "config file path")
	if err:=system.LoadConfiguration(*configFilePath);err!=nil{
		logrus.Error("err parsing config log file",err)
		return
	}
	db,err :=models.InitDB()
	if err != nil {
		logrus.Error("数据库连接失败！！！")
		logrus.Error(err.Error())
	}
	defer db.Close()
	gin.SetMode(gin.ReleaseMode)
	router:=gin.Default()
	setTemplate(router)
	setSessions(router)
	router.Use(sharedData())
	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	//don't need to auth
    //login
    router.GET("/login",controllers.SigninGet)
	router.POST("/login",controllers.SigninPost)
	router.GET("/auth/:authType", controllers.AuthGet)
	router.GET("/",controllers.Blog)
	router.GET("/page",controllers.Page)
	router.POST("/go",controllers.Go)//自动部署
	//register
	if system.GetConfiguration().SignupEnabled {
		router.GET("/register", controllers.SignupGet)
		router.POST("/register", controllers.SignupPost)
	}

	authorized:=router.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		//need to auth
		authorized.POST("/upload",controllers.Upload)
		authorized.GET("/logout",controllers.LogoutGet)
		authorized.GET("/TagHtml",controllers.TagHtml)
		authorized.GET("/listTag",controllers.ListTag)
		authorized.POST("/addTag",controllers.AddTag)
		authorized.GET("/toAddPage",controllers.ToAddPageHTML)
		authorized.POST("/addPage",controllers.AddPage)
		authorized.GET("/updatePage",controllers.UpdatePageGet)
		authorized.POST("/updatePage",controllers.UpdatePagePost)
		authorized.POST("/deletePage",controllers.DeletePage)
		authorized.GET("/index",controllers.Index)
		authorized.GET("/listPage",controllers.ListPage)
		authorized.GET("/listUser",controllers.ListUser)
		authorized.GET("/userPage",controllers.ToUserPage)

	}
	//logrus.Info(getCurrentDirectory())
	router.Run(":80")

}

func setTemplate(engine *gin.Engine) {

	funcMap := template.FuncMap{
		"dateFormat": helpers.DateFormat,
		"substring":  helpers.Substring,
		"isOdd":      helpers.IsOdd,
		"isEven":     helpers.IsEven,
		"truncate":   helpers.Truncate,
		"add":        helpers.Add,
		"minus":      helpers.Minus,
		"listtag":    helpers.ListTag,
	}

	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "./views/**/*"))
}

//setSession initializes sessions & csrf middlewares
func setSessions(router *gin.Engine){
	config:=system.GetConfiguration()
	store:=sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly:true,MaxAge:7*86400,Path:"/"})
	router.Use(sessions.Sessions("gin-session",store))

}
func sharedData() gin.HandlerFunc{
   return func(c *gin.Context){
   	  session:=sessions.Default(c)
   	  if uID:=session.Get(controllers.SESSION_KEY);uID!=nil {
   	  	 user,err:=models.GetUser(uID)
   	  	 if err == nil {
   	  	 	c.Set(controllers.CONTEXT_USER_KEY,user)
		 }
	  }
	  if system.GetConfiguration().SignupEnabled {
	  	 c.Set("SignupEnabled",true)
	  }
	  c.Next()
   }
}
//AuthRequired grants access to authenticated users, requires SharedData middleware
func AdminScopeRequired() gin.HandlerFunc {
    return func(c *gin.Context){
       if user,_:=c.Get(controllers.CONTEXT_USER_KEY);user!=nil{
           if _,ok:=user.(*models.User);ok{
           	  c.Next()
           	  return
		   }
	   }
	   logrus.Warnf("User not authorized to visit%s",c.Request.RequestURI)
       c.HTML(http.StatusForbidden,"login.html",gin.H{
       	  "message":"User not authorized",
	   })
        c.Abort()
	}
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Error(err.Error())
	}
    dir="/root/go/src/blog"
	//dir="/Users/kongmin/go/src/blog"
	return strings.Replace(dir, "\\", "/", -1)
}
