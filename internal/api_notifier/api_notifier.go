package api_notifier

import (
	"github.com/TheDoctor028/version-guard-operator/internal/model"
)

// ApiNotifier is a notifier that sends notifications to an API
type ApiNotifier struct {
	ApiUrl string
}

func (a ApiNotifier) SendNotification(data model.VersionChangeData) error {
	panic("implement me")
}

func NewApiNotifier() error {
	panic("implement me")
}
