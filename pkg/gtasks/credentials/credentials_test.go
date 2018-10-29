package credentials_test

import (
	"testing"

	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks/credentials"
	"github.com/stretchr/testify/assert"
)

func TestNewFromJSONFile(t *testing.T) {
	credentials, err := credentials.NewFromJSONFile("test_client_id.json")
	assert.NoError(t, err)
	assert.Equal(t, "somerandomid.apps.googleusercontent.com", credentials.ClientID)
	assert.Equal(t, "somesecret", credentials.ClientSecret)
}

func TestNewFromJSON(t *testing.T) {
	credentials, err := credentials.NewFromJSON([]byte(`{
		"installed": {
		  "client_id": "somerandomid.apps.googleusercontent.com",
		  "project_id": "someprojectid",
		  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
		  "token_uri": "https://www.googleapis.com/oauth2/v3/token",
		  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		  "client_secret": "somesecret",
		  "redirect_uris": [
			"urn:ietf:wg:oauth:2.0:oob",
			"http://localhost"
		  ]
		}
	  }`))
	assert.NoError(t, err)
	assert.Equal(t, "somerandomid.apps.googleusercontent.com", credentials.ClientID)
	assert.Equal(t, "somesecret", credentials.ClientSecret)
}
