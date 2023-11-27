package signaturex

import (
	"errors"
)

const (
	HMAC_SHA256_SECRET = "hmac-sha256-secret"
)

var (
	ErrorSign = errors.New("sign error")
)

// hmac-sha256-secret
type SignatureHmacSha256Secret struct {
	name   string
	secret string
}

func NewSignatureHmacSha256Secret(secret string) *SignatureHmacSha256Secret {
	return &SignatureHmacSha256Secret{
		name:   HMAC_SHA256_SECRET,
		secret: secret,
	}
}

func (sign *SignatureHmacSha256Secret) Name() string {
	return sign.name
}