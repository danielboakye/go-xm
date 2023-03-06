package main

import (
	"net/http"
	"time"

	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/handlers"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const testUUID = "af7c1fe6-d669-414e-b066-e9733f0de7a8"

type HTTPServer struct {
	handler IHTTPHandler
	cfg     config.Configurations
}

type IHTTPHandler interface {
	http.Handler
	Run(addr ...string) error
}

func newHTTPHandler(handler *handlers.Handler, cfg config.Configurations) IHTTPHandler {
	//  Creates a handler without any middleware by default
	engine := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	engine.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	engine.Use(gin.Recovery())

	engine.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:8080"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300 * time.Second,
	}))

	// Generate token for test
	engine.GET("/api/v1/token", func(ctx *gin.Context) {
		accessToken, err := helpers.GenerateAccessToken(cfg, testUUID)
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
	private.POST("", handler.CreateCompany)
	private.PATCH("/:company-id", handler.UpdateCompany)
	private.DELETE("/:company-id", handler.DeleteCompany)

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Resource not found",
		})
	})

	return engine
}

func newHTTPServer(handler IHTTPHandler, cfg config.Configurations) HTTPServer {
	return HTTPServer{handler: handler, cfg: cfg}
}

func (s HTTPServer) Start() error {
	err := s.handler.Run(":" + s.cfg.HTTPPort)
	return err
}

// Stop Server
