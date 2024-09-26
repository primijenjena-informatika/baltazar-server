package config

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/primijenjan-informatika/baltazar-server/controllers"
	"github.com/spf13/viper"
)

func StartServer() {
	e := echo.New()

	v1 := e.Group("/v1")

	v1.POST("/chat/completions", controllers.HandleMessages)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("api.port"))))
}
