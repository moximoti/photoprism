package server

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/api"
	"github.com/photoprism/photoprism/internal/config"
)

func registerRoutes(router *gin.Engine, conf *config.Config) {
	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// Static assets like js, css and font files.
	router.Static(conf.BaseUri(config.StaticUri), conf.StaticPath())
	router.StaticFile(conf.BaseUri("/favicon.ico"), filepath.Join(conf.ImgPath(), "favicon.ico"))

	// PWA Manifest.
	router.GET(conf.BaseUri("/manifest.json"), func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Content-Type", "application/json")

		clientConfig := conf.PublicConfig()
		c.HTML(http.StatusOK, "manifest.json", gin.H{"config": clientConfig})
	})

	// PWA Service Worker.
	router.GET(conf.BaseUri("/sw.js"), func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.File(filepath.Join(conf.BuildPath(), "sw.js"))
	})

	// Rainbow Page.
	router.GET(conf.BaseUri("/rainbow"), func(c *gin.Context) {
		clientConfig := conf.PublicConfig()
		c.HTML(http.StatusOK, "rainbow.tmpl", gin.H{"config": clientConfig})
	})

	/*	// Redirect to external Identity Provider.
		var (
			oauthConfig = &oauth2.Config{
				RedirectURL:    "http://localhost:2342/auth/callback",
				ClientID:     "photoprism-dev",
				ClientSecret: "341e8af4-4ed7-40cc-bd2c-b29a5e0cd40c",
				Scopes:       []string{"profile", "email", "openid"},
				Endpoint:     oauth2.Endpoint{
					AuthURL:   "https://keycloak.timovolkmann.de/auth/realms/master/protocol/openid-connect/auth",
					TokenURL:  "https://keycloak.timovolkmann.de/auth/realms/master/protocol/openid-connect/token",
					AuthStyle: 0,
				},
			}
			// Some random string, random for each request
			oauthStateString = "random"
		)
		if service.Config().Settings().Auth.AuthProvider == config.ProviderNone {

		}
		router.GET("/auth/external", func(c *gin.Context) {
		})
		router.GET("/auth/callback", func(c *gin.Context) {
			//clientConfig := conf.PublicConfig()
			//randomIdentifier := "fsado8yghtfds9hy5r"
			st := c.Query("state") //c.PostForm("state")
			log.Infof("State from Callback %s", st)
			if st != oauthStateString {
				log.Errorf("callback not valid")
			}
			token, err := oauthConfig.Exchange(c, c.Query("code"))

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
	*/
	// JSON-REST API Version 1
	v1 := router.Group(conf.BaseUri(config.ApiUri))
	{
		api.GetStatus(v1)
		api.GetErrors(v1)

		api.GetConfig(v1)
		api.GetConfigOptions(v1)
		api.SaveConfigOptions(v1)

		api.GetSettings(v1)
		api.SaveSettings(v1)

		api.ChangePassword(v1)
		api.CreateSession(v1)
		api.DeleteSession(v1)

		api.GetThumb(v1)
		api.GetDownload(v1)
		api.GetVideo(v1)
		api.CreateZip(v1)
		api.DownloadZip(v1)

		api.GetGeo(v1)
		api.GetPhoto(v1)
		api.GetPhotoYaml(v1)
		api.UpdatePhoto(v1)
		api.GetPhotos(v1)
		api.GetPhotoDownload(v1)
		api.GetPhotoLinks(v1)
		api.CreatePhotoLink(v1)
		api.UpdatePhotoLink(v1)
		api.DeletePhotoLink(v1)
		api.ApprovePhoto(v1)
		api.LikePhoto(v1)
		api.DislikePhoto(v1)
		api.AddPhotoLabel(v1)
		api.RemovePhotoLabel(v1)
		api.UpdatePhotoLabel(v1)
		api.GetMomentsTime(v1)
		api.GetFile(v1)
		api.DeleteFile(v1)
		api.UpdateFileMarker(v1)
		api.PhotoPrimary(v1)
		api.PhotoUnstack(v1)

		api.GetLabels(v1)
		api.UpdateLabel(v1)
		api.GetLabelLinks(v1)
		api.CreateLabelLink(v1)
		api.UpdateLabelLink(v1)
		api.DeleteLabelLink(v1)
		api.LikeLabel(v1)
		api.DislikeLabel(v1)
		api.LabelCover(v1)

		api.GetFoldersOriginals(v1)
		api.GetFoldersImport(v1)
		api.GetFolderCover(v1)

		api.Upload(v1)
		api.StartImport(v1)
		api.CancelImport(v1)
		api.StartIndexing(v1)
		api.CancelIndexing(v1)

		api.BatchPhotosApprove(v1)
		api.BatchPhotosArchive(v1)
		api.BatchPhotosRestore(v1)
		api.BatchPhotosPrivate(v1)
		api.BatchPhotosDelete(v1)
		api.BatchAlbumsDelete(v1)
		api.BatchLabelsDelete(v1)

		api.GetAlbum(v1)
		api.CreateAlbum(v1)
		api.UpdateAlbum(v1)
		api.DeleteAlbum(v1)
		api.DownloadAlbum(v1)
		api.GetAlbums(v1)
		api.GetAlbumLinks(v1)
		api.CreateAlbumLink(v1)
		api.UpdateAlbumLink(v1)
		api.DeleteAlbumLink(v1)
		api.LikeAlbum(v1)
		api.DislikeAlbum(v1)
		api.AlbumCover(v1)
		api.CloneAlbums(v1)
		api.AddPhotosToAlbum(v1)
		api.RemovePhotosFromAlbum(v1)

		api.GetAccounts(v1)
		api.GetAccount(v1)
		api.GetAccountFolders(v1)
		api.ShareWithAccount(v1)
		api.CreateAccount(v1)
		api.DeleteAccount(v1)
		api.UpdateAccount(v1)

		api.SendFeedback(v1)

		api.GetSvg(v1)

		api.Websocket(v1)

		api.AuthEndpoints(v1)
	}

	// Configure link sharing.
	s := router.Group(conf.BaseUri("/s"))
	{
		api.Shares(s)
		api.SharePreview(s)
	}

	// WebDAV server for file management, sync and sharing.
	if conf.DisableWebDAV() {
		log.Info("webdav: server disabled")
	} else {
		WebDAV(conf.OriginalsPath(), router.Group(conf.BaseUri(WebDAVOriginals), BasicAuth()), conf)
		log.Infof("webdav: %s/ enabled, waiting for requests", conf.BaseUri(WebDAVOriginals))

		if conf.ImportPath() != "" {
			WebDAV(conf.ImportPath(), router.Group(conf.BaseUri(WebDAVImport), BasicAuth()), conf)
			log.Infof("webdav: %s/ enabled, waiting for requests", conf.BaseUri(WebDAVImport))
		}
	}

	// Default HTML page for client-side rendering and routing via VueJS.
	router.NoRoute(func(c *gin.Context) {
		clientConfig := conf.PublicConfig()
		c.HTML(http.StatusOK, conf.TemplateName(), gin.H{"config": clientConfig})
	})
}
