package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/valyala/fasthttp"
	"io/ioutil"
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


func TestBlance(name string,cookie string,httpsUrlBlance string) {
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
		logs.Info("UserName %s USDT Balance %s", name,string(bodyBytes))
	}
}

func main()  {
	var kvs = GetConfig()
	for name, v := range kvs{
		TestBlance(name,v,UrlConfig()["balance"])
	}
}