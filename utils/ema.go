package utils

import "github.com/ethereum/go-ethereum/log"

var limit = 10000

// Ema is the main object
type Ema struct {
	period int
	points []point
}

type point struct {
	Timestamp int64
	Value     float64
	Ema       float64
}

// NewEma creates new ema instance.
func NewEma(period int) *Ema {
	ema := &Ema{period: period, points: make([]point, 0)}
	return ema
}

// GetPoints returns all the points in this ema up to recent 10k points.
func (ema *Ema) GetPoints() []point {
	return ema.points
}

// Add adds a new Value to Ema
func (ema *Ema) Add(timestamp int64, value float64) {
	p := point{Timestamp: timestamp, Value: value}

	alpha := 2.0 / (float64(ema.period) + 1.0)

	// fmt.Println(alpha)

	emaTminusOne := value
	if len(ema.points) > 0 {
		emaTminusOne = ema.points[len(ema.points)-1].Ema
	}

	emaT := alpha*value + (1-alpha)*emaTminusOne
	p.Ema = emaT
	ema.points = append(ema.points, p)

	if len(ema.points) > limit {
		ema.points = ema.points[limit/2:]
	}
}

// Last2 returns last2 ema value.
func (ema *Ema) Last2() float64 {
	return ema.points[len(ema.points)-2].Ema
}

// Last3 returns last ema value.
func (ema *Ema) Last3() float64 {
	return ema.points[len(ema.points)-3].Ema
}

// Current returns current ema value.
func (ema *Ema) Current() float64 {
	return ema.points[len(ema.points)-1].Ema
}

// IsGoldCross check if two ema gold cross or not.
func IsGoldCross(fma *Ema, sma *Ema, currentPrice float64) bool {
	if fma.Last3() < sma.Last3() {
		if fma.Last2() > sma.Last2() {
			log.Info("fma 金叉前:", fma.Last3(), " sma 金叉前:", sma.Last3())
			log.Info("fma 最新:", fma.Last2(), " sma 最新:", sma.Last3())
			return true
		}
	}
	return false
}

// IsDeadCross checks if two ema death cross or not.
func IsDeadCross(fma *Ema, sma *Ema, currentPrice float64) bool {
	if fma.Last3() > sma.Last3() {
		if fma.Last2() < sma.Last2() {
			log.Info("fma 死叉前:", fma.Last3(), " sma 死叉前:", sma.Last3())
			log.Info("fma 最新:", fma.Last2(), " sma 最新:", sma.Last3())
			return true
		}
	}
	return false
}
