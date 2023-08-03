package notifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/notifications/templates"
)

// POP3/SMTP协议
// 接收邮件服务器：pop.exmail.qq.com ，使用SSL，端口号995
// 发送邮件服务器：smtp.exmail.qq.com ，使用SSL，端口号465
func TestEmail(t *testing.T) {

	em := NewEmail(EmailxOptions{
		SmtpServer:   "smtp.exmail.qq.com",
		SmtpPort:     465,
		SmtpUsername: "liuyonglong@uc1024.com",
		SmtpPassword: "WFALviTcWu5WCSAU",
	})

	// 生成验证码模板
	temp, err := templates.GenerateCodeEmail(&templates.CodeTemplatesData{
		Code: "123456",
	})

	assert.NoError(t, err)

	err = em.Send("372572571@qq.com", "邮箱短信验证", temp)
	assert.NoError(t, err)
}
