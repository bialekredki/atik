package handlers

import (
	"bialekredki/atik/web/templates"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getProfileFromSession(session sessions.Session) map[string]interface{} {
	profile, ok := session.Get("profile").(map[string]interface{})
	if !ok {
		profile = make(map[string]interface{})
	}
	return profile
}

func (h *Handler) HomeHandler(ctx *gin.Context) {
	// var directories []models.MetadataDirectory
	// var files []models.MetadataFile
	session := sessions.Default(ctx)
	owner_id, ok := session.Get("owner_id").(uint)
	if !ok {
		ctx.Redirect(http.StatusTemporaryRedirect, "/welcome")
		return
	}
	h.logger.Debug("Home", zap.Uint("owner_id", owner_id), zap.Any("session", session))
	// directories, files = h.metadataRepository.ListContentsOfDirectory(nil, owner_id, 20)
	objects, err := h.metadataRepository.ListMetadataObjects(owner_id, nil, 20)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "", templates.Home(templates.HomeTemplateParams{
		Profile:   getProfileFromSession(session),
		Directory: nil,
		Objects:   objects,
	}))
}

func MainPageHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := getProfileFromSession(session)
	ctx.HTML(http.StatusOK, "", templates.Main(profile))
}

func AliveHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "I'm alive")
}
