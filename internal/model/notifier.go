package model

type Notifier interface {
	SendChangeNotification(data VersionChangeData) error
}
