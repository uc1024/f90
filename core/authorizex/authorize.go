package authorizex

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthorizeOtpions struct {
	secret         string
	ExpireDuration time.Duration
}

type Authorize struct {
	options AuthorizeOtpions
	parser  *jwt.Parser
}

func NewAuthorize(secret string, fun ...func(*AuthorizeOtpions)) *Authorize {
	options := AuthorizeOtpions{}
	options.ExpireDuration = time.Hour * 24 * 7
	options.secret = secret
	for _, f := range fun {
		f(&options)
	}
	return &Authorize{
		parser: jwt.NewParser(jwt.WithJSONNumber()),
	}
}

func (auth *Authorize) ParserCtxToken(ctx context.Context, ex Extractor) (resuls *jwt.Token, err error) {
	token, err := ex.Extract(ctx)
	_ = token
	if err != nil {
		return
	}
	resuls, err = auth.parser.Parse(token, auth.hs256Keyfunc())
	return
}

func (auth *Authorize) ParserRequestToken(r *http.Request, ex Extractor) (resuls *jwt.Token, err error) {
	token, err := ex.ExtractRequest(r)
	_ = token
	if err != nil {
		return
	}
	resuls, err = auth.parser.Parse(token, auth.hs256Keyfunc())

	return
}

func (auth *Authorize) ParserStringToken(token string) (resuls *jwt.Token, err error) {
	resuls, err = auth.parser.Parse(token, auth.hs256Keyfunc())
	return
}

func (auth *Authorize) hs256Keyfunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.options.secret), nil
	}
}

func (auth *Authorize) GenerateToken(claims jwt.MapClaims) (string, error) {
	_, ok := claims[jwtExpire]
	if !ok {
		claims[jwtExpire] = time.Now().Add(auth.options.ExpireDuration).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.options.secret))
}
