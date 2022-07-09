package internal

import (
	"github.com/ExchangeRates/strategy-analyser/internal/api"
	"github.com/ExchangeRates/strategy-analyser/internal/config"
	"github.com/ExchangeRates/strategy-analyser/internal/controller"
	"github.com/ExchangeRates/strategy-analyser/internal/feign"
	"github.com/ExchangeRates/strategy-analyser/internal/service"
)

func Start(config *config.Config) error {

	feignRateResults := feign.NewRateResultFeign(config.RateResultUrl)
	processorService := service.NewProcessStrategy()
	analyseService := service.NewAnalyserService(feignRateResults, processorService)
	analyseController := controller.NewAnalyserController(analyseService)

	srv := api.NewServer(analyseController)

	return srv.GracefullListenAndServe(config.Port)
}
