package gtasks_test

import (
	"testing"

	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks"
	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks/credentials"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAll(t *testing.T) {
	checkPreconditions(t)
	client, err := gtasks.NewClientFromEnvVars()
	assert.NoError(t, err)
	err = client.DeleteAll()
	assert.NoError(t, err)
}

func checkPreconditions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped TestDeleteAll: short tests only should be run.")
	}
	if !credentials.EnvVarsSet() {
		t.Skipf("Skipped TestDeleteAll: environment variables %v and %v are missing.", credentials.ClientID, credentials.ClientSecret)
	}
}
