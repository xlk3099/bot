package ok

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var okEp = "https://www.okex.com/api/v1"

// TradeType is a string alias
type TradeType = string

const (
	// Long 做多
	Long TradeType = "1"
	// Short 做空
	Short TradeType = "2"
	// CloseLong 平多
	CloseLong TradeType = "3"
	// CloseShort 平空
	CloseShort TradeType = "4"
)

// Pair is a ok restful api handler.
type Pair struct {
	symbol       string
	contractType string
	apiKey       string
	apiSign      string
	c            *http.Client
}

// FutureTicker represents the futureTicker
type FutureTicker struct {
	Date   string `json:"date"`
	Ticker ticker `json:"ticker"`
}

// ticker represents ticker
type ticker struct {
	Last       float64 `json:"last"`
	Buy        float64 `json:"buy"`  // 买一价
	Sell       float64 `json:"sell"` // 卖一价
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Vol        float64 `json:"vol"` // 24 小时成交量
	ContractID int64   `json:"contract_id"`
	UnitAmount float64 `json:"unit_amount"`
}

// Kline K线历史数据
type Kline struct {
	TimeStamp  int64
	Open       float64
	Highest    float64
	Lowest     float64
	Close      float64
	UnitAmount int
	Amount     float64
}

// FuturePosResp represents the futurepo response.
type FuturePosResp struct {
	Result   bool      `json:"result"`
	Holdings []Holding `json:"holding"`
}

// Holding represents future pos
type Holding struct {
	BuyAmount           int     `json:"buy_amount"`
	BuyAvailable        int     `json:"buy_available"`
	BuyBond             float64 `json:"buy_bond"`
	BuyFlatprice        string  `json:"buy_flatprice"`
	BuyPriceAvg         float64 `json:"buy_price_avg"`
	BuyPriceCost        float64 `json:"buy_price_cost"`
	BuyProfitLossratio  string  `json:"buy_profit_lossratio"`
	BuyProfitReal       float64 `json:"buy_profit_real"`
	ContractID          int64   `json:"contract_id"`
	ContractType        string  `json:"contract_type"`
	CreateDate          int64   `json:"create_date"`
	SellAmount          int     `json:"sell_amount"`
	SellAvailable       int     `json:"sell_available"`
	SellBond            float64 `json:"sell_bond"`
	SellFlatprice       string  `json:"sell_flatprice"`
	SellPriceAvg        float64 `json:"sell_price_avg"`
	SellPriceCost       float64 `json:"sell_price_cost"`
	SellProfitLossratio string  `json:"sell_profit_lossratio"`
	SellProfitReal      float64 `json:"sell_profit_real"`
	Symbol              string  `json:"symbol"`
	LevelRate           int     `json:"lever_rate"`
}

// FutureIndex represents the future index in rest resp.
type FutureIndex struct {
	Index float64 `json:"future_index"`
}

// EtcFutureInfoResp represents future current state
type EtcFutureInfoResp struct {
	Info   pair `json:"info"`
	Result bool `json:"result"`
}

type pair struct {
	Etc pairDetails `json:"etc"`
}

type pairDetails struct {
	Balance   float64    `json:"balance"`
	Contracts []contract `json:"contracts"`
	Rights    float64    `json:"rights"`
}

type contract struct {
	Available    float64 `json:"available"`
	Balance      float64 `json:"balance"`
	Bond         float64 `json:"bond"`
	ContractID   int64   `json:"contract_id"`
	ContractType string  `json:"contract_type"`
	Freeze       float64 `json:"freeze"`
	Profit       float64 `json:"profit"`
	UpProfit     float64 `json:"upprofit"`
}

// NewPair creates a new ok rest handler.
func NewPair(symbol, contractType, apiKey, apiSign string) *Pair {
	return &Pair{
		symbol:       symbol,
		contractType: contractType,
		apiKey:       apiKey,
		apiSign:      apiSign,
		c:            &http.Client{Timeout: 20 * time.Second},
	}
}

