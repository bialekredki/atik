package main

import (
	"bialekredki/atik/config"
	"bialekredki/atik/lib/aws"
	"bialekredki/atik/lib/metadata"
	"bialekredki/atik/lib/renderer"
	"bialekredki/atik/pkg/auth"
	"bialekredki/atik/web/handlers"
	"context"
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func bootstrapRoutes(engine *gin.Engine, handler *handlers.Handler) {
	engine.GET("/status", handlers.AliveHandler)
	engine.GET("/welcome", handlers.MainPageHandler)
	engine.GET("/login", handler.LoginHandler)
	engine.GET("/login/callback", handler.LoginCallbackHandler)
	engine.GET("/logout/callback", handlers.LogoutCallbackHandler)

	authorized := engine.Group("/")
	{
		authorized.GET("/", handler.HomeHandler)
		authorized.GET("/directory/:objectId", handler.GetDirectoryRow)
		authorized.GET("/logout", handlers.LogoutHandler)
		authorized.POST("/directory", handler.CreateDirectory)
	}
}

func Server(lc fx.Lifecycle, handler *handlers.Handler) *gin.Engine {
	engine := gin.Default()

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			store := cookie.NewStore([]byte(os.Getenv("COOKIE_SESSION_SECRET_KEY")))
			engine.Use(sessions.Sessions("auth-session", store))

			engine.SetTrustedProxies(nil)
			bootstrapRoutes(engine, handler)
			engine.Static("/public", "./assets")

			ginHtmlRenderer := engine.HTMLRender
			engine.HTMLRender = &renderer.HTMLTemplRenderer{
				FallbackHtmlRenderer: ginHtmlRenderer,
			}
			go func() {
				server.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
	return engine
}

func main() {
	godotenv.Load(".env")
	gob.Register(map[string]interface{}{})

	app := fx.New(
		fx.Provide(
			Server,
			aws.LoadConfig,
			aws.NewS3Repository,
			handlers.NewHandler,
			auth.New,
			config.DatabaseConnection,
			metadata.NewMetadataRepository,
			zap.NewDevelopment,
		),
		fx.Invoke(func(*gin.Engine) {}),
	)
	app.Run()
}
