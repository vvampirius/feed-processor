package email

import (
	"log"
	"net/mail"
)

func SendMessages(hostPort string, username string, password string, from *mail.Address, to *mail.Address, messages [][]byte) error {
	connection, client, err := ConnectTLS(hostPort, username, password)
	if err != nil { return err }

	for _, message := range messages {
		SendMessage(client, from, to, message)
	}

	if err := connection.Close(); err!=nil { log.Println(err) }
	return nil
}