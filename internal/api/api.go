package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ExchangeRates/strategy-analyser/internal/controller"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Server interface {
	GracefullListenAndServe(port int)
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

func (s *server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down")
	return nil
}

func (s *server) bindingAddressFromPort(port int) string {
	return fmt.Sprintf(":%d", port)
}

func (s *server) GracefullListenAndServe(port int) error {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return http.ListenAndServe(s.bindingAddressFromPort(port), s)
	})
	g.Go(func() error {
		<-gCtx.Done()
		return s.Shutdown(context.Background())
	})

	return g.Wait()
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
