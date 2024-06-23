package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

var app App

// Test main
func TestMain(m *testing.M) {
	app = App{}
	app.Init("users.db")
	app.DB.Exec("DELETE FROM users")
	test := m.Run()
	os.Exit(test)
}

// Test REST API
func TestServeHTTPTableDriven(t *testing.T) {
	var tests = []struct {
		method      string
		url         string
		data        []byte
		wantCode    int
		wantMessage string
	}{
		{"POST", "/save", []byte(`{"id":"123"}`), 200, "user with id:123 saved"},                                         // Test saving user 123
		{"POST", "/save", []byte(`{"id":123}`), 400, "bad data format"},                                                  // Test saving user with bad data
		{"POST", "/save", []byte(`{"id":""}`), 400, "bad data format, please add id(required),name,email,date_of_birth"}, // Test saving user with no id
		{"GET", "/idonotexist", nil, 404, "user id:idonotexist not found"},                                               // Test user not found
		{"GET", "/123", nil, 200, "user id:123 found"},                                                                   // Test getting user 123
	}

	for _, tt := range tests {

		log.Println("Testing ", tt.method, tt.url, string(tt.data))

		request := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(tt.data))
		response := httptest.NewRecorder()
		app.Router.ServeHTTP(response, request)

		// Marhsal body
		var output output
		err := json.Unmarshal(response.Body.Bytes(), &output)
		if err != nil {
			t.Errorf("cannot unmarshal body (err:%s)", err.Error())
		}

		// Check response code
		if output.Response_code != tt.wantCode {
			t.Errorf("Expected response_code %d. Got %d", tt.wantCode, output.Response_code)
		}

		// Check message
		if output.Message != tt.wantMessage {
			t.Errorf("Expected message %s. Got %s", tt.wantMessage, output.Message)
		}
	}
}
