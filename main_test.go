package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/gin-gonic/gin.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const adminAccount = "admin"
const adminPassword = "ElPsyCongroo"

func init() {
	setupDB("sqlite3", "test.db")
	initAdmin(adminAccount, adminPassword)
}

func Authenticate(router *gin.Engine) string {
	b, _ := json.Marshal(map[string]string{
		"account":  adminAccount,
		"password": adminPassword,
	})
	url := "/api/authenticate"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	respJSON := make(map[string]string)
	b, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &respJSON)
	return respJSON["token"]
}

func Test_PostAuthenticate(t *testing.T) {
	url := "/api/authenticate"
	b, _ := json.Marshal(map[string]string{
		"account":  adminAccount,
		"password": adminPassword,
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	resp := httptest.NewRecorder()

	router := GetMainEngine("JWT_Secret")
	router.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Errorf("Response code should be 200, was: %d.", resp.Code)
		t.Errorf("Response.Body : %v", resp.Body)
	}

	// test wrong password
	b, _ = json.Marshal(map[string]string{
		"account":  adminAccount,
		"password": adminPassword + "_wrong",
	})
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(b))
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != 401 {
		t.Errorf("Response code should be 401, was: %d.", resp.Code)
	}
}

// func Test_getBallonStatus(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/board", nil)
// 	response := httptest.NewRecorder()

// 	router := GetMainEngine()
// 	router.ServeHTTP(response, request)
// 	if response.Code != 200 {
// 		t.Errorf("Response code should be 200, was: %d", response.Code)
// 	}
// }

// func Test_getBallonStatus(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/ballon", nil)
// 	response := httptest.NewRecorder()

// 	router := GetMainEngine()
// 	router.ServeHTTP(response, request)
// 	if response.Code != 200 {
// 		t.Errorf("Response code should be 200, was: %s", response.Code)
// 	}
// }
