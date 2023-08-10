package signaturex

import (
	"bytes"
	"net/http"
	"testing"

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
	_, err = sign.SigntureRequest(req)
	assert.NoError(t, err)
	
	// Call the SigntureChckeRequest method
	err = sign.SigntureChckeRequest(req)
	assert.NoError(t, err)
}
