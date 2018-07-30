package ok_test

import (
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/xlk3099/ok-trading/ok"
	"github.com/xlk3099/ok-trading/utils"
)

var etcPair = ok.NewPair("etc_usd", "this_week", "", "")

func TestGetFuturePos(t *testing.T) {
	ftr, _ := etcPair.GetFuturePos4Fix()
	log.Info(len(ftr.Holdings))
	log.Info(ftr.Holdings[0].BuyAvailable)
}

// func TestGetFutureTicker(t *testing.T) {
// 	ft := etcPair.GetFutureTicker()
// 	log.Info("买一价：", ft.Ticker.Buy)
// 	log.Info("卖一价:", ft.Ticker.Sell)
// }

func TestFutureTrade(t *testing.T) {
	ft := etcPair.GetFutureTicker()
	etcPair.FutureTrade(utils.Float64ToString(ft.Ticker.Sell), strconv.Itoa(1), ok.Long, true)
	ftr, _ := etcPair.GetFuturePos4Fix()
	log.Info(ftr)
	log.Info(len(ftr.Holdings))
}

func TestGetFutureUserInfo4Fix(t *testing.T) {
	etcInfo := etcPair.GetFutureUserInfo4Fix()
	log.Info(etcInfo.Info)
}
