package api_notifier

import (
	"github.com/TheDoctor028/version-guard-operator/internal/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSendChangeNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	apiNotifier := &ApiNotifier{
		ApiUrl: server.URL,
	}

	data := model.VersionChangeData{}

	err := apiNotifier.SendChangeNotification(data)
	if err != nil {
		t.Errorf("SendChangeNotification failed: %v", err)
	}
}

func TestNewApiNotifier(t *testing.T) {
	apiUrl := "http://test-api.com"
	err := os.Setenv("API_URL", apiUrl)
	if err != nil {
		t.Errorf("Failed to set API_URL environment variable: %v", err)
		return
	}
	defer os.Unsetenv("API_URL")

	notifier, err := NewApiNotifier()
	if err != nil {
		t.Errorf("NewApiNotifier failed: %v", err)
		return
	}

	if notifier.ApiUrl != apiUrl {
		t.Errorf("ApiUrl mismatch: expected %s, got %s", apiUrl, notifier.ApiUrl)
	}
}
