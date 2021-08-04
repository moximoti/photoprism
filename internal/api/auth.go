package api

import (
	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/authn"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/internal/session"
	"net/http"
)

// GET /api/v1/auth/
func AuthEndpoints(router *gin.RouterGroup) {
	conf := service.Config()

	ap := conf.AuthProvider()
	if ap == authn.ProviderNone || ap == "" {
		return
	}
	if err := authn.Init(conf); err != nil {
		log.Errorf(err.Error())
	}

	router.GET("/auth/reset_password", func(c *gin.Context) {
		_ = c.Query("code")
		_ = c.Query("username")
		// TODO user.resetPassword(code)
		// and implement frontend success/error message
		c.Redirect(http.StatusTemporaryRedirect, "/login?msg=password_success")
	})

	router.GET("/auth/activate_user", func(c *gin.Context) {
		_ = c.Query("code")
		_ = c.Query("username")
		// TODO user.activate(code)
		// and implement frontend success/error message
		c.Redirect(http.StatusTemporaryRedirect, "/login?msg=activation_error")
	})

	router.GET("/auth/external", func(c *gin.Context) {
		err := authn.StartAuthFlow(c.Writer, c.Request)
		if err != nil {
			log.Errorf("External Auth Error: %s", err.Error())
		}
	})

	router.GET("/auth/callback", func(c *gin.Context) {
		// "SignInOAuthCallback"
		userInfo, err := authn.FinalizeAuthFlow(c.Writer, c.Request)
		if err != nil {
			log.Errorf(err.Error())
			c.HTML(http.StatusUnauthorized, "callback.tmpl", gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		log.Infof("UserInfo: %s %s", userInfo.Email, userInfo.UserID)
		log.Debugf("IDToken: %s", userInfo.IDToken)
		log.Debugf("AToken: %s", userInfo.AccessToken)

		user := entity.FindUserByExternalUID(userInfo.UserID)
		if user == nil {
			c.HTML(http.StatusOK, "callback.tmpl", gin.H{
				"status": "linkUser",
				"linkUser": map[string]string{
					"NickName": userInfo.NickName,
					"Name":     userInfo.Name,
					"Email":    userInfo.Email,
					"IdToken":  userInfo.IDToken,
				},
				"config": conf.UserConfig(),
			})
			return
		}
		log.Infof("user '%s' logged in", user.UserName)
		var data = session.Data{
			User: *user,
		}
		id := service.Session().Create(data)
		c.HTML(http.StatusOK, "callback.tmpl", gin.H{
			"status": "ok",
			"id":     id,
			"data":   data,
			"config": conf.UserConfig(),
		})

	})
}
