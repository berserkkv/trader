package calculator

import (
	"testing"

	"github.com/berserkkv/trader/model/enum/order"
	"github.com/stretchr/testify/assert"
)

func TestCalculateStopLoss(t *testing.T) {
	assert.Equal(t, 110.0, CalculateStopLoss(100.0, 10.0, order.SHORT))
	assert.Equal(t, 90.0, CalculateStopLoss(100.0, 10.0, order.LONG))
}

func TestCalculateTakeProfit(t *testing.T) {
	assert.Equal(t, 90.0, CalculateTakeProfit(100, 10, order.SHORT))
	assert.Equal(t, 110.0, CalculateTakeProfit(100, 10, order.LONG))
}

func TestCalculateBuyQuantity(t *testing.T) {
	assert.Equal(t, 10.0, CalculateBuyQuantity(10, 100))
	assert.Equal(t, 0.0, CalculateBuyQuantity(10, 0))
	assert.Equal(t, 0.0, CalculateBuyQuantity(0, 100))
}

func TestCalculateTakerFee(t *testing.T) {
	assert.Equal(t, 0.04, CalculateTakerFee(100))
	assert.Equal(t, 0.0, CalculateTakerFee(0))
}

func TestCalculateMakerFee(t *testing.T) {
	assert.Equal(t, 0.02, CalculateMakerFee(100))
	assert.Equal(t, 0.0, CalculateMakerFee(0))
}

func TestCalculatePNL(t *testing.T) {
	assert.Equal(t, 10.0, CalculatePNL(101, 1000, 10, order.LONG))
	assert.Equal(t, 100.0, CalculatePNL(90, 1000, 10, order.SHORT))
	assert.Equal(t, 0.0, CalculatePNL(100, 1000, 10, order.SHORT))

}

func TestCalculateRoe(t *testing.T) {
	assert.InDelta(t, 10.0, CalculateRoe(100, 101, 10, order.LONG), 1e-6)
	assert.InDelta(t, 10, CalculateRoe(100, 99, 10, order.SHORT), 1e-6)
}
