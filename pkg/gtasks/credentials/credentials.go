package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	// DefaultPath is the default filepath for OAuth JSON credentials.
	DefaultPath = "~/.todo.txt-googletasks_credentials.json"
	// ClientID is the environment variable to read an OAuth client ID.
	ClientID = "CLIENT_ID"
	// ClientSecret is the environment variable to read an OAuth client ID.
	ClientSecret = "CLIENT_SECRET"
)

// Credentials represents Google OAuth credentials.
type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Wrapper data structure around Credentials.
type credentialsFile struct {
	Credentials Credentials `json:"installed"`
}

// NewFromJSONFile reads Google OAuth credentials from the provided file, and
// creates a new Credentials object with the values read.
func NewFromJSONFile(path string) (*Credentials, error) {
	path, err := evaluateSymlinks(path)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewFromJSON(file)
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

// NewFromJSON reads Google OAuth credentials from the provided JSON bytes,
// and creates a new Credentials object with the values read.
func NewFromJSON(jsonBytes []byte) (*Credentials, error) {
	credentials := &credentialsFile{}
	if err := json.Unmarshal(jsonBytes, &credentials); err != nil {
		return nil, err
	}
	return &credentials.Credentials, nil
}

// NewFromEnvVars reads the CLIENT_ID and CLIENT_SECRET environment variables,
// and creates a new Credentials object with the values read.
func NewFromEnvVars() (*Credentials, error) {
	clientID := os.Getenv(ClientID)
	if clientID == "" {
		return nil, fmt.Errorf("%v is required but missing", ClientID)
	}
	clientSecret := os.Getenv(ClientSecret)
	if clientSecret == "" {
		return nil, fmt.Errorf("%v is required but missing", ClientSecret)
	}
	return &Credentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

// EnvVarsSet checks if the CLIENT_ID and CLIENT_SECRET environment variables,
// and returns true if so, or false otherwise.
func EnvVarsSet() bool {
	return os.Getenv(ClientID) != "" && os.Getenv(ClientSecret) != ""
}
