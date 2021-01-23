package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"gopkg.in/mail.v2"
)

// SMTPConfig used for saving SMTP mail server config
type SMTPConfig struct {
	MailServer string
	Port       string
	Address    string
	Password   string
}

// GetSupportConfig this function gets support mail server config from environment variables
func GetSupportConfig() SMTPConfig {
	return SMTPConfig{
		MailServer: os.Getenv("MAIL_SERVER"),
		Port:       os.Getenv("MAIL_SERVER_PORT"),
		Address:    os.Getenv("SUPPORT_MAIL_ADDRESS"),
		Password:   os.Getenv("SUPPORT_MAIL_PASS"),
	}
}

// SendSupport send HTML email from support mail server
func SendSupport(subject, msgPlain, msgHTML, email string) error {
	conf := GetSupportConfig()

	m := mail.NewMessage()
	m.SetHeader("From", conf.Address)
	m.SetHeader("To", email)

	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", msgPlain)
	m.AddAlternative("text/html", msgHTML)

	// m.SetAddressHeader("Cc", "support@edgecom.io", "Edgecom Support")
	// m.Attach("/home/Alex/lolcat.jpg")

	d := mail.NewDialer(conf.MailServer, 587, conf.Address, conf.Password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return d.DialAndSend(m)
}

// SendWelcome sends a welcome email to the new registered user
func SendWelcome(name, email string) error {
	subject := "Welcome To Edgecom Energy Portal\n"

	if name == "" {
		name = "Customer"
	}

	data := struct {
		Name string
	}{
		Name: name,
	}

	msgHTML, err := readWelcomeHTML(data)
	if err != nil {
		return fmt.Errorf("unable to read HTML welcome message: %v", err)
	}

	msgPlain, err := readWelcomePlain()
	if err != nil {
		return fmt.Errorf("unable to read plain welcome message: %v", err)
	}

	return SendSupport(subject, msgPlain, msgHTML, email)
}

// SendResetEmail sends reset password email
func SendResetEmail(userID int64, userName, link, email string) error {
	subject := "Reset Edgecom password\n"

	if userName == "" {
		userName = "Customer"
	}

	data := struct {
		Name string
		Link string
	}{
		Name: userName,
		Link: link,
	}

	msgHTML, err := readResetPassHTML(data)
	if err != nil {
		return fmt.Errorf("unable to read HTML welcome message: %v", err)
	}

	msgPlain := "Reset your password with following link"

	if err := SendSupport(subject, msgPlain, msgHTML, email); err != nil {
		return err
	}

	if err := AddHistory(userID, "Reset Password"); err != nil {
		log.Println("error while saving email history: ", err)
	}

	return nil
}

// SendResetEmail sends reset password email
func SendWelcomeAndResetEmail(userID int64, userName, username, link, email string) error {
	subject := "Welcome To Edgecom Energy Portal\n"

	if userName == "" {
		userName = "Customer"
	}

	data := struct {
		Name     string
		Link     string
		Username string
	}{
		Name:     userName,
		Link:     link,
		Username: username,
	}

	msgHTML, err := readWelcomeAndResetPassHTML(data)
	if err != nil {
		return fmt.Errorf("unable to read HTML welcome message: %v", err)
	}

	msgPlain := "Reset your password with following link"

	if err := SendSupport(subject, msgPlain, msgHTML, email); err != nil {
		return err
	}

	if err := AddHistory(userID, "Welcome"); err != nil {
		log.Println("error while saving email history: ", err)
	}

	return nil
}

func readWelcomeHTML(data interface{}) (string, error) {
	welcome, err := template.ParseFiles("./email/welcome.html")
	if err != nil {
		return "", fmt.Errorf("unable to read welcome HTMl template: %v", err)
	}

	var tpl bytes.Buffer
	if err := welcome.Execute(&tpl, data); err != nil {
		return "", fmt.Errorf("error while rendering HTML template: %v", err)
	}

	return tpl.String(), nil
}

func readResetPassHTML(data interface{}) (string, error) {
	reset, err := template.ParseFiles("./email/reset_pass.html")
	if err != nil {
		return "", fmt.Errorf("unable to read reset password template: %v", err)
	}

	var tpl bytes.Buffer
	if err := reset.Execute(&tpl, data); err != nil {
		return "", fmt.Errorf("error while rendering HTML template: %v", err)
	}

	return tpl.String(), nil
}

func readWelcomeAndResetPassHTML(data interface{}) (string, error) {
	reset, err := template.ParseFiles("./email/welcome_and_reset.html")
	if err != nil {
		return "", fmt.Errorf("unable to read welcome and reset password template: %v", err)
	}

	var tpl bytes.Buffer
	if err := reset.Execute(&tpl, data); err != nil {
		return "", fmt.Errorf("error while rendering HTML template: %v", err)
	}

	return tpl.String(), nil
}

func readWelcomePlain() (string, error) {
	// No need to close the file.
	content, err := ioutil.ReadFile("./email/welcome.txt")
	if err != nil {
		return "", err
	}

	// Convert []byte to string
	return string(content), nil
}
