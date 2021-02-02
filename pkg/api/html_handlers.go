package api

import (
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/store"
)

const (
	successURI = "/profile.html"
	failureURI = "/failure.html"
)

func renderTemplate(w io.Writer, tmplPath string, data interface{}) error {
	// Load the template file.
	tmpl := template.Must(template.ParseFiles(tmplPath))

	// Render the template file.
	return tmpl.Execute(w, data)
}

func HomeHTMLHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := renderTemplate(ctx.Writer, "public/index.html", nil); err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// SignInHTMLHandler displays the sign-in form.
func SignInHTMLHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session.
		ss := NewSessionStore(cfg.SessionSecret)
		s, err := ss.GetSession(ctx.Request)
		if err != nil {
			ctx.Redirect(http.StatusFound, failureURI)
			return
		}
		if s != nil {
			claims, err := auth.ParseJWT(s.OAuth2Token.AccessToken, cfg.AccessSecret)
			if err != nil {
				ctx.Redirect(http.StatusFound, failureURI)
				return
			}
			if claims.Expiry.Time().After(time.Now()) {
				ctx.Redirect(http.StatusFound, successURI)
			}
		}

		if err := renderTemplate(ctx.Writer, "public/signin.html", nil); err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// SignInFormHTMLHandler process the sign-in form submission.
func SignInFormHTMLHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var errMsg strings.Builder

		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		if username == "" {
			errMsg.WriteString("username is empty")
		}
		if password == "" {
			errMsg.WriteString("password is empty")
		}

		user, err := auth.Authenticate(ctx, ds, username, password)
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidCredentials {
			ctx.Redirect(http.StatusFound, failureURI)
			return
		} else if err != nil {
			ctx.Redirect(http.StatusFound, failureURI)
			return
		}

		if errMsg.Len() > 0 {
			content := &struct {
				Error string
			}{
				Error: errMsg.String(),
			}

			if terr := renderTemplate(ctx.Writer, "public/signin.html", content); terr != nil {
				ctx.Redirect(http.StatusInternalServerError, failureURI)
				return
			}

			return
		}

		// Generate oauth2 object and save.
		aTTL := time.Duration(cfg.AccessTTL) * time.Second
		rTTL := time.Duration(cfg.RefreshTTL) * time.Second
		otoken, err := auth.CreateOAuthToken(ctx, ds, user.ID, cfg.AccessSecret, aTTL, rTTL)
		if err != nil {
			ctx.Redirect(http.StatusInternalServerError, failureURI)
			return
		}

		// Save session to the cookie.
		session := UserSession{
			OAuth2Token: *otoken,
		}
		ss := NewSessionStore(cfg.SessionSecret)
		if err := ss.SetSession(ctx.Writer, ctx.Request, &session); err != nil {
			ctx.Redirect(http.StatusInternalServerError, failureURI)
			return
		}

		ctx.Redirect(http.StatusFound, successURI)
	}
}

// ProfileHTMLHandler displays the sign-in form.
func ProfileHTMLHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := renderTemplate(ctx.Writer, "public/profile.html", nil); err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// SignoutFormHTMLHandler process the sign-out form submission.
func SignoutFormHTMLHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session.
		ss := NewSessionStore(cfg.SessionSecret)
		if err := ss.ClearSession(ctx.Writer, ctx.Request); err != nil {
			ctx.Redirect(http.StatusFound, failureURI)
			return
		}

		ctx.Redirect(http.StatusFound, "/index.html")
	}
}
