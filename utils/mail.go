package utils

import (
	"auxpi/bootstrap"
	"bytes"
	"fmt"
	"html/template"
	"strconv"

	"github.com/astaxie/beego/logs"
	"gopkg.in/gomail.v2"
)

type MailInfo struct {
	SiteName   string
	Logo       string
	SiteLink   string
	UserIndex  string
	UserCenter string
	Active     string
	Content    string
}

type Result struct {
	output string
}

func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": bootstrap.SiteConfig.MailConfig.User,
		"pass": bootstrap.SiteConfig.MailConfig.Pass,
		"host": bootstrap.SiteConfig.MailConfig.Host,
		"port": bootstrap.SiteConfig.MailConfig.Port,
	}

	port, _ := strconv.Atoi(mailConn["port"])

	m := gomail.NewMessage()
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

func RenderMail(tpl string, args ...string) string {

	if len(args) == 0 {
		return ""
	}
	token := ""
	active := ""
	mailContent := ""

	switch args[0] {
	case "register":
		token = args[1]
		active = bootstrap.SiteConfig.SiteUrl + "register/active/" + token
		mailContent = `感谢您的注册，请点击下方的链接完成注册(如果不能再浏览器中打开，您可以直接复制到浏览器中然后进行访问) : `
	case "reset":
		token = args[1]
		mailContent = `如果您不知道这封邮件是做什么请直删除邮件即可,点击下方链接找回您的密码(如果不能再浏览器中打开，您可以直接复制到浏览器中然后进行访问) : `
		active = bootstrap.SiteConfig.SiteUrl + "reset/" + token
		//共用这个就可以
		tpl = "register.tpl"

	}

	var mailInfo = MailInfo{
		SiteName:   bootstrap.SiteConfig.SiteName,
		SiteLink:   bootstrap.SiteConfig.SiteUrl,
		UserIndex:  bootstrap.SiteConfig.SiteUrl + "users/index",
		UserCenter: bootstrap.SiteConfig.SiteUrl + "users/edit",
		Logo:       "x",
		Active:     active,
		Content:    mailContent,
	}
	dir := "views/mail/" + tpl
	var t *template.Template
	buf := new(bytes.Buffer)

	t, err := t.ParseFiles(dir)
	if err != nil {
		logs.Alert(err)
		return ""
	}
	err = t.Execute(buf, mailInfo)
	if err != nil {
		logs.Alert(err)
		return ""
	}
	return buf.String()

}

func (p *Result) Write(b []byte) (n int, err error) {
	fmt.Println("called by template")
	p.output += string(b)
	return len(b), nil
}
