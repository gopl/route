package main

import (
	"fmt"

	"github.com/iris-contrib/middleware/basicauth"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name"`
	Email string `json:"email" xml:"email" form:"email"`
}

type Response struct {
	ReturnCode string `json:"returnCode"`
	Result     string `json:"result"`
}

func pathParameters(ctx *iris.Context) {
	user := ctx.Param("user") // 路径参数
	params := ctx.Params
	ctx.Write("user=%v, Params=%v", user, params)
}

func queryParameters(ctx *iris.Context) {
	user := ""                   // 没法直接获取请求参数的值？
	queryArgs := ctx.QueryArgs() // 从RequestURI获取请求参数
	ctx.Write("user=%v, QueryArgs", user, queryArgs)
}

func handlingRequest(ctx *iris.Context) {
	var u User
	if err := ctx.ReadJSON(&u); err != nil {
		panic(err.Error())
	}
	fmt.Println(u)
	ctx.JSON(iris.StatusOK, Response{"OK", "Handle Request"})
}

func main() {
	auth := basicauth.Default(
		map[string]string{
			"user": "password",
		})
	iris.Use(logger.New(iris.Logger))
	iris.Use(auth)
	iris.Get("/parameters/path/:user", pathParameters)
	iris.Get("/parameters/query", queryParameters)
	iris.Put("/handling/request", handlingRequest)

	iris.ListenTLS(":12345", "config/server-cert.pem", "config/server-key.pem")
}
