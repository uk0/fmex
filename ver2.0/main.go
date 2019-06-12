package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"time"
)

var outTime int64 = 200;

var (
	tokenClient *TokenClient
	buyClient   *BuyClient
)

var TokenThreadNum = 4;

type UsdtTemplate struct {
	ID              string `json:"id"`
	Amount          string `json:"amount"`
	PaymentCurrency string `json:"payment_currency"`
	Token           string `json:"token"`
}

const (
	httpsUrlToken  = "https://www.fcoin.pro/openapi/auth/v1/lightning_deals/A-oY7xm0zzPtE3g7wFDx6A/token"
	httpsUrlBlance = "https://exchange.fcoin.pro/openapi/v3/assets/wallet/balances"
	httpsUrlBuy    = "https://www.fcoin.pro/openapi/auth/v1/lightning_deals/A-oY7xm0zzPtE3g7wFDx6A/buy"
)

var ticker = time.NewTicker(time.Duration(outTime) * time.Microsecond) // --- A
type TokenMessage struct {
	url string
}

type BuyClient struct {

}

type TokenClient struct {
	TaskTokenChan chan *TokenMessage
}

var TokenResponseCh = make(chan *string, 10000000)

func NewHttpTokenr() (tkn *TokenClient, err error) {
	tkn = &TokenClient{
		TaskTokenChan: make(chan *TokenMessage, 1000),
	}
	if err != nil {
		fmt.Printf("Failed to create Connetcion: %s\n", err)
		os.Exit(1)
	}
	for i := 0; i < TokenThreadNum; i++ {
		go tkn.Token()
	}

	return
}

func NewHttpBuy() (buy *BuyClient, err error) {
	buy = &BuyClient{

	}
	if err != nil {
		fmt.Printf("Failed to create Connetcion: %s\n", err)
		os.Exit(1)
	}
	for i := 0; i < TokenThreadNum; i++ {
		go buy.Buy()
	}

	return
}

func InitToken() (err error) {
	tokenClient, err = NewHttpTokenr()
	return err
}


func InitBuy() (err error) {
	buyClient, err = NewHttpBuy()
	return err
}


func GetConfig() map[string]string {
	bytes, _ := ioutil.ReadFile("./cookie.json")
	var kv = map[string]string{}
	json.Unmarshal(bytes, &kv)
	return kv
}


func (k *TokenClient) Token() {
	logs.Info("[start get Token]")
	for v := range k.TaskTokenChan {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer func(){
			// 用完需要释放资源
			fasthttp.ReleaseResponse(resp)
			fasthttp.ReleaseRequest(req)
		}()
		req.Header.Set("Cookie", GetConfig()["zhangjianxin"])
		req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
		req.Header.SetContentType("application/json")
		req.Header.SetMethod("POST")
		req.SetRequestURI(v.url)
		if err := fasthttp.Do(req, resp); err != nil {
			fmt.Println("请求失败:", err.Error())
			return
		}
		b := resp.Body()
		fmt.Println("result:\r\n", string(b))
	}
}

func (k *TokenClient) AddTask(url string){
	k.TaskTokenChan <- &TokenMessage{url: url}
	return
}


func (k *BuyClient)  Buy() {
	for i := range TokenResponseCh {
		fmt.Println(i)
	}
}


func main() {
	_ = InitToken()
	_ = InitBuy()
	loopWorker()
}

func loopWorker() {
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			tokenClient.AddTask(httpsUrlToken)
		}
	}
}
