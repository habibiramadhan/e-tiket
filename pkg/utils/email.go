//pkg/utils/email.go

package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	FromName string
}

type EmailData struct {
	To      []string
	Subject string
	Body    string
}

func SendEmail(config SMTPConfig, data EmailData) error {
	if config.Host == "" || config.Port == "" {
		return fmt.Errorf("konfigurasi SMTP tidak lengkap")
	}

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	
	to := strings.Join(data.To, ",")
	
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s\r\n", to, data.Subject, data.Body))
	
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	
	// Kirim email dengan TLS jika port adalah 465 (SMTPS)
	if config.Port == "465" {
		// Buat koneksi TLS
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         config.Host,
		}
		
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("gagal membuat koneksi TLS: %w", err)
		}
		defer conn.Close()
		
		client, err := smtp.NewClient(conn, config.Host)
		if err != nil {
			return fmt.Errorf("gagal membuat client SMTP: %w", err)
		}
		defer client.Close()
		
		// Autentikasi
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("gagal autentikasi: %w", err)
		}
		
		// Set pengirim dan penerima
		if err = client.Mail(config.Username); err != nil {
			return fmt.Errorf("gagal set pengirim: %w", err)
		}
		
		for _, recipient := range data.To {
			if err = client.Rcpt(recipient); err != nil {
				return fmt.Errorf("gagal set penerima %s: %w", recipient, err)
			}
		}
		
		// Kirim email
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("gagal persiapan pengiriman: %w", err)
		}
		
		_, err = w.Write(msg)
		if err != nil {
			return fmt.Errorf("gagal menulis pesan: %w", err)
		}
		
		err = w.Close()
		if err != nil {
			return fmt.Errorf("gagal menutup pesan: %w", err)
		}
		
		return client.Quit()
	} else {
		// Gunakan metode standar untuk port lain (587, 25, dll)
		return smtp.SendMail(addr, auth, config.Username, data.To, msg)
	}
}

func ParseTemplate(templatePath string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}