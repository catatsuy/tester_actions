package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/catatsuy/bento/config"
	"github.com/catatsuy/bento/util"
)

var (
	templates *template.Template
)

func init() {
	templates = template.Must(template.ParseFiles(
		"./public/translate.html",
	))
}

type reqTranslate struct {
	Input string `json:"input"`
}

type resTranslate struct {
	Output string `json:"output"`
}

func (s *Server) postTranslate(w http.ResponseWriter, r *http.Request) {
	rt := reqTranslate{}

	err := json.NewDecoder(r.Body).Decode(&rt)
	if err != nil {
		outputErrorMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	input := util.TrimUnnecessary(rt.Input)

	if len(input) == 0 {
		return
	}

	oldnew, newold := config.Replacer(s.words)
	replacerNoun := strings.NewReplacer(oldnew...)
	revertNoun := strings.NewReplacer(newold...)

	input = replacerNoun.Replace(input)

	conf, exist, err := config.LoadCache()
	if err != nil {
		outputError(w, http.StatusInternalServerError, err)
		return
	}

	token := ""
	if !exist {
		token, err = s.Session.GetToken()
		if err != nil {
			outputError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		s.Session.SetCacheCookie(conf.Cookies)
		token = conf.Token
	}
	s.Session.SetToken(token)

	isJP := util.AutoDetectJP(input)

	output, err := s.Session.PostTranslate(input, isJP)
	if err != nil {
		outputError(w, http.StatusInternalServerError, err)
		return
	}

	output = revertNoun.Replace(output)

	if !exist {
		ccs := s.Session.DumpCookies()
		err = config.DumpCache(config.Cache{
			Cookies: ccs,
			Token:   token,
		})
		if err != nil {
			outputError(w, http.StatusInternalServerError, err)
			return
		}
	}

	res := resTranslate{
		Output: output,
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(res)
}

func (s *Server) getTranslate(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "translate.html", struct{}{})
}
