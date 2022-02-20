package service

import (
	"github.com/ExchangeRates/strategy-analyser/internal/feign"
	"github.com/ExchangeRates/strategy-analyser/internal/model"
	"time"
)

type AnalyserService interface {
	Process(strategy string, from, to time.Time, major, minor, step float64) (float64, float64, error)
}

type analyserServiceImpl struct {
	feign     feign.RateResultFeign
	processor ProcessStrategy
}

func NewAnalyserService(feign feign.RateResultFeign, processor ProcessStrategy) AnalyserService {
	return &analyserServiceImpl{
		feign:     feign,
		processor: processor,
	}
}

func (a *analyserServiceImpl) Process(strategy string, from, to time.Time, major, minor, step float64) (float64, float64, error) {

	count, onePageSize, err := a.feign.Count(strategy, from, to)
	if err != nil {
		return 0, 0, err
	}

	actions := make([]model.Rate, 0)
	for i := 0; i < 1+count/onePageSize; i++ {
		loaded, err := a.feign.Actions(strategy, i, from, to)
		if err != nil {
			return 0, 0, nil
		}
		for _, action := range loaded {
			actions = append(actions, action)
		}
	}

	majorRes, minorRes := a.processor.Process(
		actions,
		major,
		minor,
		step,
	)

	return majorRes, minorRes, nil
}
