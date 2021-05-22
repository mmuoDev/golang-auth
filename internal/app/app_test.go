package app_test

import (
	"fmt"
	"golang-auth/internal"
	"golang-auth/internal/app"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

//mongoDBProvider mocks mongo DB
func mongoDBProvider() *mongo.Database {
	return nil
}

func TestCreateUserReturns201(t *testing.T) {
	insertIntoDBInvoked := false

	mockUserInsert := func(o *app.OptionalArgs) {
		o.AddUser = func(u internal.User) error {
			insertIntoDBInvoked = true
			return nil
		}
	}

	//optional args
	opts := []app.Options{
		mockUserInsert,
	}

	ap := app.New(mongoDBProvider, opts...)
	serverURL, cleanUpServer := app.NewTestServer(ap.Handler())
	defer cleanUpServer()

	reqPayload, _ := os.Open(filepath.Join("testdata", "add_user_request.json"))
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", serverURL), reqPayload)

	client := &http.Client{}
	res, _ := client.Do(req)

	t.Run("Http Status Code is 201", func(t *testing.T) {
		assert.Equal(t, res.StatusCode, http.StatusCreated)
	})

	t.Run("Insert to DB invoked", func(t *testing.T) {
		assert.True(t, insertIntoDBInvoked)
	})
}
