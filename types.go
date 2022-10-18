package main

const (
	// PutAction - Put artifacts
	PutAction = "put"

	// DeleteAction - Delete artifacts
	DeleteAction = "delete"

	// GetAction - Get artifacts
	GetAction = "get"
)

type (
	// Action - Input params
	Action struct {
		Action    string
		Bucket    string
		S3Class   string
		Key       string
		Artifacts []string
	}
)
