package tests

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func FuzzCreateUserCheck(f *testing.F) {
	client := getTestClient()
	f.Fuzz(func(t *testing.T, name string, domain string) {
		email := name + "@" + domain + ".com"
		resp, err := client.createUser(name, email)

		if err != nil {
			switch {
			case errors.Is(err, ErrForbidden):
				fallthrough
			case errors.Is(err, ErrBadRequest):
				return
			default:
				panic(err.Error())
			}
		}
		id := resp.Data.ID
		assert.NoError(t, err)
		assert.Equal(t, name, resp.Data.Nickname)
		assert.Equal(t, email, resp.Data.Email)

		resp, err = client.getUser(id)

		assert.NoError(t, err)
		assert.Equal(t, id, resp.Data.ID)
		assert.Equal(t, name, resp.Data.Nickname)
		assert.Equal(t, email, resp.Data.Email)
	})
}
