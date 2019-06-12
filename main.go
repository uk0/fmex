package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"time"
)

var tokenChan = make(chan string, 100)
var BuyRequestChan = make(chan string, 100)

var cookie = ""

type UsdtTemplate struct {
	ID              string `json:"id"`
	Amount          string `json:"amount"`
	PaymentCurrency string `json:"payment_currency"`
	Token           string `json:"token"`
}

const (
	//httpsUrlToken  = "https://www.fcoin.pro/openapi/v1/lightning_deals/DBF9uRg3WBPwBCLKs0HrOQ/token"
	httpsUrlToken  = "https://www.fcoin.pro/openapi/auth/v1/lightning_deals/A-oY7xm0zzPtE3g7wFDx6A/token"
	httpsUrlBlance = "https://exchange.fcoin.pro/openapi/v3/assets/wallet/balances"
	httpsUrlBuy    = "https://www.fcoin.pro/openapi/auth/v1/lightning_deals/A-oY7xm0zzPtE3g7wFDx6A/buy"
)

func GetConfig() map[string]string {
	bytes, _ := ioutil.ReadFile("./cookie.json")
	var kv = map[string]string{}
	json.Unmarshal(bytes, &kv)
	return kv
}

func TestBlance() {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(httpsUrlBlance)
	req.Header.Set("Cookie", cookie)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.SetMethod("GET")
	resp := fasthttp.AcquireResponse()
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	client := &fasthttp.Client{
		TLSConfig: cfg,
	}
	if err := client.Do(req, resp); err != nil {
		println("Error:", err.Error())
	} else {
		bodyBytes := resp.Body()
		available := gjson.GetBytes(bodyBytes, "data.balances.5.available")
		fmt.Println(fmt.Sprintf(" USDT Blance %s", available.String()))
	}
}

func main() {

	// ./ssynflood --name "zhangjianxin" --amount "50000.00000000" --type "usdt"  --id "DBF9uRg3WBPwBCLKs0HrOQ"
	var name, amount, typek, ids string
	flag.StringVar(&name, "name", "zhangjianxin", "name")
	flag.StringVar(&amount, "amount", "50000.00000000", "amount")
	flag.StringVar(&typek, "type", "usdt", "type")
	flag.StringVar(&ids, "id", "DBF9uRg3WBPwBCLKs0HrOQ", "id")

	flag.Parse()

	fmt.Println("用户名:", name)
	fmt.Println("使用资金类型:", typek)
	fmt.Println("买入货币ID:", ids)
	fmt.Println("买入数量:", amount)

	data := &UsdtTemplate{
		ID:              ids,
		PaymentCurrency: typek,
		Amount:          amount,
	}

	var config = GetConfig()
	cookie = config[name]

	TestBlance()

	var i1 string
	for {
		select {
		case i1 = <-tokenChan:
			if len(i1) != 0 {
				go BuyRequest(i1, cookie, data)
			}
		default:
			//fmt.Println(fmt.Sprintf("Get Token.... %d", len(i1)))
		}
		go GetToken(cookie)
	}
}

func GetToken(cookie string) {
	e := time.Now()
	timer := time.AfterFunc(time.Microsecond*1000, func() {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI(httpsUrlToken)
		req.Header.Set("Cookie", cookie)
		req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
		req.Header.SetMethod("POST")
		resp := fasthttp.AcquireResponse()
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		}
		client := &fasthttp.Client{
			TLSConfig: cfg,
		}
		if err := client.Do(req, resp); err != nil {
			println("Error:", err.Error())
		} else {
			bodyBytes := resp.Body()
			var reData = map[string]string{}
			json.Unmarshal(bodyBytes, &reData)
			fmt.Println(string(bodyBytes))
			tokenChan <- reData["data"]
		}
		fmt.Println("time ",time.Since(e))

	})
	defer timer.Stop()

}

func BuyRequest(token string, cookie string, data *UsdtTemplate) {
	fmt.Println(fmt.Sprintf("token = %s", token))

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(httpsUrlBuy)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.SetMethod("POST")
	// 构造发送数据
	data.Token = token
	b, _ := json.Marshal(data)

	fmt.Println(fmt.Sprintf("data = %s", b))
	// 发送数据体
	req.SetBodyString(string(b))

	resp := fasthttp.AcquireResponse()

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	client := &fasthttp.Client{
		TLSConfig: cfg,
	}
	if err := client.Do(req, resp); err != nil {
		println("Error:", err.Error())
	} else {
		bodyBytes := resp.Body()
		fmt.Println(fmt.Sprintf("response %s ", string(bodyBytes)))
		BuyRequestChan <- string(bodyBytes)
	}
}
