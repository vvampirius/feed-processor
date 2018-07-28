package email

import (
	"net/smtp"
	"net/mail"
	"log"
)

func SendMessage(c *smtp.Client, from *mail.Address, to *mail.Address, message []byte) error {
	if err := c.Mail(from.Address); err != nil { return err }

	if err := c.Rcpt(to.Address); err != nil { return err }

	w, err := c.Data()
	if err != nil {	return err }

	if _, err := w.Write(message); err != nil { log.Println(err) }

	return w.Close()
}
