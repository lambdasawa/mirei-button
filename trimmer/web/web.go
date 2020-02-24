package web

import (
	"mb-trimmer/web/endpoint/sound"
	"mb-trimmer/web/endpoint/video"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/assets", "frontend/dist")
	e.Static("/", "frontend/dist/index.html")
	e.GET("video", video.Get)
	e.GET("sound", sound.Get)

	return e.Start(":3011")
}
