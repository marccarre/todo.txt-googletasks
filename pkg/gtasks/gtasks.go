package gtasks

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
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

// DeleteAll deletes all tasks in all lists.
func (c Client) DeleteAll() error {
	log.Info("Fetching all lists")
	lists, err := c.api.Tasklists.List().Do()
	if err != nil {
		return err
	}
	for _, list := range lists.Items {
		hasNextPage := true
		for hasNextPage {
			log.WithField("id", list.Id).WithField("title", list.Title).Info("Fetching all tasks")
			tasks, err := c.api.Tasks.List(list.Id).Do()
			if err != nil {
				return err
			}
			log.WithField("count", len(tasks.Items)).Info("Deleting tasks")
			for _, task := range tasks.Items {
				if err := c.api.Tasks.Delete(list.Id, task.Id).Do(); err != nil {
					return err
				}
			}
			hasNextPage = tasks.NextPageToken != "" // continue until no next page.
		}
	}
	log.Info("Done: deletion of all tasks in all lists completed successfully.")
	return nil
}
