package httpserver

import (
	"github.com/alexeyzer/httpServer/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// httpServer
type httpserver struct {
	config *config
	logger *logrus.Logger
	router *mux.Router
	store *store.Store
}

func NewHttpServer(config *config) *httpserver {
	server := httpserver{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
	return &server
}

func (s *httpserver) Start() error{
	err := s.ConfigLogger()
	if err != nil {
		return err
	}
	s.ConfigureRouter()
	if err := s.ConfigStore(); err != nil{
		return err
	}
	s.logger.Info("Staring httpserver")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *httpserver) ConfigLogger() error{
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *httpserver) ConfigStore() error{
	store := store.NewStore(s.config.Store)
	if err := store.Open(); err != nil {
		return err
	}

	s.store = store
	return nil
}

func (s *httpserver) ConfigureRouter() {
	s.router.HandleFunc("/hello", handleHello())
}

func handleHello() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		io.WriteString(w, "First response in go")
	}
}