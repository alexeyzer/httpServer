package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/alexeyzer/httpServer/internal/app/model"
	"github.com/alexeyzer/httpServer/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
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
	adv := &model.Adv{}
	listRef :=  []model.Ref{}

	adv.Name = r.FormValue("name")
	adv.Description = r.FormValue("description")
	adv.Price, _ = strconv.Atoi(r.FormValue("price"))
	if err := adv.Check(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	str := r.FormValue("reference")
	list := strings.Split(str, ",")

	newAdv, err := s.store.Adv().Create(adv)
	if err != nil {
		s.logger.Errorf("error adding new adverb: %v", err)
		return
	}
	if len(list) > 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Maximum count of references is 3"))
		return
	}
	for _,a := range list {
		newReference := model.Ref{}
		newReference.Ref = a
		newReference.AdvId = newAdv.ID
		listRef = append(listRef, newReference)
	}
	for _, elem := range listRef {
		_, err := s.store.Ref().Create(&elem)
		if err != nil {
			s.logger.Errorf("error adding new reference: %v", err)
			return
		}
	}
	response := fmt.Sprintf("{\"id\": %v, result: %v}", adv.ID, "success")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
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
		listRef, err := s.store.Ref().GetList(model.ID)
		if err != nil {
			s.logger.Errorf("error in pq: %v", err)
			return
		}
		if len(listRef) > 0{
			if optional == true  {
				model.Ref = listRef
			} else {
				model.Ref = listRef[:1]
			}
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
	p := r.FormValue("page")
	var n int = 1
	if p != ""{
		var err error
		n, err = strconv.Atoi(p)
		if err != nil {
			s.logger.Errorf("error converting page: %v", err)
			return
		}
	}
	list, err := s.store.Adv().List(sort)
	if err != nil {
		s.logger.Errorf("error in getting all adv: %v", err)
	}
	for c, a := range list {
		listRef, err := s.store.Ref().GetList(a.ID)
		if err != nil {
			s.logger.Errorf("error in pq: %v", err)
			return
		}
		if len(listRef) > 0{
			list[c].Ref = listRef[:1]
		}
	}
	page := model.PageAdv{}
	count := len(list)
	countPages := count / 10 + count % 10
	if count > 0 {
		if countPages < n {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Page"))
			return
		}
		if count < n*10 {
			if n == 1{
				n = 0
			}
			page.ListAdv = list[n*10 : count]
			page.NextPage = false
		} else {
			page.ListAdv = list[n*10 : n*10+10]
			page.NextPage = true
		}
	}
	response, err := json.Marshal(page)
	if err != nil {
		s.logger.Errorf("error in marshaling json: %v", err)
		return
	}
	s.logger.Info(string(response))
	w.Write(response)
}