package service

import "github.com/ExchangeRates/strategy-analyser/internal/model"

type ProcessStrategy interface {
	Process(actions []model.Rate, major, minor, step float64) (float64, float64)
}

type processStrategyImpl struct{}

func NewProcessStrategy() ProcessStrategy {
	return &processStrategyImpl{}
}

func (p *processStrategyImpl) Process(actions []model.Rate, major, minor, step float64) (float64, float64) {
	for _, action := range actions {
		rate := action.GetRate()
		if action.GetIsBuy() {
			major += step * rate
			minor -= step / rate
		} else {
			major -= step * rate
			minor += step / rate
		}
	}

	return major, minor
}
