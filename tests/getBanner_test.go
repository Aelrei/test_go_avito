package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetBanner(t *testing.T) {
	expectedResponses := []string{`"{\"text\":\"some_text1\",\"title\":\"some_title1\",\"url\":\"some_url1\"}"`,
		`"{\"text\":\"some_text1\",\"title\":\"some_title1\",\"url\":\"some_url1\"}"`,
		`"{\"text\":\"some_text999\",\"title\":\"some_title999\",\"url\":\"some_url999\"}"`}

	urls := []string{"http://localhost:8085/user_banner?tag_id=1&feature_id=2&use_last_revision=true",
		"http://localhost:8085/user_banner?feature_id=2&tag_id=1",
		"http://localhost:8085/user_banner?feature_id=6&tag_id=999"}

	client := &http.Client{}

	for i, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("token", "user_token")

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
		err = json.Unmarshal(body, &data)
		if err != nil {
			t.Fatalf("Failed to unmarshal response JSON: %v", err)
		}

		jsonString, err := json.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to marshal response JSON: %v", err)
		}

		jsonBytes := append(jsonString, '\n')
		jsonBytes, _ = json.MarshalIndent(string(jsonString), "", " ")

		expectedResponse := expectedResponses[i]

		if string(jsonBytes) != expectedResponse {
			t.Fatalf("Response does not match expected JSON\nExpected: %s\nActual: %s", expectedResponse, string(jsonBytes))
		}
	}

	fmt.Println("Test passed successfully!")
}
