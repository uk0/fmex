package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/codegangsta/cli"
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

var buyCount = 0
var tokenCount = 0

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
	config["filename"] = "./femx.log"
	config["level"] = 6;
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println(" json.Marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

func loopWorker(data *UsdtTemplate, tokenT string, BuyT string) {

	tokenChan := make(chan string)
	T1, _ := strconv.ParseInt(tokenT, 10, 64)
	T2, _ := strconv.ParseInt(BuyT, 10, 64)

	t1 := time.NewTicker(time.Duration(T1) * time.Millisecond)

	t2 := time.NewTicker(time.Duration(T2) * time.Millisecond)

	go func() {
		for t := range t1.C {
			fmt.Println("GetToken at", t)
			go GetToken(tokenChan, cookie)
		}
	}()

	go func() {
		for t := range t2.C {
			fmt.Println("BuyRequest at", t)
			go BuyRequest(tokenChan, cookie, data)
		}
	}()
	time.Sleep(30 * time.Minute)
	t1.Stop()
	t2.Stop()
}

func GetBlance(name string, cookie string, httpsUrlBlance string) {
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
		logs.Info("UserName %s USDT Balance %s", name, string(bodyBytes))
	}
}

func NewBalanceCommand() cli.Command {
	return cli.Command{
		Name:  "balance",
		Usage: "get all user balance",
		Action: func(c *cli.Context) error {
			var kvs = GetConfig()
			for name, v := range kvs{
				GetBlance(name,v,UrlConfig()["balance"])
			}
			return nil
		},
	}
}

func NewServerCommand() cli.Command {
	return cli.Command{
		Name:  "testserver",
		Usage: "start test server",
		Action: func(c *cli.Context) error {
			StartTestServer(c.String("port"))
			return nil
		}, Flags: []cli.Flag{
			cli.StringFlag{Name: "port", Usage: "server is running on port"},          // name
		},
	}
}

func NewFmexCommand() cli.Command {
	return cli.Command{
		Name:  "fmex",
		Usage: "start fmex services",
		Action: func(c *cli.Context) error {


			NewFemx(c.String("name"),c.String("token_time"),c.String("buy_time"))
			return cli.ShowAppHelp(c)
		}, Flags: []cli.Flag{
			cli.StringFlag{Name: "name", Usage: "user name"},          // name
			cli.StringFlag{Name: "token_time", Usage: "token timer"},   // token
			cli.StringFlag{Name: "buy_time", Usage: "buy timer"}, // buy
		},
	}
}

func NewFemx(name string,tokenT string,BuyT string){
	err := initLog()
	if err != nil {
		return
	}

	if (name=="" || len(name)==0) && (tokenT=="" || len(tokenT)==0) && (BuyT=="" || len(BuyT)==0){
		logs.Error(" params is null")
		return
	}
	logs.Info("Start Successfully")

	logs.Info("Token Timer:", tokenT)
	logs.Info("Buy Timer:", BuyT)
	logs.Info("用户名:", name)
	logs.Info("使用资金类型:", typek)
	logs.Info("启动时间:", time.Now().Format(format))

	var config = GetConfig()
	var url = UrlConfig()
	// 赋值
	httpsUrlToken = url["token"]
	httpsUrlBlance = url["balance"]
	httpsUrlBuy = url["buy"]
	// 赋值
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
	loopWorker(data, tokenT, BuyT)
}

func GetToken(tokenChan chan string, cookie string) {
	logs.Info("Token 可用数量： %d", len(tokenChan))
	logs.Info("Token Start %s ", time.Now().Format(format))
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
	logs.Info("[BuyTokenSince] %s ", time.Since(e).String())
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
		var kvMapBuy = map[string]string{}
		_ = json.Unmarshal(bodyBytes, &kvMapBuy)

		// 请求过快 Token错误
		if strings.Contains(string(bodyBytes), "lightning_deal_token_access_limit") {
			tokenCount++
			logs.Info("BuyResponseFail  %s", string(bodyBytes))
			logs.Info("BuyResponseFailConsumption %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
			return
		}
		// 请求过快
		if strings.Contains(string(bodyBytes), "too_many_request") {
			tokenCount++
			logs.Info("BuyResponseFail  %s", string(bodyBytes))
			logs.Info("BuyResponseFailConsumption %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
			return
		}
		// 活动没有开始
		if strings.Contains(string(bodyBytes), "lightning_deal_not_valid_range") {
			tokenCount++
			logs.Info("BuyResponseFail  %s", string(bodyBytes))
			logs.Info("BuyResponseFailConsumption %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
			return
		}
		// 活动结束
		if strings.Contains(string(bodyBytes), "lightning_deal_finished") {
			tokenCount++
			logs.Info("BuyResponseFail  %s", string(bodyBytes))
			logs.Info("BuyResponseFailConsumption %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
			return
		}

		buyCount++
		logs.Info("BuyResponseSuccess  %s", string(bodyBytes))
		logs.Info("BuyResponseSuccessConsumption %s [BuyRequest Start] %s End %s ", time.Since(e).String(), e.Format(format), time.Now().Format(format))
		BuyRequestChan <- string(bodyBytes)
	}

	logs.Info("TokenSuccessBuyFailCount %d", tokenCount)
	logs.Info("BuySuccessCount %d", buyCount)

}
