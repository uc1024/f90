package signaturex

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestSignatureHmacSha256Secret_SigntureChckeRequest(t *testing.T) {
	// Create a new instance of SignatureHmacSha256Secret with a test secret
	sign := NewSignatureHmacSha256Secret("test-secret")

	// Create a new http.Request with test data
	req, err := http.NewRequest(http.MethodPost, "http://example.com/test", bytes.NewBufferString("test-body"))
	assert.NoError(t, err)

	// Set the required headers
	req.Header.Set("AppId", "test-app-id")
	req.Header.Set("UnixMilli", cast.ToString(time.Now().UnixMilli()))
	s, err := sign.SigntureRequest(req)
	assert.NoError(t, err)
	req.Header.Set("Sign", s)

	// Call the SigntureChckeRequest method
	err = sign.SigntureChckeRequest(req)
	assert.NoError(t, err)
}
