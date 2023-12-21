package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jwt-authentication-golang/movies"
	"net/http"
	"testing"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func TestFlow(t *testing.T) {
	var jsonStr = []byte(`{"email":"mukesh@go.com","password":"123465789"}`)
	//fmt.Println(string(jsonStr[:]))
	request, err := http.NewRequest("POST", "http://localhost:8080/api/token", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
	body, _ := ioutil.ReadAll(res.Body)
	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		panic(err)
	}
	fmt.Printf("token %v", tokenResponse.Token)
	//GET http://{{host}}/api/secured/ping HTTP/1.1
	//content-type: application/json
	//authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im11a2VzaC5tdXJ1Z2FuIiwiZW1haWwiOiJtdWtlc2hAZ28uY29tIiwiZXhwIjoxNzAyMDI5OTI3fQ._e1_2PgKeVlbMq5Gv9stfcuGb9A5MJo0fw51fuRQeGM
	request, err = http.NewRequest("GET", "http://localhost:8080/api/secured/ping", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authorization", tokenResponse.Token)

	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
	body, _ = ioutil.ReadAll(res.Body)
	expected := `{"message":"pong"}`
	if expected != string(body[:]) {
		t.Error("bad response")
	}
	request, err = http.NewRequest("GET", "http://localhost:8080/api/movies/cities", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authorization", tokenResponse.Token)
	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
	body, _ = ioutil.ReadAll(res.Body)
	type jsonCities struct {
		Key    string
		Cities []movies.City
	}
	var cities jsonCities
	json.Unmarshal(body, &cities)
	if len(cities.cities) < 100 {
		t.Fatal("expected at least 100 cities")
	}
}
