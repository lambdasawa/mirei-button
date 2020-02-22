package web

import (
	"mb-trimmer/web/endpoint/kirinuki"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("kirinuki", kirinuki.Get)

	return e.Start(":3011")
}
