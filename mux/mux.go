package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/base64"
	"strings"

	"time"

	"github.com/gorilla/mux"
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
	params := mux.Vars(r)
	user := params["user"]
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

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		fmt.Printf(
			"%s\t%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// 认证
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != "user" && pair[1] != "password" {
			http.Error(w, "Not authorized", 401)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").
		Path("/parameters/path/{user:user.*}"). // 支持正则表达式URL
		Handler(logger(basicAuth(pathParameters), "pathParameters"))
	router.Methods("GET").
		Path("/parameters/query").
		Handler(logger(basicAuth(queryParameters), "queryParameters"))
	router.Methods("PUT").
		Path("/handling/request").
		Handler(logger(basicAuth(handlingRequest), "handlingRequest"))

	log.Fatal(http.ListenAndServeTLS(":12345", "config/server-cert.pem",
		"config/server-key.pem", router))
}
