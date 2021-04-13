package main

import (
	"OAuth2/OAuth2_01/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"log"
	"net/http"
	"time"
)

var srv *server.Server

func main() {
	//创建管理对象
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	err := clientStore.Set("clienta", &models.Client{
		ID:     "clienta",
		Secret: "123",
		Domain: "http://localhost:8080",
	})
	if err != nil {
		log.Fatal(err)
	}
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	r := gin.New()
	r.Use(utils.ErrorHandler())
	//响应授权码,通过userAuthorizeHandler
	r.GET("/auth", func(context *gin.Context) {
		err := srv.HandleAuthorizeRequest(context.Writer, context.Request)
		if err != nil {
			log.Println(err)
		}
	})

	//根据授权码获取token
	r.POST("/token", func(context *gin.Context) {
		err := srv.HandleTokenRequest(context.Writer, context.Request)
		if err != nil {
			panic(err.Error())
		}
	})

	//如果没有登录 则跳转登录界面
	r.Any("/login", func(context *gin.Context) {
		data := map[string]string{
			"error": "",
		}
		if context.Request.Method == "POST" {
			uname, upass := context.PostForm("username"), context.PostForm("userPass")
			if uname+upass == "shneyi123" {
				utils.SaveUserSesison(context, uname)
				context.Redirect(302, "/auth?"+context.Request.URL.RawQuery)
				return
			} else {
				data["error"] = "用户名密码错误"
			}
		}
		context.HTML(200, "login.html", data)
	})

	r.POST("/info", func(context *gin.Context) {
		token, err2 := srv.ValidationBearerToken(context.Request)
		if err2 != nil {
			panic(err2.Error())
		}
		ret := gin.H{
			"user_id": token.GetUserID(),
			"expire":  int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		}
		context.JSON(200, ret)
	})

	//加载模版
	r.LoadHTMLGlob("public/*.html")
	r.Run(":80")
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	//得到当前的用户id
	if userID = utils.GetUserSession(r); userID == "" {
		w.Header().Set("Location", "/login?"+r.URL.RawQuery) //进入login页面登录
		w.WriteHeader(302)
	}
	return
}
