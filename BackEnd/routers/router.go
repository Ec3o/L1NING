package routers

import (
	"BackEnd/actionmodules"
	"BackEnd/config"
	"BackEnd/controllers"
	"BackEnd/models"
	"fmt"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims models.Claims

func Run() {
	r := gin.Default()
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AddAllowHeaders("Authorization")
	CORSConfig.AllowCredentials = true
	r.Use(cors.New(CORSConfig))
	r.LoadHTMLGlob("templates/*")
	r.GET("/", ShowindexPage)
	r.GET("/api/WordKiller", actionmodules.StartKill)
	r.GET("/sunnysport", actionmodules.GetData)
	auth := r.Group("/api")
	{
		auth.POST("/login", controllers.Login)

		auth.POST("/register", controllers.Register)
	}
	Todo := r.Group("/api")
	Todo.Use(Authorize())
	{
		// 添加待办事项
		Todo.POST("/todos", controllers.AddTodo)

		// 获取所有待办事项
		Todo.GET("/todos", controllers.GetTodo)

		// 获取单个待办事项
		Todo.GET("/todos/:id", controllers.GetTodo)

		// 更新待办事项
		Todo.PUT("/todos/:id", controllers.UpdateTodo)

		// 删除待办事项
		Todo.DELETE("/todos/:id", controllers.DeleteTodo)
	}
	r.Run(config.Runningport)

}
func ShowindexPage(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, "请求头中 auth 为空")
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, "无效的 Token")
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
				c.JSON(http.StatusUnauthorized, "Token 已过期")
				c.Abort()
				return
			}
			fmt.Println(claims)
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, "令牌错误")
			c.Abort()
			return
		}
	}
}
