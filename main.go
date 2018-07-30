package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	"github.com/xlk3099/ok-trading/ok"
	"github.com/xlk3099/ok-trading/utils"
)

var addr = flag.String("addr", "real.okex.com:10440", "http service address")

func init() {
	// init logrus
	formatter := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	}

	log.SetFormatter(formatter)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	filenameHook := filename.NewHook()
	filenameHook.Field = "source"
	log.AddHook(filenameHook)
}

func main() {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	tradeEMA5()
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	updateEMA1()
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	updateEMA5()
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	updateEMA30()
	// }()
	// wg.Wait()
}

func trade() {
	// 比较ma20, ma
}

// func updateEMA1() {
// 	// calculate ema1
// 	btc := ok.NewPair("bch_usd", "this_week", "", "")
// 	ema := func() *utils.Ema {
// 		klines := btc.GetFutureKlineData("1min")
// 		ema := utils.NewEma(12)
// 		for _, k := range klines {
// 			ema.Add(k.TimeStamp, k.Close)
// 		}
// 		return ema
// 	}
// 	log.Info("1min:", ema().Latest())
// 	ticker := time.NewTicker(1 * time.Minute)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			ema := ema()
// 			log.Info("1min:", ema.Latest())
// 		}
// 	}
// }

func tradeEMA5() {
	// calculate EMA5
	// calculate ema1
	etc := ok.NewPair("etc_usd", "this_week", "", "")
	fma := func() *utils.Ema {
		klines := etc.GetFutureKlineData("5min")
		fma := utils.NewEma(12)
		for _, k := range klines {
			fma.Add(k.TimeStamp, k.Close)
		}
		return fma
	}
	sma := func() *utils.Ema {
		klines := etc.GetFutureKlineData("5min")
		sma := utils.NewEma(50)
		for _, k := range klines {
			sma.Add(k.TimeStamp, k.Close)
		}
		return sma
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ema12 := fma()
			ema50 := sma()
			fpr := doGetFurturePos4Fix(etc)
			userInfo := etc.GetFutureUserInfo4Fix()
			amtToTrade := int(userInfo.Info.Etc.Balance / 5 * 20)
			ft := etc.GetFutureTicker()
			if utils.IsGoldCross(ema12, ema50, ft.Ticker.Last) {
				log.Info("卧槽，金叉了。。。")
				if len(fpr.Holdings) > 0 {
					amtClose := fpr.Holdings[0].SellAvailable
					if amtClose > 0 {
						success := doTrade(etc, utils.Float64ToString(ft.Ticker.Buy), strconv.Itoa(amtClose), ok.CloseShort, true)
						if success {
							log.Error("大爷止损成功")
						} else {
							log.Error("大爷止损失败")
						}
					}
					if cha := fpr.Holdings[0].BuyAmount; cha < amtToTrade*3/4 {
						success := doTrade(etc, utils.Float64ToString(ft.Ticker.Sell), strconv.Itoa(amtToTrade-cha), ok.Long, true)
						if success {
							log.Info("大爷增单成功：", amtToTrade-cha)
						} else {
							log.Error("大爷增单失败")
						}
						continue
					}
				}
				success := doTrade(etc, utils.Float64ToString(ft.Ticker.Sell), strconv.Itoa(amtToTrade), ok.Long, true)
				if success == true {
					log.Info("大爷开多成功,张数：", amtToTrade)
				} else {
					log.Info("大爷开多失败....很遗憾，检查程序bug吧。。。")
				}
			}

			if utils.IsDeadCross(ema12, ema50, ft.Ticker.Last) {
				log.Info("卧槽，死叉了。。。")
				ft := etc.GetFutureTicker()
				if len(fpr.Holdings) > 0 {
					amtClose := fpr.Holdings[0].BuyAvailable
					if amtClose > 0 {
						success := doTrade(etc, utils.Float64ToString(ft.Ticker.Buy), strconv.Itoa(amtClose), ok.CloseLong, true)
						if success {
							log.Info("大爷止损成功")
						} else {
							log.Error("大爷止损失败")
						}
					}
					if cha := fpr.Holdings[0].SellAmount; cha < amtToTrade*3/4 {
						success := doTrade(etc, utils.Float64ToString(ft.Ticker.Sell), strconv.Itoa(amtToTrade-cha), ok.Short, true)
						if success {
							log.Info("大爷增单成功：", amtToTrade-cha)
						} else {
							log.Error("大爷增单失败")
						}
						continue
					}
				}
				success := doTrade(etc, utils.Float64ToString(ft.Ticker.Sell), strconv.Itoa(amtToTrade), ok.Short, true)
				if success {
					log.Info("大爷开空成功", amtToTrade)
				} else {
					log.Info("大爷开空失败....很遗憾，检查程序bug吧。。。")
				}
			}
		}
	}
}

// func updateEMA30() {
// 	// calculate EMA30
// 	btc := ok.NewPair("bch_usd", "this_week", "", "")
// 	ema := func() *utils.Ema {
// 		klines := btc.GetFutureKlineData("30min")
// 		ema := utils.NewEma(12)
// 		for _, k := range klines {
// 			ema.Add(k.TimeStamp, k.Close)
// 		}
// 		return ema
// 	}
// 	log.Info("30min:", ema().Latest())
// 	ticker := time.NewTicker(1 * time.Minute)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			ema := ema()
// 			log.Info("30min:", ema.Latest())
// 		}
// 	}
// }

// func updateEMA60() {
// 	// calculate EMA60
// 	btc := ok.NewPair("bch_usd", "this_week", "", "")

// 	btc.GetFutureKlineData("5min")
// 	ticker := time.NewTicker(60 * time.Minute)
// 	defer ticker.Stop()
// }

func doGetFurturePos4Fix(pair *ok.Pair) *ok.FuturePosResp {
	var futurePos *ok.FuturePosResp
	var errNotFound = errors.New("unable to get current holdings")
	err := utils.Do(func(attempt int) (bool, error) {
		var err error
		futurePos, err = pair.GetFuturePos4Fix()
		if len(futurePos.Holdings) < 0 {
			err = errNotFound
		}
		time.Sleep(300 * time.Millisecond)
		return attempt < 10, err // try 10 times
	})

	if err != nil {
		log.WithError(err).Error("unable to get future current future holdings")
		return nil
	}
	return futurePos
}

func doTrade(pair *ok.Pair, price string, amount string, tradeType ok.TradeType, matching bool) bool {
	var errTradeFail = errors.New("failed to trade")
	err := utils.Do(func(attempt int) (bool, error) {
		var err error
		resp := pair.FutureTrade(price, amount, tradeType, matching)
		if !resp.Result {
			err = errTradeFail
		}
		time.Sleep(300 * time.Millisecond)
		return attempt < 10, err // try 5 times
	})
	if err != nil {
		return false
	}
	return true
}
