package web

import (
	"mb-trimmer/web/endpoint/register"
	"mb-trimmer/web/endpoint/sound"
	"mb-trimmer/web/endpoint/twitter"
	"mb-trimmer/web/endpoint/video"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("MB_SESSION_SECRET")))))

	e.Static("/assets", "frontend/dist")
	e.Static("/", "frontend/dist/index.html")
	e.GET("/twitter/sign-in", twitter.SignIn)
	e.GET("/twitter/callback", twitter.Callback)
	e.GET("/twitter/status", twitter.Status)
	e.GET("/twitter/sign-out", twitter.SignOut)
	e.GET("/video", video.Get)
	e.GET("/sound", sound.Get)
	e.POST("/register", register.Register)
	e.POST("/bulk-register", register.BulkRegister)

	return e.Start(":3011")
}
