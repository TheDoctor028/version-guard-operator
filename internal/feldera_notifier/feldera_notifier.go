package feldera_notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TheDoctor028/version-guard-operator/internal/model"
	"net/http"
	"os"
	"strings"
)

// FelderaNotifier is a notifier that sends notifications to an API
type FelderaNotifier struct {
	ApiUrl       string
	PipelineUUID string
}

// SendChangeNotification sends a notification to the API
func (a FelderaNotifier) SendChangeNotification(data model.VersionChangeData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	payload := []byte(fmt.Sprintf("{\"insert\": %s}", string(jsonData)))

	url := fmt.Sprintf("%s/v0/pipelines/%s/ingress/%s?format=json",
		a.ApiUrl, a.PipelineUUID, strings.ToUpper(data.Kind+"ChangeData"))
	// Call the API endpoint with the JSON data
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
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

// NewApiNotifier creates a new FelderaNotifier
func NewApiNotifier() (*FelderaNotifier, error) {
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		return nil, fmt.Errorf("API_URL environment variable is not set")
	}

	pipelineUUID := os.Getenv("PIPELINE_UUID")
	if pipelineUUID == "" {
		return nil, fmt.Errorf("PIPELINE_UUID environment variable is not set")
	}

	return &FelderaNotifier{
		ApiUrl:       apiUrl,
		PipelineUUID: pipelineUUID,
	}, nil
}
