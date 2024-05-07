package main

import (
	"bialekredki/atik/config"
	"bialekredki/atik/lib/aws"
	"bialekredki/atik/lib/renderer"
	"bialekredki/atik/models"
	"bialekredki/atik/pkg/auth"
	"bialekredki/atik/web/handlers"
	"encoding/gob"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	engine := gin.Default()
	connection := config.DatabaseConnection()

	if err := connection.AutoMigrate(
		&models.MetadataOwner{},
		&models.MetadataFile{},
		&models.MetadataDirectory{},
	); err != nil {
		panic(err)
	}

	gob.Register(map[string]interface{}{})
	aws.LoadConfig()

	store := cookie.NewStore([]byte(os.Getenv("COOKIE_SESSION_SECRET_KEY")))
	engine.Use(sessions.Sessions("auth-session", store))

	authenticator, err := auth.New()
	if err != nil {
		panic(err)
	}

	ginHtmlRenderer := engine.HTMLRender
	engine.Static("/public", "./assets")
	engine.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	engine.SetTrustedProxies(nil)

	engine.GET("/status", handlers.AliveHandler)
	engine.GET("/login", handlers.LoginHandlerFactory(authenticator))
	engine.GET("/login/callback", handlers.LoginCallbackHandlerFactory(authenticator))
	engine.GET("/logout/callback", handlers.LogoutCallbackHandler)

	authorized := engine.Group("/")

	{
		authorized.GET("/", handlers.HomeHandler)
		authorized.GET("/logout", handlers.LogoutHandler)
	}

	engine.Run()
}
