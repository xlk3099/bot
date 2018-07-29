package utils_test

import (
	"testing"

	"github.com/xlk3099/ok-trading/utils"
)

func TestAdd(t *testing.T) {
	sma5 := utils.NewSMA(5)
	sma5.Add(1)
	if sma5.Avg() != 1 {
		t.Error(sma5.Avg())
	}
	sma5.Add(2)
	if sma5.Avg() != 1.5 {
		t.Error(sma5.Avg())
	}

	sma5.Add(3.0)
	sma5.Add(4.0)
	sma5.Add(5.0)
	sma5.Add(6.0)
	if sma5.Avg() != 4.0 {
		t.Error(sma5.Avg())
	}
}
