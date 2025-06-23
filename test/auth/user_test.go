package auth

import (
	"encoding/json"
	"go_shurtiner/test/credential"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "v1/users", nil)
	assert.NoError(t, err)
	request.Header.Set("Authorization", credential.BasicAuthHeaderValue)

	r := httptest.NewRecorder()

	respBody, err := json.Marshal("")
	assert.NoError(t, err)

	assert.Equal(t, 200, r.Code)
	assert.NotEqual(t, respBody, r.Body.Bytes())
}
