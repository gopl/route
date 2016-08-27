package main

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" xml:"email" form:"email" binding:"required"`
}

type Response struct {
	ReturnCode string `json:"returnCode"`
	Result     string `json:"result"`
}

func pathParameters(c *gin.Context) {
	user := c.Param("user") // 获取路径参数
	params := c.Params      // 路径参数key value对切片
	c.String(http.StatusOK, fmt.Sprintf("user=%v, params=%v", user, params))
}

func queryParameters(c *gin.Context) {
	user := c.Query("user") // 支持Query参数
	c.String(http.StatusOK, "user="+user)
}

func handlingRequest(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusExpectationFailed, Response{"Failed", err.Error()})
	}
	fmt.Println(u)
	c.JSON(http.StatusOK, Response{"OK", "Handle Request"})
}

func main() {
	router := gin.Default()

	// 分组路由
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"user": "password",
	}))
	v1 := authorized.Group("/parameters")
	{
		v1.GET("/path/:user", pathParameters)
		v1.GET("/query", queryParameters)
	}
	router.PUT("/handling/request", handlingRequest)

	router.RunTLS(":12345", "config/server-cert.pem", "config/server-key.pem")
}
