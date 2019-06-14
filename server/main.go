package main

import (
	"encoding/json"
	"fmt"
	"github.com/didip/tollbooth"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

const (
	request_limit = "{\"status\":\"error\",\"err_code\":\"too_many_request\",\"err_msg\":\"Your order is too fast. Please try again later.\",\"data\":\"null\"}"
)

var tempUUID = "";

func main() {
	mux := http.NewServeMux()
	mux.Handle("/token", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(2, nil).SetMessage(request_limit), tokenHandler))
	mux.Handle("/buy", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(2, nil).SetMessage(request_limit), buyHandler))
	_ = http.ListenAndServe(":4000", mux)
}

func GetUUID() string {
	// 创建
	u1 := uuid.NewV4()
	return u1.String()
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	var Uid = GetUUID();
	tempUUID = Uid

	_, _ = w.Write([]byte(fmt.Sprintf("{\"status\":\"ok\",\"data\":\"%s\"}", Uid)))
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(fmt.Sprintf("Buy Body %s", string(body)))

	var kvms = map[string]string{}
	_ = json.Unmarshal(body, &kvms)

	if kvms["token"] == tempUUID {
		_, _ = w.Write([]byte("{\"status\":\"error\",\"err_code\":\"lightning_deal_finished\",\"err_msg\":\"Lightning deal finished.\",\"data\":\"null\"}"))
	} else {
		_, _ = w.Write([]byte("{\"status\":\"error\",\"err_code\":\"token is error \",\"err_msg\":\"bad_params.\"}"))
	}
}
