package api

import (
	"github.com/photoprism/photoprism/internal/authn"
	"github.com/photoprism/photoprism/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/acl"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/i18n"
	"github.com/photoprism/photoprism/internal/service"
)

// PUT /api/v1/users/:uid/password
func ChangePassword(router *gin.RouterGroup) {
	router.PUT("/users/:uid/password", func(c *gin.Context) {
		conf := service.Config()

		if conf.Public() || conf.DisableSettings() {
			Abort(c, http.StatusForbidden, i18n.ErrPublic)
			return
		}

		s := Auth(SessionID(c), acl.ResourcePeople, acl.ActionUpdateSelf)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		uid := c.Param("uid")
		m := entity.FindUserByUID(uid)

		if m == nil {
			Abort(c, http.StatusNotFound, i18n.ErrUserNotFound)
			return
		}

		f := form.ChangePassword{}

		if err := c.BindJSON(&f); err != nil {
			Error(c, http.StatusBadRequest, err, i18n.ErrInvalidPassword)
			return
		}

		if m.InvalidPassword(f.OldPassword) {
			Abort(c, http.StatusBadRequest, i18n.ErrInvalidPassword)
			return
		}

		if err := m.SetPassword(f.NewPassword); err != nil {
			Error(c, http.StatusBadRequest, err, i18n.ErrInvalidPassword)
			return
		}

		c.JSON(http.StatusOK, i18n.NewResponse(http.StatusOK, i18n.MsgPasswordChanged))
	})
}

func UserManagement(router *gin.RouterGroup) {
	conf := service.Config()

	// GET /api/v1/users Retrieves a list of all users
	router.GET("/users", func(c *gin.Context) {

	})

	// POST /api/v1/users Creates new user and returns session if no email confirmation && no admin confirmation
	router.POST("/users", func(c *gin.Context) {
		var f form.Register

		if err := c.BindJSON(&f); err != nil {
			AbortBadRequest(c)
			return
		}
		log.Info(f)
		user := &entity.User{
			RoleAdmin:    true, // TODO change back to false when implementing access control
			UserName:     f.UserName,
			FullName:     f.FullName,
			PrimaryEmail: f.Email,
			Password:     f.Password,
		}
		log.Info(user)

		if f.IdToken != "" {
			externalUid, err := authn.ValidateAndExtractID(f.IdToken)
			if err != nil {
				c.Error(err)
				return
			}
			user.ExternalUID = externalUid
			err = user.Create()
			if err != nil {
				c.Error(err)
				return
			}
			log.Infof("user sucessfully registered and linked: %s", user.UserName)

		} else {
			err := user.CreateAndValidate(conf.AuthConfig())
			if err != nil {
				c.Error(err)
				return
			}
			log.Infof("user sucessfully registered ")
		}
		// TODO only login directly if confirmation (admin/email) not required
		var data = session.Data{
			User: *user,
		}
		id := service.Session().Create(data)

		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": id, "data": data, "config": conf.UserConfig()})
		return
	})

	// GET /api/v1/users/:uid Retrieves a users info
	router.GET("/users/:uid", func(c *gin.Context) {

	})

	// PUT /api/v1/users/:uid Updates a users info
	router.PUT("/users/:uid", func(c *gin.Context) {

	})

	// DELETE /api/v1/users/:uid Deletes a user
	router.DELETE("/users/:uid", func(c *gin.Context) {

	})
}
