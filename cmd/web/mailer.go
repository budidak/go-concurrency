package main

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Mail is the data type for the mail server
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Wait        *sync.WaitGroup
	MailerChan  chan Message // arka planda maili buraya göndermek istiyoruz.
	ErrorChan   chan error   // hata oluştuğunda
	DoneChan    chan bool    // işlem tamamlandığında
}

// Message is the data type for the mails
type Message struct {
	From          string
	FromName      string
	To            string
	Subject       string
	Attachments   []string
	AttachmentMap map[string]string
	Data          any            // body
	DataMap       map[string]any // getting data to the template
	Template      string         // bu verileri göndereceğimiz template
}

// listenForMail continously listens for messages on the MailerChan, or error on the ErrorChan etc.
func (app *Config) listenForMail() {
	for {
		select {
		// MailerChan içinde mail varsa alacağız ve bu case gerçekleşirse;
		// arkaplanda goroutine ile sendMail() ile maili göndereceğiz.
		case msg := <-app.Mailer.MailerChan:
			go app.Mailer.sendMail(msg, app.Mailer.ErrorChan)
		// olası bir hata oluşması durumunda ErrorChan'den error alacağı ve handle edeceğiz.
		// sendMail() fonksiyonunu yazarken öyle handle etmiştik errorları.
		case err := <-app.Mailer.ErrorChan:
			app.ErrorLog.Println(err)
		// eğer DoneChan'e tamamlanmayla ilgili bilgi gelirse programı sonlandıracağız.
		case <-app.Mailer.DoneChan:
			return
		}
	}
}

// createMail function sets mail configurations and returns a Mail
// buradan dönen m değeri, Config içerisindeki Mail tipindeki değişkene atanacak.
// böylelikle aynı ayarları listenForMail() içerisinde yazmadan Config içinden çekip kullanabileceğiz.
func (app *Config) createMail() Mail {
	// create channels
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100) // buffered channel, can hold 100 Messages
	doneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromName:    "Info",
		FromAddress: "info@mycompany.com",
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    doneChan,
	}

	return m
}

// sendMail function sets up SMTP configurations and sends an email
// bunu programda otomatik olarak mail göndermek için kullanacağız, hatalı girişler veya üyelik gibi...
func (m *Mail) sendMail(msg Message, errorChan chan error) {

	defer m.Wait.Done() // decrement wait group counter by 1 at the end of function

	// if there is no template for the given message, use default
	if msg.Template == "" {
		msg.Template = "mail"
	}

	// if there is no FromAddress for the given message, use receiver's as default
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	// if there is no FromName for the given message, use receiver's as default
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// if there is no AttachmentMap, create empty map
	if msg.AttachmentMap == nil {
		msg.AttachmentMap = make(map[string]string)
	}

	// if there is no data in DataMap, create empty map
	if len(msg.DataMap) == 0 {
		msg.DataMap = make(map[string]any)
	}

	msg.DataMap["message"] = msg.Data
	// "message" kullandık; çünkü mail.html.gotmpl/mail.plain.gotmpl dosyalarında o isimde kullandık.

	// build html mail
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}

	// build plain text mail
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		errorChan <- err
	}

	// create a SMTP server and connect it
	server := mail.NewSMTPClient()

	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}

	// create a new email message and set its properties
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// if there is any attachment add it to the email
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	if len(msg.AttachmentMap) > 0 {
		for key, val := range msg.AttachmentMap {
			email.AddAttachment(val, key)
		}
	}

	// send email
	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}
}

// buildHTMLMessage function creates a html style for the mail
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.html.gotmpl", msg.Template)
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer // execute the template into that
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String() // contents of buffer -> string
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// buildPlainTextMessage function creates a plain text for the mail
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.plain.gotmpl", msg.Template)
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer // execute the template into that
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String() // contents of buffer -> string
	return plainMessage, nil
}

// inlineCSS function makes the CSS more acceptable to various email clients
func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

// getEncryption function allows to get the encrytion method for a particular server
func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none": // use this in "development phase"
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
