package model

import "time"

type ApplicationData struct {
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Image     string    `json:"image"`
	Timestamp time.Time `json:"timestamp"`
}
