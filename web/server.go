package web

import (
	"fmt"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	vv9 "gopkg.in/go-playground/validator.v9"
)

type WebServer struct {
	*echo.Echo
	host string
	port int
}

func NewWebServer(host string, port int) (*WebServer, error) {
	webServer := new(WebServer)
	webServer.host = host
	webServer.port = port

	e := echo.New()
	// 自定义错误
	e.HTTPErrorHandler = customHTTPErrorHandler
	// 验证
	e.Validator = &CustomValidator{
		validator: vv9.New(),
	}
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		ExposeHeaders: []string{"*"},
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowMethods:  []string{"*"},
	}))
	// RequestID
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = middleware.DefaultRequestIDConfig.Generator()
				c.Request().Header.Set(echo.HeaderXRequestID, requestID)
			}
			return next(c)
		}
	})

	webServer.Echo = e

	return webServer, nil
}

func (w *WebServer) Run() error {
	w.RegisterMiddleware()
	w.RegisterRouting()

	addr := fmt.Sprintf("%s:%d", w.host, w.port)
	w.Server.Addr = addr
	w.HideBanner = true
	w.Debug = true
	w.Logger.Fatal(gracehttp.Serve(w.Server))
	//w.Logger.Fatal(w.Start(w.Server.Addr))
	return nil
}

func (w *WebServer) RegisterMiddleware() {
	w.Use(middleware.Logger())
	w.Use(middleware.Recover())
	w.Use(middleware.Secure())
	w.Use(middleware.RequestID())
}

func (w *WebServer) RegisterRouting() {}

func customHTTPErrorHandler(e error, ctx echo.Context) {}
