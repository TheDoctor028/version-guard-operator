package api_notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TheDoctor028/version-guard-operator/internal/model"
	"net/http"
	"os"
)

// ApiNotifier is a notifier that sends notifications to an API
type ApiNotifier struct {
	ApiUrl string
}

// SendChangeNotification sends a notification to the API
func (a ApiNotifier) SendChangeNotification(data model.VersionChangeData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Call the API endpoint with the JSON data
	resp, err := http.Post(a.ApiUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// NewApiNotifier creates a new ApiNotifier
func NewApiNotifier() (*ApiNotifier, error) {
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		return nil, fmt.Errorf("API_URL environment variable is not set")
	}

	return &ApiNotifier{
		ApiUrl: apiUrl,
	}, nil
}
