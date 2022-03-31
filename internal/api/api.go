package api

import (
	"fmt"
	"github.com/ExchangeRates/strategy-analyser/internal/controller"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server interface {
	BindingAddressFromPort(port int) string
}

type server struct {
	router     *mux.Router
	logger     *logrus.Logger
	controller *controller.AnalyserController
}

func NewServer(controller *controller.AnalyserController) *server {
	s := &server{
		router:     mux.NewRouter(),
		logger:     logrus.New(),
		controller: controller,
	}

	s.configureRouter()

	logrus.Info("starting api server")

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) BindingAddressFromPort(port int) string {
	return fmt.Sprintf(":%d", port)
}

func (s *server) configureRouter() {
	s.router.Path("/analyse/{strategy}").
		Handler(s.controller.HandleAnalyseStrategy()).
		Queries(
			"from", "{from}",
			"to", "{to}",
		).
		Methods(http.MethodPost)
}
