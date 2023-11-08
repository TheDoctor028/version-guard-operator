package model

type Notifier interface {
	SendNotification(data VersionChangeData) error
}
