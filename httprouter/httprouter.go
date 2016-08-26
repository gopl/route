package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name"`
	Email string `json:"email" xml:"email" form:"email"`
}

type Response struct {
	ReturnCode string `json:"returnCode"`
	Result     string `json:"result"`
}

func pathParameters(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "user="+user)
}

func queryParameters(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")
	io.WriteString(w, fmt.Sprintf("user=%v, params=%v", user, ps))
}

func handlingRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	router := httprouter.New()
	router.GET("/parameters/path/:user", pathParameters)
	router.GET("/parameters/query", queryParameters)
	router.PUT("/handling/request", handlingRequest)

	log.Fatal(http.ListenAndServeTLS(":12345", "config/server-cert.pem",
		"config/server-key.pem", router))
}
