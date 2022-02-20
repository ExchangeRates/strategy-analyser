package model

type Rate interface {
	GetIsBuy() bool
	GetRate() float64
}

type rateObj struct {
	isBuy bool
	rate  float64
}

func NewRate(isBuy bool, rate float64) Rate {
	return &rateObj{
		isBuy: isBuy,
		rate:  rate,
	}
}

func (r *rateObj) GetIsBuy() bool {
	return r.isBuy
}

func (r *rateObj) GetRate() float64 {
	return r.rate
}
