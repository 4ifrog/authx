package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"
)

func NewValidate(trans ut.Translator) *validator.Validate {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Panicf("failed to cast to *validator.Validate")
	}
	if err := ent.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Panicf("failed to register validation translator: %v", err)
	}

	return validate
}

func NewEnglishTranslator() ut.Translator {
	locale := "en"
	english := en.New()
	uni := ut.New(english, english)
	trans, ok := uni.GetTranslator(locale)
	if !ok {
		log.Panicf("failed to get %s translator for translating validation errors", locale)
	}

	return trans
}

func NewTemplate(tmplDir string) *template.Template {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.gohtml", tmplDir))
	if err != nil {
		log.Panicf("failed to get files in %s: %v", tmplDir, err)
	}

	// Load the template file.
	return template.Must(template.ParseFiles(files...))
}

func setErrorStatus(ctx *gin.Context, err error, code int) {
	ctx.Status(code)
	_ = ctx.Error(err)
	ctx.Next()
}

func isAcceptingHTML(r *http.Request) bool {
	accepts := r.Header["Accept"]
	for _, accept := range accepts {
		for _, mime := range strings.Split(accept, ",") {
			switch mime {
			case "text/html", "application/xhtml+xml", "application/xml":
				return true
			}
		}
	}

	return false
}
