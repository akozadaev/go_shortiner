package credential

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicAuthHeaderValue(t *testing.T) {
	credentials := GetTestCredentials()
	token := credentials
	value := "Basic " + base64.StdEncoding.EncodeToString([]byte(token))
	assert.True(t, value == BasicAuthHeaderValue)
}
