package httpserver

import (
	"encoding/json"
	"github.com/alexeyzer/httpServer/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	s.router.HandleFunc("/adv/list", s.listAdv)
	s.router.HandleFunc("/adv/current", s.currentAdv)
	s.router.HandleFunc("/adv/create", s.advCreate)
}

func (s *httpserver) advCreate(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("HandlerAdvCreate")

}

func (s *httpserver) currentAdv(w http.ResponseWriter, r *http.Request){
	s.logger.Info("HandlerCurrentAdv")

	fields := r.FormValue("fields")
	id := r.FormValue("id")
	var optional bool = false

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}
	advId, _ := strconv.Atoi(id)
	if fields != "" {
		optional = true
	}
	model, err := s.store.Adv().FindById(advId, optional)
	if err != nil {
		s.logger.Errorf("error in find pq: %v", err)
		return
	}
	if model != nil {
		model.NextPage = false
		listRef, err := s.store.Ref().GetList(model.ID)
		if err != nil {
			s.logger.Errorf("error in pq: %v", err)
			return
		}
		if optional == true  {
			model.Ref = listRef
		} else {
			model.Ref = listRef[:1]
		}
	}
	response, err := json.Marshal(model)
	if err != nil {
		s.logger.Errorf("error in marshaling json: %v", err)
		return
	}
	s.logger.Info(string(response))
	w.Write(response)
}

func (s *httpserver) listAdv(w http.ResponseWriter, r *http.Request){
	s.logger.Info("HandlerListAdv")

	sort := r.FormValue("sort")
	list, err := s.store.Adv().List(sort)
	if err != nil {
		s.logger.Errorf("error in getting all adv: %v", err)
	}
	for _, a := range list {
		listRef, err := s.store.Ref().GetList(a.ID)
		if err != nil {
			s.logger.Errorf("error in pq: %v", err)
			return
		}
		a.Ref =listRef[:1]
	}
	response, err := json.Marshal(list)
	if err != nil {
		s.logger.Errorf("error in marshaling json: %v", err)
		return
	}
	//count := len(list)
	s.logger.Info(string(response))
	w.Write(response)
}