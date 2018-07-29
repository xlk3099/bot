package utils

// code  is from https://github.com/mxmCherry/movavg/blob/master/sma.go

// SMA - Simple moving average
type SMA struct {
	window int // Simple Moving Average window (number of latest values to calculate average value)

	vals []float64 // stored values, circular list
	n    int       // number of actually stored values (vals slice usage)
	i    int       // last-set value index; first-set (oldest) value is (i+1)%window when n == window

	avg float64 // current Simple Moving Average value
}

// NewSMA constructs new Simple Moving Average calculator.
// Window arg must be >= 1.
func NewSMA(window int) *SMA {
	if window <= 0 {
		panic("movavg.NewSMA: window should be > 0")
	}
	return &SMA{
		window: window,
		vals:   make([]float64, window),
	}
}

// Add recalculates Simple Moving Average value and returns it.
func (a *SMA) Add(v float64) float64 {
	if a.n == a.window {
		// filled window - most frequent case:
		// https://en.wikipedia.org/wiki/Moving_average#Simple_moving_average
		a.i = (a.i + 1) % a.window
		a.avg += (v - a.vals[a.i]) / float64(a.n)
		a.vals[a.i] = v
	} else if a.n != 0 {
		// partially-filled window - second most frequent case:
		// https://en.wikipedia.org/wiki/Moving_average#Cumulative_moving_average
		a.i = (a.i + 1) % a.window
		a.avg = (v + float64(a.n)*a.avg) / float64(a.n+1)
		a.vals[a.i] = v
		a.n++
	} else {
		// empty window - least frequent case (occurs only once, on first value added):
		// simply assign given value as current average:
		a.avg = v
		a.vals[0] = v
		a.n = 1
	}
	return a.avg
}

// Avg returns current Simple Moving Average value.
func (a *SMA) Avg() float64 {
	return a.avg
}
