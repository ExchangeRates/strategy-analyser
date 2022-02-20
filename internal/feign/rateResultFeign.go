package feign

import (
	"encoding/json"
	"fmt"
	"github.com/ExchangeRates/strategy-analyser/internal/model"
	"log"
	"net/http"
	"time"
)

type RateResultFeign interface {
	Actions(strategy string, page int, from, to time.Time) ([]model.Rate, error)
	Count(strategy string, from, to time.Time) (int, int, error)
}

type rateResultFeignImpl struct {
	url        string
	dateLayout string
}

func NewRateResultFeign(url string) RateResultFeign {
	return &rateResultFeignImpl{
		url:        url,
		dateLayout: "01.02.2006",
	}
}

type action struct {
	Major    string    `json:"major"`
	Minor    string    `json:"minor"`
	Action   string    `json:"action"`
	Rate     float64   `json:"rate"`
	Created  time.Time `json:"created"`
	Strategy string    `json:"strategy"`
}

func (rr *rateResultFeignImpl) Actions(strategy string, page int, from, to time.Time) ([]model.Rate, error) {

	url := fmt.Sprintf(
		"%s/rate-actions/%s?page=%d&from=%s&to=%s",
		rr.url, strategy, page,
		from.Format(rr.dateLayout), to.Format(rr.dateLayout),
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	actions := &[]action{}
	if err := json.NewDecoder(resp.Body).Decode(actions); err != nil {
		return nil, err
	}

	result := make([]model.Rate, 0)
	for _, action := range *actions {
		rate := model.NewRate(
			action.Action == "put",
			action.Rate,
		)
		result = append(result, rate)
	}

	return result, nil
}

type count struct {
	Count       int `json:"count"`
	OnePageSize int `json:"onePageSize"`
}

func (rr *rateResultFeignImpl) Count(strategy string, from, to time.Time) (int, int, error) {

	url := fmt.Sprintf(
		"%s/rate-actions/%s/count?from=%s&to=%s",
		rr.url, strategy,
		from.Format(rr.dateLayout), to.Format(rr.dateLayout),
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	count := &count{}
	if err := json.NewDecoder(resp.Body).Decode(count); err != nil {
		return 0, 0, err
	}

	return count.Count, count.OnePageSize, nil
}
