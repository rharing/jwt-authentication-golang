package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	models "jwt-authentication-golang/models"
	"net/http"
	"testing"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func TestFlow(t *testing.T) {
	asserts := assert.New(t)
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

	asserts.Nil(err)
	body, _ = io.ReadAll(res.Body)
	expected := `{"message":"pong"}`
	if expected != string(body[:]) {
		t.Error("bad response")
	}
	request, err = http.NewRequest("GET", "http://localhost:8080/api/movies/cities?use4Testing=1", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authorization", tokenResponse.Token)
	res, err = http.DefaultClient.Do(request)

	asserts.Nil(err)
	body, _ = io.ReadAll(res.Body)
	type jsonCities struct {
		Key    string
		Cities []models.City
	}
	var cities jsonCities
	json.Unmarshal(body, &cities)
	if len(cities.Cities) < 100 {
		t.Fatal("expected at least 100 cities")
	}
	request, err = http.NewRequest("GET", "http://localhost:8080/api/movies/haarlem?use4Testing=1", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authorization", tokenResponse.Token)
	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
	body, _ = ioutil.ReadAll(res.Body)
	type jsonPlays struct {
		Key   string
		Plays []models.PlayDTO
	}
	var plays jsonPlays
	json.Unmarshal(body, &plays)
	if len(plays.Plays) < 100 {
		t.Fatal("expected at least 100 plays")
	}
	playDTO := plays.Plays[0]
	assert.NotNil(t, playDTO.Title)
	assert.NotNil(t, playDTO.MovieId)
	assert.NotNil(t, playDTO.Tickethref)
	assert.NotNil(t, playDTO.Moviehref)
	assert.NotNil(t, playDTO.Start)
	request, err = http.NewRequest("GET", "http://localhost:8080/api/movie/oppenheimer?use4Testing=1", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authorization", tokenResponse.Token)
	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

}
