package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type ResponseJSON struct {
	ReturnCode string `json:"returnCode"`
	Result     string `json:"result"`
}

type User struct {
	Name  string `json:"name" xml:"name" form:"name"`
	Email string `json:"email" xml:"email" form:"email"`
}

func PathParameters(c echo.Context) error {
	user := c.Param("user")
	paramNames := c.ParamNames()   // 路径参数名列表
	paramValues := c.ParamValues() // 路径参数值列表
	return c.String(http.StatusOK, fmt.Sprintf("user=%v, paramNames=%v, paramValues=%v",
		user, paramNames, paramValues))
}

func QueryParameters(c echo.Context) error {
	path := c.Path()              // 获取handler的注册路径
	user := c.QueryParam("user")  // 获取Query参数的第一个值
	parameters := c.QueryParams() //获取所有Query参数map
	return c.String(http.StatusOK, fmt.Sprintf("path=%v, user=%v, params=%v", path, user, parameters))
}

func HandlingRequest(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusExpectationFailed, ResponseJSON{"Failed", err.Error()})
	}
	fmt.Println(u)
	return c.JSON(http.StatusOK, ResponseJSON{"OK", "Handle Request"})
}
