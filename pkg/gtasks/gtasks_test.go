package gtasks_test

import (
	"testing"

	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks"
	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks/credentials"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAll(t *testing.T) {
	if testing.Short() {
		t.Skip("Long running test TestDeleteAll skipped.")
	}
	client, err := gtasks.NewClient(credentials.DefaultPath)
	assert.NoError(t, err)
	err = client.DeleteAll()
	assert.NoError(t, err)
}
