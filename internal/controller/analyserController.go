package controller

import (
	"encoding/json"
	"github.com/ExchangeRates/strategy-analyser/internal/feign"
	"github.com/ExchangeRates/strategy-analyser/internal/service"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type AnalyserController struct {
	service service.AnalyserService
	feign   feign.RateResultFeign
}

func NewAnalyserController(service service.AnalyserService) *AnalyserController {
	return &AnalyserController{
		service: service,
	}
}

func (ac *AnalyserController) HandleAnalyseStrategy() http.HandlerFunc {
	dateLayout := "01.02.2006"
	type request struct {
		MajorStart float64 `json:"majorStart"`
		MinorStart float64 `json:"minorStart"`
		Step       float64 `json:"step"`
	}
	type response struct {
		Major float64 `json:"major"`
		Minor float64 `json:"minor"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body := &request{}
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			// TODO error respond
			return
		}

		strategy := mux.Vars(r)["strategy"]
		from, err := time.Parse(dateLayout, r.FormValue("from"))
		if err != nil {
			// TODO error respond
			return
		}
		to, err := time.Parse(dateLayout, r.FormValue("to"))
		if err != nil {
			// TODO error respond
			return
		}

		majorSum, minorSum, err := ac.service.Process(
			strategy,
			from,
			to,
			body.MajorStart,
			body.MinorStart,
			body.Step,
		)
		if err != nil {
			// TODO
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response{
			Major: majorSum,
			Minor: minorSum,
		}); err != nil {
			// TODO error respond
			return
		}
	}
}
