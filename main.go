package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var BuyRequestChan = make(chan string, 10000000)


var cookie = ""
var httpsUrlToken = ""
var httpsUrlBlance = ""
var httpsUrlBuy = ""

type HttpClient struct {
	Cookie string
}

type UsdtTemplate struct {
	ID              string `json:"id"`
	Amount          string `json:"amount"`
	PaymentCurrency string `json:"payment_currency"`
	Token           string `json:"token"`
}
var buyCount=0
const (
	typek  = "usdt"

	format = "2006-01-02 15:04:05.000"
)

func GetConfig() map[string]string {
	bytes, _ := ioutil.ReadFile("./cookie.json")
	var kv = map[string]string{}
	json.Unmarshal(bytes, &kv)
	return kv
}

func UrlConfig() map[string]string {
	bytes, _ := ioutil.ReadFile("./url.json")
	var kv = map[string]string{}
	json.Unmarshal(bytes, &kv)
	return kv
}

func TestBlance() int64 {
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
		var kvMap = map[string]map[string]string{}
		json.Unmarshal(bodyBytes, &kvMap)
		balance, _ := strconv.ParseFloat(kvMap["data"]["available"], 32/64)
		if int64(balance) <= 5000 {
			logs.Info("Count %d", int64(balance)*10)
			return int64(balance) * 10
			/*
			 抢购数额加一个逻辑，余额高于5000刀，数额就是50000，余额低于5000刀，数额等于余额×10，保留到整数，用int类型，以免位数太多数值溢出
			*/
		}
		logs.Info("USDT Balance %f", balance)
	}
	return 50000
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
	tokenChan := make(chan string)
	t, _ := strconv.ParseInt(timeOut, 10, 64)
	ticker := time.NewTicker(time.Duration(t) * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			go GetToken(tokenChan, cookie)

			go BuyRequest(tokenChan, cookie, data)
		}
	}()
	time.Sleep(30 * time.Minute)
	ticker.Stop()
}

func main() {
	err := initLog()
	if err != nil {
		return
	}
	logs.Info("Start Successfully")
	// ./ssynflood --time "1000" --name "zhangjianxin"
	var name, timeout string
	flag.StringVar(&name, "name", "zhangjianxin", "name")
	flag.StringVar(&timeout, "time", "1000", "time")

	flag.Parse()
	logs.Info("Timer:", timeout)
	logs.Info("用户名:", name)
	logs.Info("使用资金类型:", typek)
	logs.Info("启动时间:", time.Now().Format(format))


	var config = GetConfig()
	var url = UrlConfig()
	// 赋值
	httpsUrlToken = url["token"]
	httpsUrlBlance = url["balance"]
	httpsUrlBuy = url["buy"]

	cookie = config[name]

	amount := TestBlance()
	ids := strings.Split(httpsUrlBuy, "/")
	logs.Info("Select ID ", ids[len(ids)-2])
	data := &UsdtTemplate{
		ID:              ids[len(ids)-2],
		PaymentCurrency: typek,
		Amount:          strconv.FormatInt(amount, 10),
	}
	logs.Info("可买入数量:", amount)
	loopWorker(data, timeout)
}

func GetToken(tokenChan chan string, cookie string) {
	logs.Info("[Get]Token 可用数量： %d", len(tokenChan))
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

		logs.Info("TokenResponse ", string(bodyBytes))
		logs.Info("TokenResponseConsumption %s [GetToken Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))

		if reData != nil && reData["data"] != "" {
			tokenChan <- reData["data"]
			logs.Info("TokenResponseSuccess ", string(bodyBytes))
			logs.Info("TokenResponseSuccessConsumption %s [GetToken Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
		}
	}
}

func BuyRequest(tokenChan chan string, cookie string, data *UsdtTemplate) {
	logs.Info("[Buy]Token 可用数量： %d", len(tokenChan))
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
	logs.Info("[TokenSince] %s ", time.Since(e).String())
	b, _ := json.Marshal(data)
	// 发送数据体
	logs.Info("[BuySendData] %s", string(b))
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

		logs.Info("BuyResponse %s", string(bodyBytes))
		logs.Info("BuyResponseConsumption %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
		var kvMapBuy =map[string]string{}
		_ = json.Unmarshal(bodyBytes, &kvMapBuy)
		if !strings.Contains(string(bodyBytes),"too_many_request"){
			buyCount++
			logs.Info("BuyResponseSuccess  %s", string(bodyBytes))
			logs.Info("BuyResponseSuccessConsumption %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
		}
		BuyRequestChan <- string(bodyBytes)
	}

	logs.Info("BuySuccessCount %d",buyCount)
}
