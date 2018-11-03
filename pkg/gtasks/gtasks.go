package gtasks

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	tasks "google.golang.org/api/tasks/v1"
)

// NewClient creates a new Google Tasks client from the provided credentials JSON file.
func NewClient(path string) (*tasks.Service, error) {
	path, err := evaluateSymlinks(path)
	if err != nil {
		return nil, err
	}
	oauthClient, err := newOAuthClientFromJSONFile(path)
	if err != nil {
		return nil, err
	}
	return tasks.New(oauthClient)
}

func evaluateSymlinks(path string) (string, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
