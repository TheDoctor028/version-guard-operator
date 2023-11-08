package model

import "time"

type VersionChangeData struct {
	Kind          string    `json:"kind"`
	Name          string    `json:"name"`
	Namespace     string    `json:"namespace"`
	ContainerName string    `json:"container_name"`
	Image         string    `json:"image"`
	Timestamp     time.Time `json:"timestamp"`
}
