package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", tokenHandler)
	mux.HandleFunc("/buy", buyHandler)
	http.ListenAndServe(":4000", limit(mux))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)

	w.Write([]byte("{\"status\":\"ok\",\"data\":\"af81ced082294a22b335214b26a8cfe5\"}"))
}


func buyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	w.Write([]byte("buy OK"))
}