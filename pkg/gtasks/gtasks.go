package gtasks

import (
	tasks "google.golang.org/api/tasks/v1"
)

// NewClient creates a new Google Tasks client from the provided credentials JSON file.
func NewClient(filepath string) (*tasks.Service, error) {
	oauthClient, err := newOAuthClientFromJSONFile(filepath)
	if err != nil {
		return nil, err
	}
	return tasks.New(oauthClient)
}
