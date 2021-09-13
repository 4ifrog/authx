package api

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
)

const (
	signinTmplFile = "signin.html"
	successURI     = "/success"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

//nolint:gochecknoinits
func init() {
	locale := "en"
	english := en.New()
	uni = ut.New(english, english)
	trans, ok := uni.GetTranslator(locale)
	if !ok {
		log.Panicf("failed to get translator for %s locale", locale)
	}

	validate, ok = binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Panicf("failed to cast to *validator.Validate")
	}
	if err := ent.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Panicf("failed to register validation translator: %v", err)
	}
}

func renderTemplate(cfg *config.Config, w io.Writer, tmplFile string, data interface{}) error {
	// Load the template file.
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("%s/%s", cfg.TemplatesDir, tmplFile)))

	// Render the template file.
	return tmpl.Execute(w, data)
}

func WebSignInHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// GET  = displays the sign-in form.
		// POST = handles the sign-in form submission.
		if ctx.Request.Method == http.MethodGet {
			if err := renderTemplate(cfg, ctx.Writer, signinTmplFile, nil); err != nil {
				http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if ctx.Request.Method == http.MethodPost {
			var msg strings.Builder

			var login models.User
			if err := ctx.ShouldBind(&login); err != nil {
				vErrs, ok := err.(validator.ValidationErrors)
				if !ok {
					log.Panicf("failed to cast validator.ValidationErrors: %v", err)
				}

				trans, _ := uni.GetTranslator("en")

				for _, e := range vErrs {
					msg.WriteString(fmt.Sprintln(e.Translate(trans)))
				}
			} else {
				user, err := auth.Authenticate(ctx, ds, login.Username, login.Password)
				if err == auth.ErrUserNotFound {
					msg.WriteString("User not found")
				} else if err == auth.ErrInvalidCredentials {
					msg.WriteString("Invalid credentials")
				} else if err != nil {
					msg.WriteString(fmt.Sprintf("Internal error: %s", err))
				}

				// Generate oauth2 object and save.
				aTTL := time.Duration(cfg.AccessTTL) * time.Second
				rTTL := time.Duration(cfg.RefreshTTL) * time.Second
				otoken, err := auth.CreateOAuthToken(ctx, ds, user.ID, cfg.AccessSecret, aTTL, rTTL)
				if err != nil {
					msg.WriteString(fmt.Sprintf("Internal error: %s", err))
				} else {
					// Save session to the cookie.
					session := UserSession{
						OAuth2Token: *otoken,
					}
					ss := NewSessionStore(cfg.SessionSecret)
					if err := ss.SetSession(ctx.Writer, ctx.Request, &session); err != nil {
						msg.WriteString(fmt.Sprintf("Internal error: %s", err))
					}
				}
			}

			if msg.Len() > 0 {
				content := &struct {
					Error string
				}{
					Error: msg.String(),
				}

				if err := renderTemplate(cfg, ctx.Writer, signinTmplFile, content); err != nil {
					http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
					return
				}

				return
			}

			// Redirect if successful
			http.Redirect(ctx.Writer, ctx.Request, successURI, http.StatusFound)
		} else {
			// Other methods
			http.Error(ctx.Writer, fmt.Sprintf("%s not supported", ctx.Request.Method), http.StatusNotImplemented)
			return
		}
	}
}
