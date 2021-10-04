package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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

func handleValidationError(err error, trans ut.Translator) []string {
	var msgs []string

	vErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		log.Panicf("failed to cast validator.ValidationErrors: %v", err)
		return msgs
	}

	for _, e := range vErrs {
		msgs = append(msgs, e.Translate(trans))
	}

	return msgs
}

func handleErrorMessages(ctx *gin.Context, tmplName string, msgs []string) {
	if len(msgs) > 0 {
		content := gin.H{
			"ErrorMessages": msgs,
		}

		ctx.HTML(http.StatusOK, tmplName, content)
		return
	}

	// Redirect if successful
	ctx.Redirect(http.StatusMovedPermanently, successURI)
}

func (hh *HTMLHandlers) postSignIn(ctx *gin.Context) []string {
	var msgs []string

	var signin SignInRequest
	if err := ctx.ShouldBind(&signin); err != nil {
		return handleValidationError(err, hh.trans)
	}

	// Check if user exists
	user, err := hh.ds.GetUserByUsername(ctx, signin.Username)
	if user == nil || err == store.ErrorNotFound {
		return append(msgs, "User not found")
	} else if err != nil {
		return append(msgs, fmt.Sprintf("Internal error: %s", err))
	}

	// Authenticate
	if ok := auth.Authenticate(user.Password, signin.Password, user.Salt); !ok {
		return append(msgs, "Invalid credentials")
	}

	// Strip sensitive data like password.
	user.RemoveSensitiveData()

	// Generate oauth2 object and save.
	aTTL := time.Duration(hh.cfg.AccessTTL) * time.Second
	rTTL := time.Duration(hh.cfg.RefreshTTL) * time.Second
	otoken, err := auth.CreateOAuthToken(ctx, hh.ds, user.ID, hh.cfg.AccessSecret, aTTL, rTTL)
	if err != nil {
		return append(msgs, fmt.Sprintf("Internal error: %s", err))
	}

	// Save token to the cookie.
	token := SessionToken{
		Token:  *otoken,
		UserID: user.ID,
	}
	ss := NewCookieStore(hh.cfg.SessionSecret)
	if err := ss.SetSessionToken(ctx.Writer, ctx.Request, &token); err != nil {
		return append(msgs, fmt.Sprintf("Internal error: %s", err))
	}

	return msgs
}

func (hh *HTMLHandlers) SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// GET  = displays the page.
		// POST = handles the form submission.
		if ctx.Request.Method == http.MethodGet {
			ctx.HTML(http.StatusOK, signinTmplName, nil)
			return
		} else if ctx.Request.Method == http.MethodPost {
			msgs := hh.postSignIn(ctx)

			handleErrorMessages(ctx, signinTmplName, msgs)
			return
		}

		setErrorStatus(ctx, ErrMethodNotSupported, http.StatusMethodNotAllowed)
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
