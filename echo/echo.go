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
	return c.String(http.StatusOK, "user="+user)
}

func queryParameters(c echo.Context) error {
	user := c.QueryParam("user") // 支持Query参数
	parameters := c.QueryParams()
	return c.String(http.StatusOK, fmt.Sprintf("user=%v, params=%v", user, parameters))
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
	if username == "joe" && password == "secret" {
		return true
	}
	return false
}

func main() {
	e := echo.New()
	e.Use(middleware.BasicAuth(auth))
	e.Use(middleware.BodyLimit("80B"))
	e.Use(middleware.Logger())

	e.GET("/parameters/path/:user", pathParameters)
	e.GET("/parameters/query", queryParameters)
	e.PUT("/handling/request", handlingRequest)

	e.Run(standard.WithConfig(engine.Config{
		Address:     ":12345",
		TLSCertFile: "config/server-cert.pem",
		TLSKeyFile:  "config/server-key.pem",
	}))
}
