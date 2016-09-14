package main

import (
	"uniagent/common/net/restful"

	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
)

func main() {
	router := restful.NewRouter(routes)
	router.Run(standard.WithConfig(engine.Config{
		Address:     ":12345",
		TLSCertFile: "config/server-cert.pem",
		TLSKeyFile:  "config/server-key.pem",
	}))
}
