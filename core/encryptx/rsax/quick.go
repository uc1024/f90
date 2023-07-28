package rsax

import (
	"encoding/base64"
	"encoding/hex"
	"sort"

)

// * 公钥加密
func PublicEncrypt(data, public_key string) (string, error) {

	grsa := New(SetPublicString(public_key))

	rsadata, err := grsa.PubKeyEncode([]byte(data))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(rsadata), nil
}

// * 私钥加密
func PriKeyEncrypt(data, private_key string) (string, error) {

	grsa := New(SetPrivateString(private_key))

	rsadata, err := grsa.PriKeyEncode([]byte(data))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(rsadata), nil
}

// * 公钥解密
func PublicDecrypt(data, public_key string) (string, error) {

	databs, _ := base64.StdEncoding.DecodeString(data)

	grsa := New(SetPublicString(public_key))

	rsadata, err := grsa.PubKeyDecode(databs)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rsadata), nil
}

// * 私钥解密
func PriKeyDecrypt(data, private_key string) (string, error) {

	databs, _ := base64.StdEncoding.DecodeString(data)

	grsa := New(SetPrivateString(private_key))

	rsadata, err := grsa.PriKeyDecode(databs)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rsadata), nil
}

// * 使用RSAWithMD5算法签名
func SignMd5WithRsa(data string, private_key string) (string, error) {
	grsa := New(SetPrivateString(private_key))

	sign, err := grsa.SignMd5WithRsa(data)
	if err != nil {
		return "", err
	}

	return sign, err
}

// * 使用RSAWithSHA1算法签名
func SignSha1WithRsa(data string, private_key string) (string, error) {
	grsa := New(SetPrivateString(private_key))

	sign, err := grsa.SignSha1WithRsa(data)
	if err != nil {
		return "", err
	}

	return sign, err
}

// * 使用RSAWithSHA256算法签名
func SignSha256WithRsa(data string, private_key string) (string, error) {
	grsa := New(SetPrivateString(private_key))

	sign, err := grsa.SignSha256WithRsa(data)
	if err != nil {
		return "", err
	}
	return sign, err
}

// * 使用RSAWithMD5验证签名
func VerifySignMd5WithRsa(data string, sign_data string, public_key string) error {
	grsa := New(SetPublicString(public_key))
	return grsa.VerifySignMd5WithRsa(data, sign_data)
}

// * 使用RSAWithSHA1验证签名
func VerifySignSha1WithRsa(data string, sign_data string, public_key string) error {
	grsa := New(SetPublicString(public_key))
	return grsa.VerifySignSha1WithRsa(data, sign_data)
}

// * 使用RSAWithSHA256验证签名
func VerifySignSha256WithRsa(data string, sign_data string, public_key string) error {
	grsa := New(SetPublicString(public_key))
	return grsa.VerifySignSha256WithRsa(data, sign_data)
}

// * 根据key的排序输出拼接value
func MapSortToVal[T any](m map[string]T) []T{
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i][0] < keys[j][0]
	})
	value:=make([]T,len(m))
	for i := 0; i < len(keys); i++ {
		value[i]=m[keys[i]]
	}
	return value
}


