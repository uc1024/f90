package signaturex

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/spf13/cast"
	"github.com/uc1024/f90/core/timex"
	"golang.org/x/exp/slices"
)

const (
	APP_ID    = "AppId"
	UnixMilli = "UnixMilli"
	Sign      = "Sign"

	HMAC_SHA256_SECRET = "hmac-sha256-secret"
)

var (
	ErrorSgin = errors.New("sign error")
)

// hmac-sha256-secret
type SignatureHmacSha256Secret struct {
	secret string
}

func NewSignatureHmacSha256Secret(secret string) *SignatureHmacSha256Secret {
	return &SignatureHmacSha256Secret{
		secret: secret,
	}
}

func (sign SignatureHmacSha256Secret) ChckeHead(head *http.Header, inKeys []string) (result map[string]string, err error) {
	slices.Sort(inKeys)
	result = make(map[string]string)
	for _, v := range inKeys {
		result[v] = head.Get(v)
		if result[v] == "" {
			return nil, fmt.Errorf("head key %s is empty", v)
		}
	}
	return
}

func (sign SignatureHmacSha256Secret) ChckeContext(ctx context.Context, inKeys []string) (result map[string]string, err error) {
	slices.Sort(inKeys)
	result = make(map[string]string)
	for _, v := range inKeys {
		value, ok := ctx.Value(v).(string)
		if !ok {
			return nil, fmt.Errorf("context key %s is not string", v)
		}
		if value == "" {
			return nil, fmt.Errorf("context key %s is empty", v)
		}
		result[v] = value
	}
	return result, nil
}

func (sign SignatureHmacSha256Secret) ChckeMap(m map[string]string, inKeys []string) (result map[string]string, err error) {
	slices.Sort(inKeys)
	result = make(map[string]string)
	for _, v := range inKeys {
		if m[v] == "" {
			return nil, fmt.Errorf("map key %s is empty", v)
		}
		result[v] = m[v]
	}
	return result, nil
}

// 根据请求生成签名
func (sign SignatureHmacSha256Secret) SigntureRequest(request *http.Request) (signStr string, err error) {

	ux := time.Now().UnixMilli()
	request.Header.Set(UnixMilli, cast.ToString(ux))

	// head info
	inKeys := []string{
		APP_ID,
		UnixMilli,
	}

	head, err := sign.ChckeHead(&request.Header, inKeys)
	if err != nil {
		return "", err
	}

	var body string

	// * 如果请求存在body
	if request.Body != nil &&
		request.Method != http.MethodGet &&
		request.Method != http.MethodDelete {
		origBody, err := io.ReadAll(request.Body)
		if err != nil {
			return "", err
		}
		request.Body.Close()
		request.Body = io.NopCloser(bytes.NewBuffer(origBody))
		if len(origBody) > 0 {
			body = *(*string)(unsafe.Pointer(&origBody))
		}
	}

	// * 签名验证
	h := strings.Builder{}
	h.Write([]byte(head[APP_ID]))
	h.Write([]byte(head[UnixMilli]))
	h.Write([]byte(strings.ToUpper(request.Method)))
	h.Write([]byte(body))
	h.Write([]byte(request.URL.RequestURI()))
	signStr = sign.Signature([]byte(h.String()))
	request.Header.Set(Sign, signStr)
	return
}

// 签名
func (sign SignatureHmacSha256Secret) Signature(buf []byte) string {
	h := hmac.New(sha256.New, []byte(sign.secret))
	h.Write(buf)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 校验请求签名验证
func (sign SignatureHmacSha256Secret) SigntureChckeRequest(request *http.Request) (err error) {
	// head info
	inKeys := []string{
		APP_ID,
		UnixMilli,
		Sign,
	}
	head, err := sign.ChckeHead(&request.Header, inKeys)
	if err != nil {
		return err
	}
	_ = head

	// * 超时
	err = sign.isTimeOut(cast.ToInt64(head[UnixMilli]))
	if err != nil {
		return err
	}

	var body string

	// * 如果请求存在body
	if request.Body != nil &&
		request.Method != http.MethodGet &&
		request.Method != http.MethodDelete {
		origBody, err := io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		request.Body.Close()
		request.Body = io.NopCloser(bytes.NewBuffer(origBody))
		if len(origBody) > 0 {
			body = *(*string)(unsafe.Pointer(&origBody))
		}
	}

	// * 签名验证
	h := strings.Builder{}
	h.Write([]byte(head[APP_ID]))
	h.Write([]byte(head[UnixMilli]))
	h.Write([]byte(strings.ToUpper(request.Method)))
	h.Write([]byte(body))
	h.Write([]byte(request.URL.RequestURI()))
	signature := sign.Signature([]byte(h.String()))
	if signature != head[Sign] {
		return ErrorSgin
	}
	return nil
}

func (sign SignatureHmacSha256Secret) isTimeOut(unix int64) error {
	return timex.IsCurrentTimeWithinInterval(unix, time.Minute*5)
}
