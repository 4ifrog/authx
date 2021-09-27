package api

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/store"
)

func GetRoutesFunc() server.RegisterRoutesFunc {
	return func(router *gin.Engine, cfg *config.Config, ds store.DataStore) {
		// Initialization.
		tmpl := NewTemplate(cfg.TemplatesDir)
		trans := NewEnglishTranslator()
		validate := NewValidate(trans)

		htmlHandlers := NewHTMLHandlers(cfg, ds, trans, validate)
		authHandlers := NewAuthHandlers(cfg, ds)
		errHandlers := NewErrorHandlers(cfg)
		middleware := NewMiddleware(cfg, ds)

		router.SetHTMLTemplate(tmpl)
		router.Use(errHandlers.ErrorResponse())

		// Public web pages.
		webGrp := router.Group("/")
		webGrp.GET("/", htmlHandlers.SignIn())
		webGrp.POST("/", htmlHandlers.SignIn())
		webGrp.GET(errorsPath, errHandlers.ErrorHTML())

		// Protected web pages.
		proWebGrp := router.Group("/")
		proWebGrp.Use(middleware.SetContextFromCookie())
		proWebGrp.GET("/userinfo", htmlHandlers.UserInfo())
		proWebGrp.POST("/userinfo", htmlHandlers.UserInfo())

		// Public auth api.
		apiGrp := router.Group("/v1")
		apiGrp.POST("/signin", authHandlers.SignIn())
		apiGrp.POST("/signout", authHandlers.SignOut())
		apiGrp.GET("/avatar/:identity", authHandlers.Avatar())

		// Protected auth api.
		proAPIGrp := router.Group("/v1")
		proAPIGrp.Use(middleware.SetContextFromBearerAuth())
		proAPIGrp.GET("/userinfo", authHandlers.UserInfo())

		// Fallback to static content.
		router.Use(static.Serve("/", static.LocalFile(cfg.StaticWebDir, false)))
	}
}
