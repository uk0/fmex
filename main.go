package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strconv"
	"time"
)

var BuyRequestChan = make(chan string, 10000000)

var cookie = ""

type HttpClient struct {
	Cookie string
}

type UsdtTemplate struct {
	ID              string `json:"id"`
	Amount          string `json:"amount"`
	PaymentCurrency string `json:"payment_currency"`
	Token           string `json:"token"`
}

const (
	format = "2006-01-02 15:04:05.000"
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
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		println("Error:", err.Error())
	} else {
		bodyBytes := resp.Body()
		available := gjson.GetBytes(bodyBytes, "data.balances.5.available")
		logs.Info("USDT Balance %s", available.String())
	}
}

func initLog() (err error) {
	//初始化日志库
	config := make(map[string]interface{})
	config["filename"] = "./ssynflood.log"
	config["level"] = 6;
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println(" json.Marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

func loopWorker(data *UsdtTemplate, timeOut string) {
	tokenChan := make(chan string, 1000000)
	t, _ := strconv.ParseInt(timeOut, 10, 64)
	ticker := time.NewTicker(time.Duration(t) * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			go GetToken(tokenChan, cookie)

			go BuyRequest(tokenChan, cookie, data)
		}
	}()
	time.Sleep(60 * time.Minute)
	ticker.Stop()
}

func main() {
	err := initLog()
	if err != nil {
		return
	}
	logs.Info("Start Successfully")
	// ./ssynflood --time "1000" --name "zhangjianxin" --amount "50000.00000000" --type "usdt"  --id "DBF9uRg3WBPwBCLKs0HrOQ"
	var name, amount, typek, ids, timeout string
	flag.StringVar(&name, "name", "zhangjianxin", "name")
	flag.StringVar(&amount, "amount", "50000.00000000", "amount")
	flag.StringVar(&typek, "type", "usdt", "type")
	flag.StringVar(&ids, "id", "DBF9uRg3WBPwBCLKs0HrOQ", "id")
	flag.StringVar(&timeout, "time", "1000", "time")

	flag.Parse()
	logs.Info("TimeOut:", timeout)
	logs.Info("用户名:", name)
	logs.Info("使用资金类型:", typek)
	logs.Info("买入货币ID:", ids)
	logs.Info("买入数量:", amount)
	logs.Info("启动时间:", time.Now().Format(format))

	data := &UsdtTemplate{
		ID:              ids,
		PaymentCurrency: typek,
		Amount:          amount,
	}

	var config = GetConfig()
	cookie = config[name]

	TestBlance()
	loopWorker(data, timeout)
}

func GetToken(tokenChan chan string, cookie string) {
	logs.Info("GetToken Start %s ", time.Now().Format(format))
	e := time.Now()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(httpsUrlToken)
	req.Header.Set("Cookie", cookie)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.SetMethod("POST")
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	client := &fasthttp.Client{
	}
	if err := client.Do(req, resp); err != nil {
		println("Error:", err.Error())
	} else {
		bodyBytes := resp.Body()
		var reData = map[string]string{}
		json.Unmarshal(bodyBytes, &reData)
		logs.Info("[Get]Token 可用数量： %d", len(tokenChan))
		logs.Info("GetToken Response ", string(bodyBytes))
		tokenChan <- reData["data"]
	}
	logs.Info("[Consumption] %s [GetToken Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))

}

func BuyRequest(tokenChan chan string, cookie string, data *UsdtTemplate) {
	logs.Info("BuyRequest Start %s ", time.Now().Format(format))
	e := time.Now()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(httpsUrlBuy)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.SetMethod("POST")
	// 构造发送数据
	data.Token = <-tokenChan // 直接拿到Token
	b, _ := json.Marshal(data)
	// 发送数据体
	req.SetBodyString(string(b))
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	client := &fasthttp.Client{
	}
	if err := client.Do(req, resp); err != nil {
		println("Error:", err.Error())
	} else {
		bodyBytes := resp.Body()
		//logs.Info("Buy Response %s", gjson.GetBytes(bodyBytes, "status").String())
		logs.Info("[Buy]Token 可用数量： %d", len(tokenChan))
		logs.Info("Buy Response %s", string(bodyBytes))
		BuyRequestChan <- string(bodyBytes)
	}
	logs.Info("[Consumption] %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
}
