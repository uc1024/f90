package authorizex

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// head /dev/urandom | tr -dc A-Za-z0-9 | head -c 32 ; echo ''

// okiSxx10kH0vqiYxnbk51XhVRkMjHXXU

const mockSecret = "okiSxx10kH0vqiYxnbk51XhVRkMjHXXU"

func TestAuthorize(t *testing.T) {
	au := NewAuthorize(mockSecret)
	ac := NewAuthorizeClaims()
	ac["user_id"] = "1"
	tk1, err := au.GenerateToken(ac)
	assert.NoError(t, err)
	ptk1, err := au.ParserStringToken(tk1)
	assert.NoError(t, err)
	assert.Equal(t, ptk1.Valid, true)
	assert.Equal(t, ptk1.Claims.(jwt.MapClaims)["user_id"], "1")
}

func newMockCtx(key string, value string) context.Context {
	base_ctx := context.Background()
	newctx := context.WithValue(base_ctx, key, value)
	return newctx
}

func TestAuthorizeFromContext(t *testing.T) {
	au := NewAuthorize(mockSecret)
	ac := NewAuthorizeClaims(SetExp(time.Now().AddDate(0, 1, 0)))
	ac["user_id"] = "1"
	tk1, err := au.GenerateToken(ac)
	assert.NoError(t, err)
	mock_ctx := newMockCtx("Token", tk1)
	ptk1, err := au.ParserCtxToken(mock_ctx, TokenExtractor{})
	assert.NoError(t, err)
	assert.Equal(t, ptk1.Valid, true)
	t.Log(ptk1.Header)
	assert.Equal(t, ptk1.Claims.(jwt.MapClaims)["user_id"], "1")
}
