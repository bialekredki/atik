package handlers

import (
	"bialekredki/atik/pkg/auth"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginHandlerFactory(authenticator *auth.Authenticator) gin.HandlerFunc {

	return func(ctx *gin.Context) {
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
			authenticator.AuthCodeURL(state),
		)
	}
}

func LoginCallbackHandlerFactory(authenticator *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid State Parameter")
			return
		}

		token, err := authenticator.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token")
			return
		}

		idToken, err := authenticator.VerifyIDToken(ctx.Request.Context(), token)
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
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
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
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}
