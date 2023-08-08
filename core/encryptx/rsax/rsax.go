package rsax

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/uc1024/f90/utilx"
)

type (
	rsaxSecurityoptions struct {
		pubStr string // * 公钥字符串
		priStr string // * 私钥字符串
	}

	rsaxSecurityOptions func(*rsaxSecurityoptions)

	RSAXSecurity struct {
		options rsaxSecurityoptions
		pubkey  *rsa.PublicKey  // * 公钥
		prikey  *rsa.PrivateKey // * 私钥
	}
)

func SetPublicString(val string) rsaxSecurityOptions {
	return func(o *rsaxSecurityoptions) {
		o.pubStr = val
	}
}

func SetPrivateString(val string) rsaxSecurityOptions {
	return func(o *rsaxSecurityoptions) {
		o.priStr = val
	}
}

func New(opts ...rsaxSecurityOptions) *RSAXSecurity {

	options := &rsaxSecurityoptions{}

	for _, v := range opts {
		v(options)
	}

	ins := &RSAXSecurity{options: *options}
	_ = ins

	if ins.options.pubStr != "" {
		ins.pubkey = utilx.MustSucc(getPKPubKey([]byte(ins.options.pubStr)))
	}

	if ins.options.priStr != "" {
		ins.prikey = utilx.MustSucc(getPKPrivkey([]byte(ins.options.priStr)))
	}

	return ins
}

// * 使用公钥加密
func (rs *RSAXSecurity) PubKeyEncode(input []byte) ([]byte, error) {

	if rs.pubkey == nil {
		return []byte(""), ErrNoPrivatekeySet
	}

	output := bytes.NewBuffer(nil)
	err := pubKeyIO(rs.pubkey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// * 公钥解密
func (rs *RSAXSecurity) PubKeyDecode(input []byte) ([]byte, error) {
	if rs.pubkey == nil {
		return []byte(""), ErrNoPrivatekeySet
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(rs.pubkey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// * 私钥加密
func (rs *RSAXSecurity) PriKeyEncode(input []byte) ([]byte, error) {
	if rs.prikey == nil {
		return []byte(""), ErrNoPrivatekeySet
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(rs.prikey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// * 私钥解密
func (rs *RSAXSecurity) PriKeyDecode(input []byte) ([]byte, error) {
	if rs.prikey == nil {
		return []byte(""), ErrNoPrivatekeySet
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(rs.prikey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}

	return io.ReadAll(output)
}

// * 使用RSAWithMD5算法签名
func (rs *RSAXSecurity) SignMd5WithRsa(data string) (string, error) {
	md5Hash := md5.New()
	s_data := []byte(data)
	md5Hash.Write(s_data)
	hashed := md5Hash.Sum(nil)

	signByte, err := rsa.SignPKCS1v15(rand.Reader, rs.prikey, crypto.MD5, hashed)
	sign := base64.StdEncoding.EncodeToString(signByte)
	return string(sign), err
}

// * 使用RSAWithSHA1算法签名
func (rs *RSAXSecurity) SignSha1WithRsa(data string) (string, error) {
	sha1Hash := sha1.New()
	s_data := []byte(data)
	sha1Hash.Write(s_data)
	hashed := sha1Hash.Sum(nil)

	signByte, err := rsa.SignPKCS1v15(rand.Reader, rs.prikey, crypto.SHA1, hashed)
	sign := base64.StdEncoding.EncodeToString(signByte)
	return string(sign), err
}

// * 使用RSAWithSHA256算法签名
func (rs *RSAXSecurity) SignSha256WithRsa(data string) (string, error) {
	sha256Hash := sha256.New()
	s_data := []byte(data)
	sha256Hash.Write(s_data)
	hashed := sha256Hash.Sum(nil)

	signByte, err := rsa.SignPKCS1v15(rand.Reader, rs.prikey, crypto.SHA256, hashed)
	sign := base64.StdEncoding.EncodeToString(signByte)
	return string(sign), err
}

// * 使用RSAWithMD5验证签名
func (rs *RSAXSecurity) VerifySignMd5WithRsa(data string, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	hash := md5.New()
	hash.Write([]byte(data))
	return rsa.VerifyPKCS1v15(rs.pubkey, crypto.MD5, hash.Sum(nil), sign)
}

// * 使用RSAWithSHA1验证签名
func (rs *RSAXSecurity) VerifySignSha1WithRsa(data string, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(data))
	return rsa.VerifyPKCS1v15(rs.pubkey, crypto.SHA1, hash.Sum(nil), sign)
}

// * 使用RSAWithSHA256验证签名
func (rs *RSAXSecurity) VerifySignSha256WithRsa(data string, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write([]byte(data))

	return rsa.VerifyPKCS1v15(rs.pubkey, crypto.SHA256, hash.Sum(nil), sign)
}
