package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	j "github.com/helloeave/json"
)

type ServerConfig struct {
	Port int
}

func (c *ServerConfig) Validate() error {
	if c.Port == 0 {
		return errors.New("port is required")
	}
	return nil
}

type Server struct {
	port   int
	Router *chi.Mux
}

func NewServer(config *ServerConfig) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RealIP,
		middleware.RequestID,
		middleware.Timeout(60*time.Second),
		Recoverer,
	)

	return &Server{
		port:   config.Port,
		Router: router,
	}, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.Router)
}

func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	b, err := j.MarshalSafeCollections(
		payload,
	) // forked encoding/json repo that turns nil [] and {} into [] and {}
	if err != nil {
		panic(Error{http.StatusInternalServerError, err})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

var decoder = schema.NewDecoder()

func Parse(r *http.Request, out interface{}) {
	if r.Method == "GET" {
		if err := decoder.Decode(out, r.URL.Query()); err != nil {
			Abort(http.StatusUnprocessableEntity, err)
		}
	} else {
		var b bytes.Buffer
		bodyCpy := io.TeeReader(r.Body, &b)
		if err := json.NewDecoder(bodyCpy).Decode(out); err != nil {
			Abort(http.StatusUnprocessableEntity, err)
		}
		r.Body = io.NopCloser(&b)
	}
}

type Validator interface {
	Validate() error
}

func ParseAndValidate(r *http.Request, out Validator) {
	Parse(r, out)
	if err := out.Validate(); err != nil {
		log.Printf("err caught in ParseAndValidate %v\n", err)
		Abort(http.StatusUnprocessableEntity, err)
	}
}
