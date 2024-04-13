package main

import (
	json2 "encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetBanner(t *testing.T) {
	expectedResponse := `"{\"text\":\"some_text\",\"title\":\"some_title\",\"url\":\"some_url\"}"`

	url := "http://localhost:8085/user_banner?tag_id=1&feature_id=1"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("token", "admin_token")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var data map[string]interface{}

	err = json2.Unmarshal(body, &data)
	jsonString, err := json2.Marshal(data)
	jsonBytes := append(jsonString, '\n')
	jsonBytes, _ = json2.MarshalIndent(string(jsonString), "", " ")
	if err != nil {
		fmt.Println("Ошибка при маршалинге в JSON:", err)
		return
	}

	if string(jsonBytes) != expectedResponse {
		t.Fatalf("Response does not match expected JSON\nExpected: %s\nActual: %s", expectedResponse, string(jsonBytes))
	}

	fmt.Println("Test passed successfully!")
}
