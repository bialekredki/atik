package handlers

import (
	"bialekredki/atik/pkg/auth"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) LoginHandler(ctx *gin.Context) {
	state, err := auth.GenerateRandomState()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	session := sessions.Default(ctx)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Redirect(
		http.StatusTemporaryRedirect,
		h.auth.AuthCodeURL(state),
	)
}

func (h *Handler) LoginCallbackHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if ctx.Query("state") != session.Get("state") {
		ctx.String(http.StatusBadRequest, "Invalid State Parameter")
		return
	}

	token, err := h.auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token")
		return
	}

	idToken, err := h.auth.VerifyIDToken(ctx.Request.Context(), token)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to verify ID Token")
		return
	}
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)

	userName := profile["name"].(string)
	h.s3repository.CreateBucket(userName, ctx.Request.Context())
	owner_id := h.metadataRepository.CreateNewOwner(userName)
	fmt.Printf("Owner ID = %d", owner_id)
	session.Set("owner_id", owner_id)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func LogoutHandler(ctx *gin.Context) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("OIDC_PROVIDER_DOMAIN") + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "http"
	}
	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host + "/logout/callback")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("OIDC_PROVIDER_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()
	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

func LogoutCallbackHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("profile")
	session.Delete("access_token")
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/welcome")
}
