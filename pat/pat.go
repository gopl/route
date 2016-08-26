package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name"`
	Email string `json:"email" xml:"email" form:"email"`
}

type Response struct {
	ReturnCode string `json:"returnCode"`
	Result     string `json:"result"`
}

func pathParameters(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get(":user")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "user="+user)
}

func queryParameters(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	user := params["user"] // 返回的是切片
	io.WriteString(w, fmt.Sprintf("user=%v, params=%v", user, params))
}

func handlingRequest(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 80))
	r.Body.Close()
	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		data, _ := json.Marshal(Response{"Failed", err.Error()})
		w.Write([]byte(data))
		return
	}
	fmt.Println(u)
	data, _ := json.Marshal(Response{"OK", "Handle Request"})
	w.Write([]byte(data))
}

func main() {
	router := pat.New()
	router.Get("/parameters/path/:user", http.HandlerFunc(pathParameters))
	router.Get("/parameters/query", http.HandlerFunc(queryParameters))
	router.Put("/handling/request", http.HandlerFunc(handlingRequest))

	log.Fatal(http.ListenAndServeTLS(":12345", "config/server-cert.pem",
		"config/server-key.pem", router))
}
