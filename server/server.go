package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/catatsuy/bento/config"
	"github.com/catatsuy/bento/mirait"
)

type Server struct {
	AppVersion string
	Session    *mirait.Session

	words []string
	mux   *http.ServeMux
}

func New(appVersion string) *Server {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Print(err)
		return nil
	}

	sess, err := mirait.NewSession(cfg)
	if err != nil {
		log.Print(err)
		return nil
	}

	words, err := config.LoadWords()
	if err != nil {
		log.Print(err)
		return nil
	}

	s := &Server{
		AppVersion: appVersion,
		Session:    sess,
		words:      words,
	}

	s.mux = http.NewServeMux()

	s.mux.HandleFunc("/api/translate", s.postTranslate)

	s.mux.HandleFunc("/api/version", s.versionHandler)

	// Register pprof handlers
	runtime.SetBlockProfileRate(1)
	s.mux.HandleFunc("/debug/pprof/", pprof.Index)
	s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.mux.Handle(("/js/"), http.FileServer(http.Dir("./public/")))
	s.mux.HandleFunc("/translate.html", s.getTranslate)

	return s
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	b, _ := json.Marshal(struct {
		Version string `json:"version"`
	}{Version: s.AppVersion})

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func outputErrorMsg(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	b, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: msg})

	w.WriteHeader(status)
	w.Write(b)
}

func outputError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	b, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: err.Error()})

	w.WriteHeader(status)
	w.Write(b)
}
