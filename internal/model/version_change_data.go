package model

import (
	"fmt"
	"time"
)

const (
	DeploymentKind  = "Deployment"
	ApplicationKind = "Application"
)

type VersionChangeData struct {
	Kind          string    `json:"kind"`
	Name          string    `json:"name"`
	Namespace     string    `json:"namespace"`
	Selector      string    `json:"selector"`
	ContainerName string    `json:"container_name"`
	Image         string    `json:"image"`
	Timestamp     time.Time `json:"timestamp"`
}

type VersionChangeDataMatcher struct {
	expected VersionChangeData
}

func VersionChangeDataEQ(x VersionChangeData) VersionChangeDataMatcher {
	return VersionChangeDataMatcher{expected: x}
}

func (v VersionChangeDataMatcher) Matches(x interface{}) bool {
	got := x.(VersionChangeData)

	return got.Kind == v.expected.Kind &&
		got.Name == v.expected.Name &&
		got.Namespace == v.expected.Namespace &&
		got.ContainerName == v.expected.ContainerName &&
		got.Selector == v.expected.Selector &&
		got.Image == v.expected.Image &&
		(got.Timestamp.Unix()-v.expected.Timestamp.Unix()) < 1000*5
	// if the timestamps are within 5 seconds of each other for test purposes, they are the same
}

func (v VersionChangeDataMatcher) String() string {
	return fmt.Sprintf("%s,", v.expected)
}

// ParseSelector converts a selector map to a string
func ParseSelector(selector map[string]string) string {
	var result string
	for key, val := range selector {
		result += fmt.Sprintf("%s=%s,", key, val)
	}
	return result[:len(result)-1]
}
