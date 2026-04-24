package sendemail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"lms-be/configs"
	"lms-be/shared/logger"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

func PrimitiveSendMail(config *configs.Config, to []string, cc []string, subject, message string) error {
	body := "From: " + config.App.SMTP.SenderName + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", config.App.SMTP.AuthEmail, config.App.SMTP.AuthPassword, config.App.SMTP.Host)
	smtpAddr := fmt.Sprintf("%s:%d", config.App.SMTP.Host, config.App.SMTP.Port)

	err := smtp.SendMail(smtpAddr, auth, config.App.SMTP.AuthEmail, append(to, cc...), []byte(body))
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

func ParseTemplate(templateFileName string, data interface{}) (body string, err error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return
	}
	body = buf.String()

	return
}

func encodeBase64(dst *bytes.Buffer, r io.Reader) error {
	encoder := base64.NewEncoder(base64.StdEncoding, dst)
	if _, err := io.Copy(encoder, r); err != nil {
		return err
	}
	return encoder.Close()
}

func SendMail(config *configs.Config, to []string, cc []string, subject, message string, attachments []string) error {
	// setup smtp
	smtpAddr := fmt.Sprintf("%s:%d", config.App.SMTP.Host, config.App.SMTP.Port)
	auth := smtp.PlainAuth("", config.App.SMTP.AuthEmail, config.App.SMTP.AuthPassword, config.App.SMTP.Host)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// header umum
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", config.App.SMTP.SenderName, config.App.SMTP.AuthEmail)
	headers["To"] = strings.Join(to, ", ")
	if len(cc) > 0 {
		headers["Cc"] = strings.Join(cc, ", ")
	}
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=" + writer.Boundary()

	for k, v := range headers {
		fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(&buf, "\r\n")

	// body part (HTML)
	altWriter := multipart.NewWriter(&buf)
	altHeader := make(map[string]string)
	altHeader["Content-Type"] = "multipart/alternative; boundary=" + altWriter.Boundary()

	fmt.Fprintf(&buf, "--%s\r\n", writer.Boundary())
	for k, v := range altHeader {
		fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(&buf, "\r\n")

	// HTML part
	htmlHeader := make(map[string]string)
	htmlHeader["Content-Type"] = `text/html; charset="UTF-8"`
	htmlHeader["Content-Transfer-Encoding"] = "quoted-printable"
	fmt.Fprintf(&buf, "--%s\r\n", altWriter.Boundary())
	for k, v := range htmlHeader {
		fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(&buf, "\r\n")

	qp := quotedprintable.NewWriter(&buf)
	qp.Write([]byte(message))
	qp.Close()
	fmt.Fprintf(&buf, "\r\n")

	// end alternative
	fmt.Fprintf(&buf, "--%s--\r\n", altWriter.Boundary())

	// attachments
	for _, filePath := range attachments {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		partHeader := make(map[string]string)
		partHeader["Content-Type"] = "application/octet-stream"
		partHeader["Content-Disposition"] = fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath))
		partHeader["Content-Transfer-Encoding"] = "base64"

		fmt.Fprintf(&buf, "--%s\r\n", writer.Boundary())
		for k, v := range partHeader {
			fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
		}
		fmt.Fprintf(&buf, "\r\n")

		// encode file ke base64
		if err := encodeBase64(&buf, file); err != nil {
			return err
		}
		fmt.Fprintf(&buf, "\r\n")
	}

	// end mixed
	fmt.Fprintf(&buf, "--%s--\r\n", writer.Boundary())

	// kirim
	allRecipients := append(to, cc...)
	return smtp.SendMail(smtpAddr, auth, config.App.SMTP.AuthEmail, allRecipients, buf.Bytes())
}