//GetFutureTicker receive the latest OKEX contract data
func (p *Pair) GetFutureTicker() *FutureTicker {
	req, err := http.NewRequest("GET", okEp+"/future_ticker.do", nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("symbol", p.symbol)
	q.Add("contract_type", p.contractType)
	req.URL.RawQuery = q.Encode()
	resp, err := p.c.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var ft FutureTicker
	err = json.Unmarshal(body, &ft)
	if err != nil {
		log.Fatal(err)
	}
	return &ft
}

// func (p *Pair) GetFutureDepth(size int, merge int) {
// 	req, err := http.NewRequest("GET", okEp+"/future_depth.do?", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	q := req.URL.Query()
// 	q.Add("symbol", p.symbol)
// 	q.Add("contract_type", p.contractType)
// 	req.URL.RawQuery = q.Encode()
// 	resp, err := p.c.Do(req)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	var ft FutureTicker
// 	err = json.Unmarshal(body, &ft)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return &ft
// }

// GetFutureTradeHistory returns recent future trade history.
func (p *Pair) GetFutureTradeHistory() {
	req, err := http.NewRequest("GET", okEp+"/future_trades.do", nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("symbol", p.symbol)
	q.Add("contract_type", p.contractType)
	req.URL.RawQuery = q.Encode()

	resp, err := p.c.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(body))
}

// GetFutureIndex returns the future index of chosen pair.
func (p *Pair) GetFutureIndex() {
	req, err := http.NewRequest("GET", okEp+"/future_index.do", nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("symbol", p.symbol)
	req.URL.RawQuery = q.Encode()

	resp, err := p.c.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var fi FutureIndex
	err = json.Unmarshal(body, &fi)
	if err != nil {
		log.Error(err)
	}

}

// GetFutureKlineData returns the Kline data
func (p *Pair) GetFutureKlineData(klineType string) []*Kline {
	req, err := http.NewRequest("GET", okEp+"/future_kline.do", nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("symbol", p.symbol)
	q.Add("type", klineType)
	q.Add("contract_type", "this_week")
	q.Add("size", "2000")
	req.URL.RawQuery = q.Encode()
	resp, err := p.c.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	data := strings.Replace(string(body), "[", "", -1)
	data = data[:len(data)-1]
	data += ","
	bufArray := strings.Split(data, "],")
	var klines []*Kline
	for _, v := range bufArray {
		if v == "" {
			continue
		}
		buf := strings.Split(v, ",")
		var k Kline
		k.TimeStamp, err = strconv.ParseInt(buf[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		k.Open, err = strconv.ParseFloat(buf[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		k.Highest, err = strconv.ParseFloat(buf[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		k.Lowest, err = strconv.ParseFloat(buf[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		k.Close, err = strconv.ParseFloat(buf[4], 64)
		if err != nil {
			log.Fatal(err)
		}
		k.UnitAmount, err = strconv.Atoi(buf[5])
		if err != nil {
			log.Fatal(err)
		}
		k.Amount, err = strconv.ParseFloat(buf[6], 64)
		if err != nil {
			log.Fatal(err)
		}
		klines = append(klines, &k)
	}
	return klines
}

// GetFuturePos4Fix returns current user future info .
func (p *Pair) GetFuturePos4Fix() (*FuturePosResp, error) {
	req, err := http.NewRequest("POST", okEp+"/future_position_4fix", nil)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	params["api_key"] = apiKey
	params["contract_type"] = p.contractType
	params["symbol"] = p.symbol
	params["type"] = "1"

	var query string
	q := req.URL.Query()

	for k, v := range params {
		query += k + "=" + v + "&"
		q.Add(k, v)
	}

	query += "secret_key" + "=" + apiSceret
	hasher := md5.New()
	hasher.Write([]byte(query))
	signature := hex.EncodeToString(hasher.Sum(nil))
	q.Add("sign", strings.ToUpper(signature))

	req.URL.RawQuery = q.Encode()
	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ftr FuturePosResp
	err = json.Unmarshal(body, &ftr)
	if err != nil {
		return nil, err
	}
	return &ftr, nil
}

// TradeResp returns trade resp
type TradeResp struct {
	Result  bool  `json:"result"`
	OrderID int64 `json:"order_id"`
}

// FutureTrade trades a future
func (p *Pair) FutureTrade(price string, amount string, tt TradeType, matchPrice bool) *TradeResp {
	req, err := http.NewRequest("POST", okEp+"/future_trade.do", nil)
	if err != nil {
		log.Fatal(err)
	}
	params := make(map[string]string)

	params["amount"] = amount
	params["api_key"] = apiKey
	params["contract_type"] = p.contractType
	if matchPrice {
		params["match_price"] = "1"
	} else {
		params["match_price"] = "0"

	}
	params["price"] = price
	params["symbol"] = p.symbol
	params["type"] = tt

	// var keys []string
	// for k := range params {
	// 	keys = append(keys, k)
	// }
	// sort.Strings(keys)
	var query string
	q := req.URL.Query()

	for k, v := range params {
		query += k + "=" + v + "&"
		q.Add(k, v)
	}

	query += "secret_key" + "=" + apiSceret
	hasher := md5.New()
	hasher.Write([]byte(query))
	signature := hex.EncodeToString(hasher.Sum(nil))
	q.Add("sign", strings.ToUpper(signature))

	req.URL.RawQuery = q.Encode()
	resp, err := p.c.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var trp TradeResp
	if err := json.Unmarshal(body, &trp); err != nil {
		log.Error(err)
	}
	return &trp
}

// GetFutureUserInfo4Fix returns current future user info.
func (p *Pair) GetFutureUserInfo4Fix() *EtcFutureInfoResp {
	req, err := http.NewRequest("POST", okEp+"/future_userinfo_4fix.do", nil)
	if err != nil {
		log.Fatal(err)
	}
	params := make(map[string]string)
	params["api_key"] = apiKey

	var query string
	q := req.URL.Query()
	for k, v := range params {
		query += k + "=" + v + "&"
		q.Add(k, v)
	}

	query += "secret_key" + "=" + apiSceret
	hasher := md5.New()
	hasher.Write([]byte(query))
	signature := hex.EncodeToString(hasher.Sum(nil))
	q.Add("sign", strings.ToUpper(signature))
	req.URL.RawQuery = q.Encode()
	resp, err := p.c.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var etfInfo EtcFutureInfoResp
	err = json.Unmarshal(body, &etfInfo)
	if err != nil {
		log.Fatal(err)
	}
	return &etfInfo
}

// func GetFuturePos() {
// 	req, err := http.NewRequest("POST", okEp+"/future_userinfo.do", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	q := req.URL.Query()
// 	q.Add
// 	q.Add("api_key", apiKey)
// 	q.Add("secret_key", apiSceret)
// }

// using json decoder
// t, err := dec.Token()
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("%T: %v\n", t, t)

// // while the array contains values
// for dec.More() {
// 	var m Kline
// 	// decode an array value (Message)

// 	err := dec.Decode(&m)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("%v: %v\n", m.TimeStamp, m.Open)
// }

// // read closing bracket
// t, err = dec.Token()
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("%T: %v\n", t, t)
// // err = json.Unmarshal(body, &klines)
// // if err != nil {
// // 	log.Error(err)
// // }
