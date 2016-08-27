package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name"`
	Email string `json:"email" xml:"email" form:"email"`
}

type Response struct {
	ReturnCode string `json:"returnCode"`
	Result     string `json:"result"`
}

func pathParameters(c echo.Context) error {
	user := c.Param("user")
	paramNames := c.ParamNames()   // 路径参数名列表
	paramValues := c.ParamValues() // 路径参数值列表
	return c.String(http.StatusOK, fmt.Sprintf("user=%v, paramNames=%v, paramValues=%v",
		user, paramNames, paramValues))
}

func queryParameters(c echo.Context) error {
	path := c.Path()              // 获取handler的注册路径
	user := c.QueryParam("user")  // 获取Query参数的第一个值
	parameters := c.QueryParams() //获取所有Query参数map
	return c.String(http.StatusOK, fmt.Sprintf("path=%v, user=%v, params=%v", path, user, parameters))
}

func handlingRequest(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusExpectationFailed, Response{"Failed", err.Error()})
	}
	fmt.Println(u)
	return c.JSON(http.StatusOK, Response{"OK", "Handle Request"})
}

// 认证
func auth(username, password string) bool {
	if username == "user" && password == "password" {
		return true
	}
	return false
}

func main() {
	router := echo.New()
	router.Use(middleware.BasicAuth(auth))
	router.Use(middleware.BodyLimit("80B"))
	router.Use(middleware.Logger())

	router.GET("/parameters/path/:user", pathParameters)
	router.GET("/parameters/query", queryParameters)
	router.PUT("/handling/request", handlingRequest)

	router.Run(standard.WithConfig(engine.Config{
		Address:     ":12345",
		TLSCertFile: "config/server-cert.pem",
		TLSKeyFile:  "config/server-key.pem",
	}))
}
