package gtasks

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	tasks "google.golang.org/api/tasks/v1"
)

// Client encapsulates the Google Tasks and OAuth clients, and exposes
// high-level operations on Google Tasks.
type Client struct {
	api *tasks.Service
}

// NewClient creates a new Google Tasks client from the provided credentials
// JSON file.
func NewClient(path string) (*Client, error) {
	path, err := evaluateSymlinks(path)
	if err != nil {
		return nil, err
	}
	oauthClient, err := newOAuthClientFromJSONFile(path)
	if err != nil {
		return nil, err
	}
	api, err := tasks.New(oauthClient)
	if err != nil {
		return nil, err
	}
	return &Client{api: api}, nil
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
