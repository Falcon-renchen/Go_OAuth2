package main

import (
	"OAuth2/OAuth2_01/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	authServerURL = "http://oauth.jtthink.com"
)
var (
	oauth2Config = oauth2.Config{
		ClientID:     "clienta",
		ClientSecret: "123",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/auth",  //获得授权码地址
			TokenURL: authServerURL + "/token", //获取token地址
		},
		RedirectURL: "http://localhost:8080/getcode",
		Scopes:      []string{"all"},
	}
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("public/*")
	//codeUrl,_:=url.ParseRequestURI("http://localhost:8080/getcode")
	//loginUrl:="http://oauth.jtthink.com/auth?" +
	//	"response_type=code&client_id=clienta&redirect_uri="+
	//	codeUrl.String()
	loginUrl := oauth2Config.AuthCodeURL("myclient")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "a-index.html", map[string]string{
			"loginUrl": loginUrl,
		})
	})
	r.GET("/getcode", func(c *gin.Context) {
		code, _ := c.GetQuery("code")
		//c.JSON(200,gin.H{"code":code})
		//请求token
		token, err := oauth2Config.Exchange(c, code)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(200, token)
		}
	})

	r.GET("/info", func(context *gin.Context) {
		token := context.Query("token")
		ret := utils.GetUserInfo(authServerURL+"/info", token, false)
		context.Writer.Header().Add("Content-type", "application/json")
		context.String(200, ret)
	})
	r.Run(":8080")
}
