package handlers

import (
	"bialekredki/atik/web/templates"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HomeHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile, ok := session.Get("profile").(map[string]interface{})
	if !ok {
		profile = make(map[string]interface{})
	}
	ctx.HTML(http.StatusOK, "", templates.Home(profile))
}

func AliveHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "I'm alive")
}
