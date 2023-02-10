package main

import (
	"net/http"

	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/handlers"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/middleware"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
	cfg    config.Configurations
}

func NewServerHTTP(handler *handlers.Handler, cfg config.Configurations) *ServerHTTP {
	//  Creates a handler without any middleware by default
	engine := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	engine.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	engine.Use(gin.Recovery())

	// Generate token for test
	engine.GET("/api/v1/token", func(ctx *gin.Context) {
		accessToken, err := helpers.GenerateAccessToken(cfg, "1")
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token": accessToken,
		})
	})

	public := engine.Group("/api/v1/company")
	public.GET("/:company-id", handler.GetCompany)

	// Mount protected handlers
	private := engine.Group("/api/v1/company", middleware.NewRouteFilter(cfg))
	private.POST("/", handler.CreateCompany)
	private.PATCH("/:company-id", handler.UpdateCompany)
	private.DELETE("/:company-id", handler.DeleteCompany)

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Resource not found",
		})
	})

	return &ServerHTTP{engine: engine, cfg: cfg}
}

func (sh *ServerHTTP) Start() {
	err := sh.engine.Run(":" + sh.cfg.HTTPPort)
	if err != nil {
		panic(err)
	}
}
