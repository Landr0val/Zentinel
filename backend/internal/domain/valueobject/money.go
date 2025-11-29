package valueobject

import (
	"errors"
	"math"
)

type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}
	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) Round() Money {
	ratio := math.Pow(10, 2)
	return Money{
		amount:   math.Round(m.amount*ratio) / ratio,
		currency: m.currency,
	}
}
