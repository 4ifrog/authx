package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/store"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const (
	signinTmplName  = "signin"
	profileTmplName = "userinfo"
	successURI      = "/userinfo"
	rootURI         = "/"
)

type HTMLHandlers struct {
	trans    ut.Translator
	validate *validator.Validate
	cfg      *config.Config
	ds       store.DataStore
}

var (
	ErrUserNotFound       = errors.New("can't find user")
	ErrMethodNotSupported = errors.New("method not supported")
	ErrUserIDCast         = errors.New("can't cast user id to string type")
)

func NewHTMLHandlers(cfg *config.Config, ds store.DataStore,
	trans ut.Translator, validate *validator.Validate) *HTMLHandlers {
	handlers := new(HTMLHandlers)

	handlers.cfg = cfg
	handlers.ds = ds
	handlers.trans = trans
	handlers.validate = validate

	return handlers
}

// TODO: Refactor too many if-else statements.

func (hh *HTMLHandlers) SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// GET  = displays the page.
		// POST = handles the form submission.
		if ctx.Request.Method == http.MethodGet {
			ctx.HTML(http.StatusOK, signinTmplName, nil)
		} else if ctx.Request.Method == http.MethodPost {
			var msg strings.Builder

			var signin SignInUser
			if err := ctx.ShouldBind(&signin); err != nil {
				vErrs, ok := err.(validator.ValidationErrors)
				if !ok {
					log.Panicf("failed to cast validator.ValidationErrors: %v", err)
				}

				for _, e := range vErrs {
					msg.WriteString(fmt.Sprintln(e.Translate(hh.trans)))
				}
			} else {
				user, err := auth.Authenticate(ctx, hh.ds, signin.Username, signin.Password)
				if err == auth.ErrUserNotFound {
					msg.WriteString("User not found")
				} else if err == auth.ErrInvalidCredentials {
					msg.WriteString("Invalid credentials")
				} else if err != nil {
					msg.WriteString(fmt.Sprintf("Internal error: %s", err))
				}

				// Generate oauth2 object and save.
				if msg.Len() == 0 {
					aTTL := time.Duration(hh.cfg.AccessTTL) * time.Second
					rTTL := time.Duration(hh.cfg.RefreshTTL) * time.Second
					otoken, err := auth.CreateOAuthToken(ctx, hh.ds, user.ID, hh.cfg.AccessSecret, aTTL, rTTL)
					if err != nil {
						msg.WriteString(fmt.Sprintf("Internal error: %s", err))
					} else {
						// Save token to the cookie.
						token := SessionToken{
							Token:  *otoken,
							UserID: user.ID,
						}
						ss := NewCookieStore(hh.cfg.SessionSecret)
						if err := ss.SetSessionToken(ctx.Writer, ctx.Request, &token); err != nil {
							msg.WriteString(fmt.Sprintf("Internal error: %s", err))
						}
					}
				}
			}

			if msg.Len() > 0 {
				content := gin.H{
					"Error": msg.String(),
				}

				ctx.HTML(http.StatusOK, signinTmplName, content)

				return
			}

			// Redirect if successful
			ctx.Redirect(http.StatusMovedPermanently, successURI)
		} else {
			// Other methods
			setErrorStatus(ctx, ErrMethodNotSupported, http.StatusMethodNotAllowed)
			return
		}
	}
}

func (hh *HTMLHandlers) UserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ss := NewCookieStore(hh.cfg.SessionSecret)

		// GET  = displays the page.
		// POST = handles the form submission.
		if ctx.Request.Method == http.MethodGet {
			val, ok := ctx.Get(keyUserID)
			if !ok {
				setErrorStatus(ctx, ErrUserNotFound, http.StatusUnauthorized)
				return
			}

			uid, ok := val.(string)
			if !ok {
				setErrorStatus(ctx, ErrUserIDCast, http.StatusInternalServerError)
				return
			}

			user, err := hh.ds.GetUser(ctx, uid)
			if err == store.ErrorNotFound {
				setErrorStatus(ctx, ErrUserNotFound, http.StatusUnauthorized)
				return
			} else if err != nil {
				setErrorStatus(ctx, err, http.StatusInternalServerError)
				return
			}

			content := gin.H{
				"Username": user.Username,
			}

			ctx.HTML(http.StatusOK, profileTmplName, content)
		} else if ctx.Request.Method == http.MethodPost {
			if err := ss.ClearSessionToken(ctx.Writer, ctx.Request); err != nil {
				setErrorStatus(ctx, err, http.StatusInternalServerError)
				return
			}

			// Redirect if successful
			ctx.Redirect(http.StatusMovedPermanently, rootURI)
		}
	}
}
