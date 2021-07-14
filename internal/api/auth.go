package api

import (
	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/authn"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/internal/session"
	"io/ioutil"
	"net/http"
)

// GET /api/v1/auth/external
func AuthEndpoints(router *gin.RouterGroup) {
	conf := service.Config()
	authn.Init()

	router.GET("/auth/external", func(c *gin.Context) {
		//clientConfig := conf.PublicConfig()
		//randomIdentifier := "random23456hjklb8"
		//url := authn.OauthConfig.AuthCodeURL(authn.OauthStateString)
		url, err := authn.StartAuthFlow(c.Writer, c.Request)
		log.Infof("External Auth Url %s", url)
		//c.Header("X-Session-Wait", oauthStateString)
		c.Redirect(http.StatusTemporaryRedirect, url)

		// try to get the user without re-authenticating
		//if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		//	t, _ := template.New("foo").Parse(userTemplate)
		//	t.Execute(res, gothUser)
		//} else {
		//	gothic.BeginAuthHandler(res, req)
		//}
	})

	router.GET("/auth/callback", func(c *gin.Context) {
		//clientConfig := conf.PublicConfig()
		//randomIdentifier := "fsado8yghtfds9hy5r"
		st := c.Query("state") //c.PostForm("state")
		log.Infof("State from Callback %s", st)
		if st != authn.OauthStateString {
			log.Errorf("callback not valid")
		}
		token, err := authn.OauthConfig.Exchange(c, c.Query("code"))

		if err != nil {
			log.Errorf("couldn't get token")
		}
		log.Infof("Access Token\n%s\n", token.AccessToken)
		log.Infof("Refresh Token\n%s\n", token.RefreshToken)
		log.Infof("ID Token\n%s\n", token.Extra("id_token"))
		log.Infof("Token valid\n%v\n", token.Valid())

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, "https://keycloak.timovolkmann.de/auth/realms/master/protocol/openid-connect/userinfo", nil)
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)
		//token.SetAuthHeader(req)
		res, err := client.Do(req)
		if err != nil {
			log.Errorf("UserInfo: %s", err)
		}
		defer res.Body.Close()
		contents, _ := ioutil.ReadAll(res.Body)
		log.Infof("UserInfo: %s", contents)

		// if token.Valid() then retrieve User with matching external UID
		// if token not valid, show error message
		// if no user found start linkUser/newUser flow by setting flag in localstorage
		// if user found return new session, set it in browser localstorage and close popup
		var data = session.Data{}
		data.User = *entity.FindUserByName("timo008")
		id := service.Session().Create(data)

		//c.Data(http.StatusOK, "application/json", contents)
		//clientConfig := conf.PublicConfig()
		c.HTML(http.StatusOK, "callback.tmpl", gin.H{"status": "ok", "id": id, "data": data, "config": conf.UserConfig()})
	})
}
