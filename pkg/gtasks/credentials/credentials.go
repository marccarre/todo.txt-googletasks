package credentials

import (
	"encoding/json"
	"io/ioutil"
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

// NewFromJSONFile reads Google OAuth credentials from the provided file.
func NewFromJSONFile(filepath string) (*Credentials, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return NewFromJSON(file)
}

// NewFromJSON reads Google OAuth credentials from the provided JSON bytes.
func NewFromJSON(jsonBytes []byte) (*Credentials, error) {
	credentials := &credentialsFile{}
	if err := json.Unmarshal(jsonBytes, &credentials); err != nil {
		return nil, err
	}
	return &credentials.Credentials, nil
}
