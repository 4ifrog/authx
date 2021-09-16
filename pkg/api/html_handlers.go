package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"
)

const (
	signinTmplName  = "signin"
	profileTmplName = "profile"
	err401TmplName  = "401"
	successURI      = "/profile"
	rootURI         = "/"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	tmpl     *template.Template
)

func initValidation() {
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

func initTemplates(tmplDir string) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.gohtml", tmplDir))
	if err != nil {
		log.Panicf("failed to get files in %s: %v", tmplDir, err)
	}

	// Load the template file.
	tmpl = template.Must(template.ParseFiles(files...))
}

func renderTemplate(ctx *gin.Context, tmplName string, data interface{}) {
	if err := tmpl.ExecuteTemplate(ctx.Writer, tmplName, data); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func InitHTMLHandlers(cfg *config.Config) {
	initValidation()
	initTemplates(cfg.TemplatesDir)
}

func WebSignInHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// GET  = displays the page.
		// POST = handles the form submission.
		if ctx.Request.Method == http.MethodGet {
			renderTemplate(ctx, signinTmplName, nil)
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
						UserID:      user.ID,
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

				renderTemplate(ctx, signinTmplName, content)

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

func WebProfileOutHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ss := NewSessionStore(cfg.SessionSecret)

		// GET  = displays the page.
		// POST = handles the form submission.
		if ctx.Request.Method == http.MethodGet {
			var content interface{}

			uid, ok := ctx.Get(keyUserID)
			if !ok {
				renderTemplate(ctx, err401TmplName, nil)
				return
			}

			user, err := ds.GetUser(ctx, uid.(string))
			if err != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			if user == nil {
				renderTemplate(ctx, err401TmplName, nil)
				return
			}

			content = &struct {
				Username string
			}{
				Username: user.Username,
			}

			renderTemplate(ctx, profileTmplName, content)
		} else if ctx.Request.Method == http.MethodPost {
			if err := ss.ClearSession(ctx.Writer, ctx.Request); err != nil {
				fmt.Printf("failed to clear session: %v", err)
			}

			// Redirect if successful
			http.Redirect(ctx.Writer, ctx.Request, rootURI, http.StatusFound)
		}
	}
}
