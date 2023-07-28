package authorizex

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	defaultExpireDuration = time.Hour * 24 * 7
	defaultJwtIssuer      = "authorizex"
)

const (
	// Constant "aud" 代表受众（audience）声明。它表示JWT的预期受众。
	jwtAudience = "aud"
	// Constant "exp" 代表过期时间（expiration time）声明。它表示JWT的过期时间。
	jwtExpire = "exp"
	// Constant "jti" 代表JWT ID（JWT ID）声明。它表示JWT的唯一标识符。
	jwtId = "jti"
	// Constant "iat" 代表签发时间（issued at）声明。它表示JWT的签发时间。
	jwtIssueAt = "iat"
	// Constant "iss" 代表发行者（issuer）声明。它表示JWT的发行实体。
	jwtIssuer = "iss"
	// Constant "nbf" 代表不可用之前时间（not before）声明。它表示在此时间之前JWT不应被接受。
	jwtNotBefore = "nbf"
	// Constant "sub" 代表主题（subject）声明。它表示JWT所涉及的实体，即JWT所关注的主题。
	jwtSubject = "sub"
	// Constant 当没有提供具体细节时使用的默认原因。
	noDetailReason = "no detail reason"
)

var defaultMapClaimsKey = []string{
	"aud",
	"exp",
	"jti",
	"iat",
	"iss",
	"nbf",
	"sub",
	"no detail reason",
}

func DefaultClaimsKey() []string {
	return defaultMapClaimsKey
}

type SetClaimsOption func(jwt.MapClaims)

func NewAuthorizeClaims(opts ...SetClaimsOption) jwt.MapClaims {
	ac := jwt.MapClaims{
			jwtAudience:  []string{},
			jwtId:        "",
			jwtIssuer:    defaultJwtIssuer,
			jwtExpire:    time.Now().Add(defaultExpireDuration).Unix(),
			jwtIssueAt:   time.Now().Unix(),
			jwtNotBefore: time.Now().Unix(),
		}
	
	for _, opt := range opts {
		opt(ac)
	}
	return ac
}


// SetExp 设置过期时间
func SetExp(t time.Time) SetClaimsOption {
	return func(ac jwt.MapClaims) {
		ac[jwtExpire] = t.Unix()
	}
}

/*
	SetAUT 设置受众 在实际应用中，如果要设置"Audience"声明的默认值，
	可以根据应用的需求和设计选择适当的值。例如，如果JWT用于身份验证服务，
	可能将默认受众设置为该服务的标识符或URL。如果JWT用于特定的用户或角色，
	可以将默认受众设置为该用户或角色的标识符。
	默认值应根据具体场景和需求进行定义，并在代码中相应地设置。
*/
func SetAud(aud []string) SetClaimsOption {
	return func(ac jwt.MapClaims) {
		ac[jwtAudience] = aud
	}
}

/*
	SetJti 设置JWT ID
	JWT ID是一个字符串值，用于唯一标识JWT。它的作用是确保每个JWT都具有唯一的标识符，以防止重放攻击或重复使用JWT。通常，JWT ID的生成应遵循一定的规则和算法，以确保其全局唯一性。
*/
func SetJti(jti string) SetClaimsOption {
	return func(ac jwt.MapClaims) {
		ac[jwtId] = jti
	}
}

// SetIat 设置签发时间
func SetIat(t time.Time) SetClaimsOption {
	return func(ac jwt.MapClaims) {
		ac[jwtIssueAt] = t.Unix()
	}
}

/*
	SetNbf 设置不可用之前时间
	"nbf"声明用于指定JWT的生效时间，即在此时间之前JWT不应被接受或使用。
	它表示JWT的开始生效时间，即在此时间之前，JWT被视为无效或不可用。
*/
func SetNbf(t time.Time) SetClaimsOption {
	return func(ac jwt.MapClaims) {
		ac[jwtNotBefore] = t.Unix()
	}
}

// SetIss 设置发行者
func SetIss(iss string) SetClaimsOption {
	return func(ac jwt.MapClaims) {
		ac[jwtIssuer] = iss
	}
}



