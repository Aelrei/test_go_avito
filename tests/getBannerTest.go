package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBanner(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем URL запроса и, если необходимо, возвращаем соответствующий ответ
		if r.URL.Path == "/banner" && r.URL.RawQuery == "tag_id=2&limit=100&offset=0" {
			// Возвращаемый JSON баннера
			bannerJSON := `{"id":1,"content":{"Title":"Test Banner","Text":"This is a test banner","URL":"http://example.com"},"feature_id":1,"active":true,"created_at":"2022-04-12T12:00:00Z","updated_at":"2022-04-12T12:00:00Z"}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(bannerJSON))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	// Создаем HTTP-клиент для отправки запросов к тестовому серверу
	client := ts.Client()

	// Отправляем запрос на получение баннера
	resp, err := client.Get(ts.URL + "/banner?tag_id=2&limit=100&offset=0")
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Проверяем, что полученный JSON соответствует ожидаемому баннеру
	expectedJSON := `{"id":1,"content":{"Title":"Test Banner","Text":"This is a test banner","URL":"http://example.com"},"feature_id":1,"active":true,"created_at":"2022-04-12T12:00:00Z","updated_at":"2022-04-12T12:00:00Z"}`
	if string(body) != expectedJSON {
		t.Errorf("expected JSON %s, got %s", expectedJSON, string(body))
	}
}
