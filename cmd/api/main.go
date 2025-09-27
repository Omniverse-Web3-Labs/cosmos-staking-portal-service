package main

import (
	"app/internal/cache/redisclient"
	"app/internal/config"
	"app/internal/db/postgres"
	"app/internal/loader"
	"app/internal/middleware"
	"app/internal/router"
	"app/pkg/server"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	config.Parse("")
	redisclient.Load()
	serverConfig := config.GetConfig().ServerApi
	otelOptions := config.GetOtelOptions()
	gin.SetMode(serverConfig.Mode)
	postgres.Load()

	server := server.New(serverConfig.Name, serverConfig.Listen, otelOptions, func(engine *gin.Engine) {
		loader.Initialize()
		engine.Use(middleware.Cors())
		engine.Use(otelgin.Middleware(serverConfig.Name))
		engine.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"ip":   c.ClientIP(),
				"time": time.Now().UnixMilli(),
			})
		})
		v1 := engine.Group("/api/v1")
		authV1 := engine.Group("/api/v1")
		authV1.Use(middleware.Auth())
		router.RegisterApiV1(v1)
		router.RegisterAuthApiV1(authV1)
	})
	server.Start()
}
